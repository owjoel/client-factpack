from qdrant_client import QdrantClient
from dotenv import load_dotenv
from sentence_transformers import SentenceTransformer
from qdrant_client.http.models import PointStruct
from prefect import task
import os
import json
import traceback
from uuid import uuid5, NAMESPACE_DNS

load_dotenv()

QDRANT_URL = os.getenv("QDRANT_URL")
QDRANT_API_KEY = os.getenv("QDRANT_API_KEY")

@task
def transform_into_vector(text: str):
    model = SentenceTransformer("BAAI/bge-large-en-v1.5")
    return model.encode(text).tolist()


@task
def upsert_text_to_qdrant(profile: dict, record_id: str, collection_name: str = "clients"):
    try:
        if not profile or not record_id:
            raise ValueError("Profile and record_id must be provided.")

        text = json.dumps(profile, sort_keys=True)

        vector = transform_into_vector(text)

        client = QdrantClient(url=QDRANT_URL, api_key=QDRANT_API_KEY)
        point = PointStruct(
            id=str(uuid5(NAMESPACE_DNS, record_id)),
            vector=vector,
            payload={"mongo_id": record_id, "profile": profile},
        )

        client.upsert(collection_name=collection_name, points=[point])
        print(f"[Qdrant] Inserted vector for record ID: {record_id}")

    except Exception as e:
        print(f"[❌] Failed to upsert into Qdrant: {e}")
        traceback.print_exc()
        raise


@task
def search_profiles_by_json(query_json: dict, collection_name: str = "clients"):
    try:
        if not query_json or not isinstance(query_json, dict):
            raise ValueError("query_json must be a non-empty dictionary.")

        query_vector = transform_into_vector(json.dumps(query_json))

        client = QdrantClient(url=QDRANT_URL, api_key=QDRANT_API_KEY)

        results = client.search(
            collection_name=collection_name,
            query_vector=query_vector,
            limit=5,
            with_payload=True,
            with_vectors=False,
        )

        matches = [
            {
                "id": match.payload.get("mongo_id"),
                "similarityScore": match.score,
            }
            for match in results
        ]

        print(f"[Qdrant] Found {len(matches)} matches")
        return matches

    except Exception as e:
        print(f"[❌] Failed to search profiles in Qdrant: {e}")
        traceback.print_exc()
        raise

@task
def search_articles(summarised_article: str, collection_name: str = "articles"):
    try:
        client = QdrantClient(url=QDRANT_URL, api_key=QDRANT_API_KEY)
        query_vector = transform_into_vector(summarised_article)

        results = client.search(
            collection_name=collection_name,
            query_vector=query_vector,
            limit=5,
            with_payload=True,
            with_vectors=False,
        )

        matches = [
            {
                "id": match.payload.get("id"),
                "similarityScore": match.score,
            }
            for match in results
        ]

        print(f"[Qdrant] Found {len(matches)} matches")
        return matches

    except Exception as e:
        print(f"[❌] Failed to search articles in Qdrant: {e}")
        traceback.print_exc()
        raise


# TODO: not required, can remove
@task
def search_profile_by_mongo_id(mongo_id: str):
    """
    Search a profile in Qdrant using the original Mongo ObjectID string.
    Uses UUIDv5 derived from the mongo_id as Qdrant point ID.
    """
    qdrant_id = uuid5(NAMESPACE_DNS, mongo_id)

    client = QdrantClient(url=QDRANT_URL, api_key=QDRANT_API_KEY)

    result = client.retrieve(collection_name=COLLECTION_NAME, ids=[qdrant_id])

    if not result:
        print(f"[❌] No Qdrant point found for Mongo ID: {mongo_id}")
        return None

    point = result[0]

    return {
        "mongo_id": point.payload.get("mongo_id"),
    }
