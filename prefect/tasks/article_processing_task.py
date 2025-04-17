import os
import ssl
import json
import httpx
from datetime import datetime, timedelta

import pika
from openai import OpenAI
from dotenv import load_dotenv
from prefect import task, get_run_logger
from bs4 import BeautifulSoup
from transformers import BartForConditionalGeneration, BartTokenizer
from newsapi import NewsApiClient

from model.client import ClientProfile
from utils.mongo_utils import get_all_client_primary_names

load_dotenv()
NEWS_API_KEY = os.getenv("NEWS_API_KEY")
MODEL_NAME = os.getenv("SUMMARIZER_MODEL")

OPENAI_ASSISTANT_ID = os.getenv("OPENAI_ASSISTANT_ID")
OPENAI_MODEL = os.getenv("OPENAI_MODEL")

tokenizer = BartTokenizer.from_pretrained(MODEL_NAME)
model = BartForConditionalGeneration.from_pretrained(MODEL_NAME)
newsapi_client = NewsApiClient(api_key=NEWS_API_KEY)

RABBITMQ_HOST = os.getenv("RABBITMQ_HOST")
RABBITMQ_PORT = os.getenv("RABBITMQ_PORT")
RABBITMQ_USER = os.getenv("RABBITMQ_USER")
RABBITMQ_PASSWORD = os.getenv("RABBITMQ_PASSWORD")

client = OpenAI()

QUEUE_NAME = "notifications"


@task
def send_to_queue(news_data):
    """Sends the summarized news data to the RabbitMQ queue."""
    logger = get_run_logger()
    credentials = pika.PlainCredentials(RABBITMQ_USER, RABBITMQ_PASSWORD)
    ssl_context = ssl.create_default_context()
    parameters = pika.ConnectionParameters(
        host=RABBITMQ_HOST,
        port=RABBITMQ_PORT,
        virtual_host="/",
        credentials=credentials,
        ssl_options=pika.SSLOptions(ssl_context),
        heartbeat=600,
        blocked_connection_timeout=300,
    )
    connection = pika.BlockingConnection(parameters)
    channel = connection.channel()

    channel.queue_declare(queue=QUEUE_NAME, durable=True)

    channel.basic_publish(
        exchange="",
        routing_key=QUEUE_NAME,
        body=json.dumps(news_data),
        properties=pika.BasicProperties(delivery_mode=2),
    )

    logger.info(f"[x] Sent news to queue for client: {news_data['clientName']}")
    connection.close()


@task
def get_articles(kw: str):
    ystd = (datetime.now() - timedelta(days=1))
    start = ystd.replace(hour=0, minute=0)
    end = ystd.replace(hour=23, minute=59)
    articles = newsapi_client.get_everything(
        q=kw,
        from_param=start,
        to=end,
        language='en',
        sort_by='relevancy',
        page=1,
        page_size=2,
    )
    return articles["articles"]


@task
def scrape_article(url) -> str:
    """Scrapes the full article text from the given URL using httpx."""
    try:
        headers = {"User-Agent": "Mozilla/5.0"}
        with httpx.Client(timeout=5.0) as client:
            response = client.get(url, headers=headers)
        if response.status_code == 200:
            soup = BeautifulSoup(response.text, "html.parser")
            paragraphs = soup.find_all("p")
            return " ".join([p.get_text() for p in paragraphs]).strip()
        return f"Error fetching article: {response.status_code}"
    except Exception as e:
        return f"Scraping failed: {str(e)}"


@task
def summarize_text(text, max_length=300, min_length=150):
    """Summarizes the given text using the BART model."""
    inputs = tokenizer.encode(
        "summarize: " + text, return_tensors="pt", max_length=1024, truncation=True
    )
    summary_ids = model.generate(
        inputs,
        max_length=max_length,
        min_length=min_length,
        length_penalty=2.0,
        num_beams=4,
        early_stopping=True,
    )
    return tokenizer.decode(summary_ids[0], skip_special_tokens=True)


@task
def get_clients() -> list[str]:
    return get_all_client_primary_names()

@task
def extract_client_info(text: str) -> ClientProfile:
    print("extracting client info")
    completion = client.beta.chat.completions.parse(
        model=OPENAI_MODEL,
        messages=[
            {
                "role": "system",
                "content": "Extract relevant infomation out of the user raw text in a JSON object."
                "Ensure that all key fields are filled if they can be reasonably inferred — especially the "
                "'ownedCompanies' field. Ensure that the name is the full name.",
            },
            {"role": "user", "content": text},
        ],
        response_format=ClientProfile,
    )
    client_data = completion.choices[0].message.parsed
    print("extraction updated")
    print("extraction ended")
    return client_data
