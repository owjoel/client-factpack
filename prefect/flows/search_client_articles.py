import json

from prefect import flow, get_run_logger
from prefect.task_runners import ThreadPoolTaskRunner
from tasks.article_processing_task import *
from tasks.sentiment_task import *
from tasks.qdrant_task import search_profiles_by_json, transform_into_vector
from tasks.dedupe_task import *
from model.client_article import ClientArticle 
from utils.mongo_utils import *

@flow(task_runner=ThreadPoolTaskRunner(max_workers=3), log_prints=True)
def search_client(c: str):
    articles = get_articles.submit(c).result()
    if not articles:
        print(f"No articles for client: {c}")
        return
    
    data: list[ClientArticle] = []
    logger = get_run_logger()
    for article in articles:

        # scrape articles
        title: str = article.get("title", "Unknown Title")
        url: str = article.get("url", "Unknown URL")
        source: str = article["source"].get("name", "Unknown Source")
        article_text: str = scrape_article.submit(url).result()
        if article_text.startswith("Error"):
            continue

        # summarize, sentimet analysis
        summary = summarize_text.submit(article_text).result()
        sentiment = analyze_sentiment.submit(summary).result()
        obj = ClientArticle(
            source=source,
            title=title,
            url=url,
            summary=summary,
            sentiment=sentiment
        )

        # qdrant to match client
        matched_clients = search_profiles_by_json(obj.model_dump())
        logger.info(obj.sentiment.model_dump_json())
        for i in matched_clients:
            logger.info(json.dumps(i, ensure_ascii=False))
        data.append(obj)

        # Extract client info then dedupe
        client_info = extract_client_info.submit(summary).result()
        client_id = dedupe_against_mongo.submit(client_info.model_dump(), matched_clients).result()
        if not client_id:
            continue

        # upload article, update client, and send to queue
        article_id = put_article(obj)
        names = update_client_article(client_id, article_id)
        message = {
            "notificationType": "client",
            "title": obj.title,
            "source": obj.source,
            "clientId": client_id,
            "clientName": ';'.join(names),
            "priority": getPriority(sentiment.label),
        }
        send_to_queue.submit(message).result()
        logger.info(json.dumps(message))


@flow
def update_clients():
    logger = get_run_logger()
    clients = get_clients.submit().result()
    if len(clients) == 0:
        logger.info("No articles fetched. Exiting.")
        return
    
    for c in clients:
        search_client(c)