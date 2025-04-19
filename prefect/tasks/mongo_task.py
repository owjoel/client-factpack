from prefect import task
from pymongo import MongoClient
from dotenv import load_dotenv
import os
from bson import ObjectId
from datetime import datetime, timezone

load_dotenv()
MONGO_URI=os.getenv("MONGO_URI")

@task
def update_job_status(job_id: str, status: str, log_message: str = None):
    with MongoClient(MONGO_URI) as client:
        db = client["client-factpack"]
        collection = db["jobs"]

        job = collection.find_one({"_id": ObjectId(job_id)})

        if not job:
            raise ValueError(f"Job with ID {job_id} not found")

        # If logs is null or not a list, reset it
        if "logs" not in job or not isinstance(job["logs"], list):
            collection.update_one({"_id": ObjectId(job_id)}, {"$set": {"logs": []}})

        update_query = {
            "$set": {"status": status, "updatedAt": datetime.now(timezone.utc)}
        }

        if log_message:
            update_query["$push"] = {
                "logs": {
                    "message": log_message,
                    "timestamp": datetime.now(timezone.utc),
                }
            }

        result = collection.update_one({"_id": ObjectId(job_id)}, update_query)

        if result.modified_count == 0:
            raise RuntimeError(
                f"Failed to update job status or push logs for job ID: {job_id}"
            )


@task
def add_job_log(job_id: str, log_message: str):
    with MongoClient(MONGO_URI) as client:
        db = client["client-factpack"]
        collection = db["jobs"]

        result = collection.update_one(
            {"_id": ObjectId(job_id)},
            {
                "$push": {
                    "logs": {
                        "message": log_message,
                        "timestamp": datetime.now(timezone.utc),
                    }
                }
            },
        )

        if result.matched_count == 0:
            raise ValueError(f"Job with ID {job_id} not found")

@task
def update_job_match_results(job_id: str, match_results: list):
    with MongoClient(MONGO_URI) as client:
        db = client["client-factpack"]
        collection = db["jobs"]

        result = collection.update_one(
            {"_id": ObjectId(job_id)}, {"$set": {"matchResults": match_results}}
        )

        if result.matched_count == 0:
            raise ValueError(f"Job with ID {job_id} not found")


@task
def update_job_scrape_result(job_id: str, scrape_result: str):
    with MongoClient(MONGO_URI) as client:
        db = client["client-factpack"]
        collection = db["jobs"]

        result = collection.update_one(
            {"_id": ObjectId(job_id)}, {"$set": {"scrapeResult": scrape_result}}
        )

        if result.matched_count == 0:
            raise ValueError(f"Job with ID {job_id} not found")


@task
def get_client_names(id: str) -> list[str]:
    with MongoClient(MONGO_URI) as client:
        db = client["client-factpack"]
        collection = db["clients"]

        result = collection.find_one({"_id": ObjectId(id)}, {"data.profile.names": 1})

        if result and "data" in result and "profile" in result["data"]:
            return result["data"]["profile"].get("names", [])
        return []
