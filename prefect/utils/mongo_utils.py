import os
import logging

from pymongo import MongoClient
from bson import ObjectId
from model.client_article import ClientArticle
from dotenv import load_dotenv

load_dotenv()

mongo_client = MongoClient(os.getenv("MONGO_URI"))
db_name = "client-factpack"
db = mongo_client[db_name]

article_collection = "articles"
clients_collection = "clients"


def put_article(article: ClientArticle):
    articles = db[article_collection]
    result = articles.insert_one(article.model_dump())
    id = result.inserted_id
    return id


def update_client_article(client_id: str, article_id: ObjectId):
    try:
        _id = ObjectId(client_id)
    except Exception:
        logging.error("Invalid clientId: Unable to convert to ObjectId", exc_info=True)
        return []

    clients = db[clients_collection]
    result = clients.find_one_and_update(
        {"_id": _id},
        {"$addToSet": {"articles": article_id}},
        projection={"data.profile.names": 1},
    )

    if not result:
        logging.warning(f"Client with ID {_id} not found")

    return result["data"]["profile"]["names"]


def fetch_mongo_records_by_ids(ids: list):
    mongo_uri = os.getenv("MONGO_URI")
    db_name = os.getenv("DB_NAME")
    collection_name = os.getenv("CLIENT_COLLECTION_NAME")
    client = MongoClient(mongo_uri)
    collection = client[db_name][collection_name]
    results = list(collection.find({"_id": {"$in": ids}}))
    client.close()
    return results


def get_all_client_primary_names() -> list[str]:
    clients = mongo_client[db_name][clients_collection]

    primary_names = []
    for doc in clients.find({}, {"data.profile.names": 1}):
        try:
            primary = doc["data"]["profile"]["names"][0]
            primary_names.append(primary)
        except (KeyError, IndexError):
            continue

    return primary_names
