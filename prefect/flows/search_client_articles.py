import logging
from datetime import datetime

import pika
from prefect import flow, task
from prefect.task_runners import ThreadPoolTaskRunner
from tasks.article_processing_task import *
from model.client_article import ClientArticle 

@flow(task_runner=ThreadPoolTaskRunner(max_workers=3))
def search_client(c: str):
    articles = get_articles.submit(kw=c, date=datetime.now().strftime("%Y-%m-%d")).result()
    if not articles:
        logging.info(f"No aticles for client: {c}")
        return
    print("Hello, trying to print")
    data: list[ClientArticle] = []
    for article in articles:
        title: str = article.get("title", "Unknown Title")
        url: str = article.get("url", "Unknown URL")
        source: str = article["source"].get("name", "Unknown Source")
        article_text: str = scrape_article.submit(url).result()
        if article_text.startswith("Error"):
            continue
        summary = summarize_text.submit(article_text).result()

        obj = ClientArticle(
            source=source,
            title=title,
            url=url,
            summary=summary
        )
        data.append(obj)
    print(data)
    message = {
        "client": c,
        "articles": [a.model_dump() for a in data]
    }
    print(message)
    send_to_queue.submit(message)

@flow
def update_clients():
    clients = get_clients.submit().result()
    if len(clients) == 0:
        logging.info("No articles fetched. Exiting.")
        return
    
    for c in clients:
        search_client(c)