from prefect import flow
from tasks.scrape_task import generate_openai_response, parse_openai_response
from tasks.qdrant_task import search_profiles_by_json

# from tasks.mongo_task import update_job_status, update_job_match_results
from tasks.dedupe_task import dedupe_against_mongo


# for text matching
@flow(name="match-client", log_prints=True)
def match_client_flow(text: str, job_id: str):
    """
    1. Update job status to processing
    2. Extract structured data from incoming text
    3. Vectorise and search for matches in Qdrant
    4. Match the text against the top 3 matches returned by Qdrant
    """
    try:
        # if job_id:
        # update_job_status(job_id, "processing", "Client matching job started")

        response = generate_openai_response(text)
        print("[MATCHING] OpenAI response generated, parsing response...")

        profile_json = parse_openai_response(response)
        print("[MATCHING] OpenAI response parsed, searching for matches...")

        matches = search_profiles_by_json(profile_json, "clients")
        print("[MATCHING] Matches found, updating job results...")

        if matches:
            dedupe_match = dedupe_against_mongo(profile_json, matches)
            if dedupe_match:
                print(f"[DEDUPLICATION] Matched existing profile:\n{dedupe_match}")
            else:
                print("[DEDUPLICATION] No matching profile found.")

            # if dedupe_match:
            #     update_job_match_results(job_id, [dedupe_match])
        # else:
        #     update_job_match_results(job_id, [])

        # if job_id:
        #     update_job_status(job_id, "completed", "Client matching job completed")

        print("[MATCHING] Job completed, exiting...")

    except Exception as e:
        error_msg = f"Error while processing: {str(e)}"
        print(error_msg)
