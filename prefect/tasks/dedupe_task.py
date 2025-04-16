import pandas as pd
from pandas_dedupe import dedupe_dataframe
from prefect import task
from typing import Dict, Optional
from dedupe import variables
from bson import ObjectId
from utils import mongo_utils
import json


@task
def dedupe_against_mongo(new_record: Dict, qdrant_matches: list) -> Optional[list]:
    """
    Match a new record against existing records in MongoDB using pandas-dedupe.

    Args:
        new_record: The new client record to match (dict).
        qdrant_matches: A list of the ObjectIds of the top 5 matches from QDrant

    Returns:
        The best-matched existing record (if found), else None.
    """

    print("[DEBUG] New record for deduplication:")
    print(json.dumps(new_record, indent=2, default=str))

    mongo_ids = [ObjectId(m["id"]) for m in qdrant_matches if m.get("id")]

    # Connect to Mongo
    print(f"[MATCHING] Fetching {len(mongo_ids)} records from Mongo for dedupe...")
    candidate_records = mongo_utils.fetch_mongo_records_by_ids(mongo_ids)

    # Extract flattened data
    existing_records = [extract_dedupe_fields(doc) for doc in candidate_records]
    df = pd.DataFrame(existing_records)
    df["_is_existing"] = True

    new_flat = extract_dedupe_fields(new_record)
    new_df = pd.DataFrame([new_flat])
    new_df["_is_existing"] = False

    field_properties = [
        variables.String("Name"),
        variables.String("Primary Citizenship"),
        variables.String("Companies"),
    ]

    combined_df = pd.concat([df, new_df], ignore_index=True)

    # Remove rows where all match fields are empty
    # combined_df = combined_df.dropna(
    #     subset=["Name", "Primary Citizenship", "Companies"], how="all"
    # )

    print("[DEBUG] Combined DataFrame after dropna:")
    print(combined_df.to_string(index=False))

    deduped_df = dedupe_dataframe(
        combined_df,
        field_properties=field_properties,
        config_name="dedupe_dataframe",
        update_model=False,
        threshold=0.4,
        sample_size=0.3,
        canonicalize=True,
    )

    # print(deduped_df.head(3).to_dict())
    # print("\nClustering output:")
    # print(deduped_df[["Name", "_is_existing", "cluster id", "confidence"]])

    deduped_df["_is_existing"] = deduped_df["_is_existing"].map(
        lambda x: str(x).lower() == "true"
    )

    new_id = deduped_df[~deduped_df["_is_existing"]]["cluster id"].iloc[0]
    match = deduped_df[
        (deduped_df["_is_existing"]) & (deduped_df["cluster id"] == new_id)
    ]

    return match.iloc[0].to_dict() if not match.empty else None


def extract_dedupe_fields(doc: dict) -> dict:
    profile = doc.get("data", {}).get("profile", {})
    owned_companies = doc.get("data", {}).get("ownedCompanies", [])

    names = " , ".join(profile.get("names") or []).strip()
    nationality = (profile.get("nationality") or "").strip()
    companies = " , ".join(
        [c.get("name", "") for c in (owned_companies or []) if "name" in c]
    ).strip()

    return {
        "Name": names if names else "",
        "Primary Citizenship": nationality if nationality else "",
        "Companies": companies if companies else "",
    }


# test_records = [
#     {
#         "data": {
#             "profile": {
#                 "names": ["Elon Musk"],
#                 "nationality": "American",
#             },
#             "ownedCompanies": [
#                 {"name": "Tesla"},
#             ],
#         }
#     },
#     {
#         "data": {
#             "profile": {
#                 "names": ["Jeff Bezos"],
#                 "nationality": "American",
#             },
#             "ownedCompanies": [
#                 {"name": "Amazon"},
#             ],
#         }
#     },
#     {
#         "data": {
#             "profile": {
#                 "names": ["Bernard Arnault"],
#                 "nationality": "French",
#             },
#             "ownedCompanies": [
#                 {"name": "LVMH"},
#             ],
#         }
#     },
#     {
#         "data": {
#             "profile": {"names": ["Mark Zuckerberg"], "nationality": "American"},
#             "ownedCompanies": [
#                 {"name": "Meta"},
#             ],
#         }
#     },
#     {
#         "data": {
#             "profile": {
#                 "names": ["Larry Page", "Lawrence Edward Page"],
#                 "nationality": "American",
#             },
#             "ownedCompanies": [{"name": "Google"}, {"name": "Alphabet"}],
#         }
#     },
# ]


# if __name__ == "__main__":
#     match = dedupe_against_mongo(test_records[3])
#     print(match)
