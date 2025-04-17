import os
from prefect import task
from datetime import datetime, timezone
import json
from pymongo import MongoClient
from dotenv import load_dotenv
from bson import ObjectId
from utils import wiki_utils, openai_utils, prompt_utils, file_utils

load_dotenv()
MONGO_URI = os.getenv("MONGO_URI")


@task
def get_cleaned_target(target: str) -> str:
    return target.lower().strip().replace(" ", "_")


@task
def get_wikipedia_text(target: str) -> str:
    url = wiki_utils.get_wikipedia_url(target)
    return wiki_utils.extract_clean_wikipedia_page(url)


@task
def generate_openai_response(wiki_text: str, target: str = "Unknown") -> str:
    prompt = prompt_utils.build_prompt_no_schema(
        f"This is the profile of {target}\n {wiki_text}"
    )
    print("Querying OpenAI...")
    return openai_utils.query_gpt4o(prompt)


@task
def parse_openai_response(response: str) -> dict:
    try:
        return json.loads(response)
    except json.JSONDecodeError as e:
        raise ValueError(f"Failed to parse OpenAI response: {e}")


@task
def save_files(target: str, wiki_text: str, profile_json: dict):
    file_utils.save_to_text(wiki_text, "wikipedia_texts", f"{target}_wikipedia.txt")
    file_utils.save_to_json(profile_json, "profiles", f"{target}_profile.json")


@task
def insert_into_mongo(profile_json: dict, target: str) -> str:
    with MongoClient(MONGO_URI) as client:
        db = client["client-factpack"]
        collection = db["clients"]

        document = {
            "data": profile_json,
            "metadata": {
                "scraped": True,
                "createdAt": datetime.now(timezone.utc),
                "updatedAt": datetime.now(timezone.utc),
            },
        }

        inserted_id = collection.insert_one(document).inserted_id
        print(f"[Inserted] {target} with ID: {inserted_id}")
        return str(inserted_id)


@task
def update_client_profile(client_id: str, profile_json: dict):
    with MongoClient(MONGO_URI) as client:
        db = client["client-factpack"]
        collection = db["clients"]
        result = collection.update_one(
            {"_id": ObjectId(client_id)},
            {
                "$set": {
                    "data": profile_json,
                    "metadata.updatedAt": datetime.now(timezone.utc),
                    "metadata.scraped": True,
                    "metadata.sources": ["wikipedia"],
                }
            },
        )

        if result.modified_count == 0:
            print(f"[LOG] No updates made to client: {client_id}")
        else:
            print(f"[LOG] Updated client {client_id} with new profile.")
