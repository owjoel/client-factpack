from prefect import flow
from tasks.scrape_task import generate_openai_response, parse_openai_response
from tasks.qdrant_task import search_profiles_by_json
from tasks.mongo_task import update_job_status, update_job_match_results


# for text matching
@flow(name="match-client", log_prints=True)
def match_client_flow(text: str, job_id: str):
    """
    1. Update job status to processing
    2. Extract structured data from incoming text
    3. Vectorise and search for matches in Qdrant
    """
    try:
        if job_id:
            update_job_status(job_id, "processing", f"Client matching job started")

        response = generate_openai_response(text)
        print(f"[MATCHING] OpenAI response generated, parsing response...")

        profile_json = parse_openai_response(response)
        print(f"[MATCHING] OpenAI response parsed, searching for matches...")

        matches = search_profiles_by_json(profile_json, "clients")
        print(f"[MATCHING] Matches found, updating job results...")

        if matches:
            update_job_match_results(job_id, matches)

        # TODO: Add dedupe logic here

        if job_id:
            update_job_status(job_id, "completed", f"Client matching job completed")

        print(f"[MATCHING] Job completed, exiting...")

    except Exception as e:
        error_msg = f"Error while processing: {str(e)}"
        print(error_msg)
