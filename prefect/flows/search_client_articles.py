import json
from datetime import datetime

from prefect import flow, get_run_logger
from prefect.task_runners import ThreadPoolTaskRunner
from tasks.article_processing_task import *
from tasks.sentiment_task import *
from model.client_article import ClientArticle 

@flow(task_runner=ThreadPoolTaskRunner(max_workers=3), log_prints=True)
def search_client(c: str):
    # logger = get_run_logger()
    articles = get_articles.submit(kw=c, date=datetime.now().strftime("%Y-%m-%d")).result()
    if not articles:
        print(f"No articles for client: {c}")
        return
    data: list[ClientArticle] = []
    for article in articles:
        title: str = article.get("title", "Unknown Title")
        url: str = article.get("url", "Unknown URL")
        source: str = article["source"].get("name", "Unknown Source")
        article_text: str = scrape_article.submit(url).result()
        if article_text.startswith("Error"):
            continue
        summary = summarize_text.submit(article_text).result()
        sentiment = analyze_sentiment.submit(summary).result()
        obj = ClientArticle(
            source=source,
            title=title,
            url=url,
            summary=summary
        )
        data.append(obj)
    message = {
        "client": c,
        "articles": [a.model_dump() for a in data]
    }
    send_to_queue.submit(message).result()

@flow
def update_clients():
    logger = get_run_logger()
    clients = get_clients.submit().result()
    if len(clients) == 0:
        logger.info("No articles fetched. Exiting.")
        return
    
    for c in clients:
        search_client(c)