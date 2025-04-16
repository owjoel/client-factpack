import os
import logging

from pymongo import MongoClient
from bson import ObjectId
from model.client_article import ClientArticle
from dotenv import load_dotenv


mongo_client = MongoClient(os.getenv("MONGO_URI"))
db_name = "client-factpack"
db = mongo_client[db_name]

article_collection = "articles"
clients_collection = "clients"


def put_article(article: ClientArticle):
    articles = db[article_collection]
    result = articles.insert_one(article)
    id = result.inserted_id


def update_client_article(client_id: str, article_id: ObjectId):
    try:
        _id = ObjectId(client_id)
    except:
        logging.error("invalid clientId", exc_info=True)

    clients = db[clients_collection]
    result = clients.update_one({"_id": _id}, {"$addToSet": {"articles": article_id}})
    return result.modified_count


def fetch_mongo_records_by_ids(ids: list):
    load_dotenv()
    mongo_uri = os.getenv("MONGO_URI")
    db_name = os.getenv("DB_NAME")
    collection_name = os.getenv("CLIENT_COLLECTION_NAME")
    client = MongoClient(mongo_uri)
    collection = client[db_name][collection_name]
    results = list(collection.find({"_id": {"$in": ids}}))
    client.close()
    return results
