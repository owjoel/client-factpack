import requests
from datetime import datetime
import time
import pika
import json
import ssl
from bs4 import BeautifulSoup
from transformers import BartForConditionalGeneration, BartTokenizer
from newsapi import NewsApiClient

# RabbitMQ Configuration
RABBITMQ_HOST = "RABBIT_MQ_HOST"
RABBITMQ_PORT = 5671  # Use 5671 if connecting securely via SSL/TLS
RABBITMQ_USER = "user"
RABBITMQ_PASSWORD = "pw"
QUEUE_NAME = "news_queue"


# API Configuration
NEWS_API_KEY = "key"
newsapi = NewsApiClient(api_key=NEWS_API_KEY)

KEYWORDS = ["Elon Musk"]

# Fetch latest news articles on Elon Musk
# ARTICLES = newsapi.get_everything(
#     q='Elon Musk',
#     from_param='2025-03-09',
#     to='2025-03-09',
#     language='en',
#     sort_by='relevancy',
#     page=2
# )

# Initialize BART model for text summarization
MODEL_NAME = "facebook/bart-large-cnn"
tokenizer = BartTokenizer.from_pretrained(MODEL_NAME)
model = BartForConditionalGeneration.from_pretrained(MODEL_NAME)


def send_to_queue(news_data):
    """Sends the summarized news data to the RabbitMQ queue."""
    credentials = pika.PlainCredentials(RABBITMQ_USER, RABBITMQ_PASSWORD)
    ssl_context = ssl.create_default_context()
    parameters = pika.ConnectionParameters(
        host=RABBITMQ_HOST,
        port=RABBITMQ_PORT,
        virtual_host="/",
        credentials=credentials,
        ssl_options=pika.SSLOptions(ssl_context),
        heartbeat=600,  # Keep connection alive
        blocked_connection_timeout=300,  # Prevent timeouts
    )
    connection = pika.BlockingConnection(parameters)
    channel = connection.channel()

    # Declare the queue to ensure it exists
    channel.queue_declare(queue=QUEUE_NAME, durable=True)

    # Publish the message
    channel.basic_publish(
        exchange="",
        routing_key=QUEUE_NAME,
        body=json.dumps(news_data),
        properties=pika.BasicProperties(delivery_mode=2),  # Persistent message
    )

    print(f"[x] Sent news to queue: {news_data['title']}")
    connection.close()


def get_articles(keyword, date=datetime.now().strftime("%Y-%m-%d")):
    articles = newsapi.get_everything(
        q=keyword,
        from_param="2025-03-09",
        to="2025-03-09",
        language="en",
        sort_by="relevancy",
        page=2,
    )
    return articles["articles"]


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


def scrape_article_content(url):
    """Scrapes the full article text from the given URL."""
    try:
        headers = {"User-Agent": "Mozilla/5.0"}
        response = requests.get(url, headers=headers, timeout=5)
        if response.status_code == 200:
            soup = BeautifulSoup(response.text, "html.parser")
            paragraphs = soup.find_all("p")
            return " ".join([p.get_text() for p in paragraphs]).strip()
        return f"Error fetching article: {response.status_code}"
    except Exception as e:
        return f"Scraping failed: {str(e)}"


def run_news_scraper():
    """Fetches news articles, scrapes full content, summarizes, and prints the results."""
    print(f"Starting news scraper at {datetime.now()}...")

    for kw in KEYWORDS:
        articles = get_articles(kw)
        if not articles:
            print("No articles fetched. Exiting.")
            return

        count = 1
        for article in articles:
            # print(article)
            title = article.get("title", "Unknown Title")
            url = article.get("url", "Unknown URL")
            source = article["source"].get("name", "Unknown Source")
            print(f"\nFetching Article {count}: {title}\n")
            article_text = scrape_article_content(url)
            summary = summarize_text(article_text)

            print("Summary:")
            print(summary)
            print("-" * 80)

            news_data = {
                "source": source,
                "title": title,
                "url": url,
                "summary": summary,
            }

            send_to_queue(news_data)

            time.sleep(1)  # Prevent excessive API requests
            count += 1
            if count >= 2:
                break


if __name__ == "__main__":
    run_news_scraper()
