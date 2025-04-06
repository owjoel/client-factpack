from prefect import flow
from tasks.qdrant_task import upsert_text_to_qdrant
from tasks.scrape_task import (
    get_cleaned_target,
    get_wikipedia_text,
    generate_openai_response,
    parse_openai_response,
    save_files,
    update_client_profile,
)
from tasks.mongo_task import update_job_status, add_job_log


@flow(name="scrape-client", log_prints=True)
def scrape_client_flow(job_id: str, target: str, client_id: str):
    try:
        if job_id:
            update_job_status(
                job_id, "processing", f"Client scraping job started for {target}"
            )

        target_clean = get_cleaned_target(target)
        wiki_text = get_wikipedia_text(target)
        add_job_log(job_id, f"[{target}] Wikipedia text retrieved")
        print(f"[{target}] Wikipedia text retrieved, calling OpenAI...")

        response = generate_openai_response(target, wiki_text)
        print(f"[{target}] OpenAI response generated, parsing...")

        profile_json = parse_openai_response(response)
        print(f"[{target}] OpenAI response parsed, saving files...")

        save_files(target_clean, wiki_text, profile_json)
        print(f"[{target}] Files saved, updating client profile...")

        # inserted_id = insert_into_mongo(profile_json, target_clean)
        update_client_profile(client_id, profile_json)
        add_job_log(job_id, f"[{target}] Profile saved")

        print(f"[{target}] Client profile updated, upserting into Qdrant...")

        upsert_text_to_qdrant("clients", profile_json, client_id)
        add_job_log(job_id, f"[{target}] Upserted into Qdrant")
        print(f"[{target}] Upserted into Qdrant, job complete")

        print(f"Job complete for {target} → Mongo ID: {client_id}")
        if job_id:
            update_job_status(
                job_id, "completed", f"Client scraping job completed for {target}"
            )

    except Exception as e:
        error_msg = f"Error while processing {target}: {str(e)}"
        print(error_msg)
        if job_id:
            update_job_status(job_id, "failed", error_msg)
