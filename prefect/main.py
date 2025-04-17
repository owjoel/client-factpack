from flows.scrape_flow import scrape_client_flow
from flows.match_flow import match_client_flow
from pymongo import MongoClient
from dotenv import load_dotenv
import os
from tasks.qdrant_task import upsert_text_to_qdrant, create_clients_collection_in_qdrant

load_dotenv()
MONGO_URI = os.getenv("MONGO_URI")


def main():
    client = MongoClient(MONGO_URI)
    db = client["client-factpack"]
    collection = db["clients"]

    create_clients_collection_in_qdrant()

    for doc in collection.find():
        try:
            record_id = str(doc["_id"])
            profile = doc.get("data", {})
            print(record_id)
            upsert_text_to_qdrant(profile=profile, record_id=record_id)
        except Exception as e:
            print(f"[❌] Failed to process document {doc.get('_id')}: {e}")


if __name__ == "__main__":
    main()
