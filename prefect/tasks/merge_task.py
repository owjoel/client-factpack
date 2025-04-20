import json
import re
from rapidfuzz import fuzz
from copy import deepcopy
from prefect import task
from utils import (
    openai_utils,
    prompt_utils,
)

@task
def normalize_name(name: str) -> str:
    name = name.lower()

    name = re.sub(
        r"\b(inc|incorporated|ltd|limited|llc|corp|corporation|plc|co|company)\b",
        "",
        name,
    )
    name = re.sub(r"\.com\b|\.(org|net|co)\b", "", name)
    name = re.sub(r"[^\w\s]", "", name)
    name = re.sub(r"\s+", " ", name).strip()
    return name


@task
def fuzzy_match(s1, s2, threshold=85):
    s1_clean = normalize_name(s1)
    s2_clean = normalize_name(s2)
    print(s1_clean, s2_clean)
    return fuzz.ratio(s1_clean, s2_clean) >= threshold


@task
def dedup_list(existing, new_list, threshold=90):
    merged = existing.copy()
    for item in new_list:
        if not any(fuzzy_match(item, exist) for exist in existing):
            merged.append(item)
    return merged


@task
def dedup_object_list(existing, new_list, key="name", threshold=90):
    merged = deepcopy(existing)
    for new_item in new_list:
        if not any(
            fuzzy_match(new_item[key], old_item[key])
            for old_item in existing
            if key in old_item and key in new_item
        ):
            merged.append(new_item)
    return merged


@task
def review_with_openai(existing: dict, incoming: dict, merged: dict) -> dict:
    """
    Ask GPT-4o to review and clean up a merged profile using schema-constrained output.
    """
    prompt = prompt_utils.build_merge_review_prompt(existing, incoming, merged)
    print("Reviewing merge with OpenAI...")
    response = openai_utils.query_gpt4o(prompt)

    try:
        return json.loads(response)
    except json.JSONDecodeError as e:
        raise ValueError(f"Failed to parse GPT output as JSON: {e}")


@task
def merge_profiles(existing, incoming):
    merged = deepcopy(existing)

    p_existing = merged["profile"]
    p_incoming = incoming["profile"]

    # Merge string lists
    for field in ["names", "industries", "occupations", "pastOccupations"]:
        p_existing[field] = dedup_list(
            p_existing.get(field, []), p_incoming.get(field, [])
        )

    # Merge currentResidence, netWorth, etc. (prefer existing if non-empty, else use incoming)
    for field in ["gender", "dateOfBirth", "description", "nationality"]:
        if not p_existing.get(field) and p_incoming.get(field):
            p_existing[field] = p_incoming[field]

    if not p_existing.get("currentResidence") and p_incoming.get("currentResidence"):
        p_existing["currentResidence"] = p_incoming["currentResidence"]

    if not p_existing.get("netWorth") and p_incoming.get("netWorth"):
        p_existing["netWorth"] = p_incoming["netWorth"]

    # Merge socials
    p_existing["socials"] = dedup_object_list(
        p_existing.get("socials", []), p_incoming.get("socials", [])
    )

    # Merge career timeline
    p_existing["careerTimeline"] = dedup_object_list(
        p_existing.get("careerTimeline", []),
        p_incoming.get("careerTimeline", []),
        key="event",
    )

    # Merge ownedCompanies, subsidiaries inside each
    def merge_owned_companies(existing_comps, incoming_comps):
        merged = deepcopy(existing_comps)
        for new_company in incoming_comps:
            if not any(
                fuzzy_match(new_company["name"], ex["name"]) for ex in existing_comps
            ):
                merged.append(new_company)
            else:
                for ex in merged:
                    if fuzzy_match(new_company["name"], ex["name"]):
                        ex["subsidiaries"] = dedup_object_list(
                            ex.get("subsidiaries", []),
                            new_company.get("subsidiaries", []),
                            key="name",
                        )
        return merged

    merged["ownedCompanies"] = merge_owned_companies(
        merged.get("ownedCompanies", []),
        incoming.get("ownedCompanies", []),
    )

    # Merge investments
    merged["investments"] = dedup_object_list(
        merged.get("investments", []),
        incoming.get("investments", []),
        key="name",
    )

    # Merge family and associates
    for key in ["family", "associates"]:
        merged[key] = dedup_object_list(
            merged.get(key, []), incoming.get(key, []), key="name"
        )

    # Merge sources: deduplicate by source name
    merged["sources"] = dedup_object_list(
        merged.get("sources", []),
        incoming.get("sources", []),
        key="source",
    )

    return merged
