import os
import logging

from pymongo import MongoClient
from bson import ObjectId
from model.client_article import ClientArticle

mongo_client = MongoClient(os.getenv("MONGO_URI"))
db_name = 'client-factpack'
db = mongo_client[db_name]

article_collection = 'articles'
clients_collection = 'clients'

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
    result = clients.update_one(
        {"_id": _id},
        {"$addToSet": {"articles": article_id}}
    )
    return result.modified_count