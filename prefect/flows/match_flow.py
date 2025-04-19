from prefect import flow
from tasks.scrape_task import generate_openai_response, parse_openai_response
from tasks.qdrant_task import search_profiles_by_json
from tasks.dedupe_task import dedupe_against_mongo
from tasks.mongo_task import (
    update_job_status,
    update_job_match_results,
    get_client_names,
    add_job_log,
)
from tasks.pdf_task import decode_file, extract_text
from tasks.notification_task import (
    publish_notification,
    JobStatus,
    JobType,
    Priority,
    Notification,
    NotificationType,
)


# for text matching
@flow(name="match-client", log_prints=True)
def match_client_flow(
    file_name: str, file_bytes: str, job_id: str, target_id: str, username: str
):
    """
    1. Decode and extract text
    2. Generate and parse LLM response
    3. Search for profile matches
    4. Update job results and status
    """

    DEDUPE_WEIGHT = 0.6
    QDRANT_WEIGHT = 0.4

    try:
        if job_id:
            update_job_status(job_id, "processing", "Client matching job started")

        file_stream = decode_file(file_bytes)
        text = extract_text(file_stream, file_name)

        names = get_client_names(target_id)

        response = generate_openai_response(text, names[0], names)
        print("[MATCHING] OpenAI response generated, parsing response...")
        add_job_log(job_id, "OpenAI response generated, parsing response...")

        profile_json = parse_openai_response(response)
        print("[MATCHING] OpenAI response parsed, searching for matches...")
        add_job_log(job_id, "OpenAI response parsed, searching for matches...")

        matches = search_profiles_by_json(profile_json, "clients")
        print("[MATCHING] Matches found, updating job results...")
        add_job_log(job_id, "Matches found, updating job results...")

        if matches:
            dedupe_match = dedupe_against_mongo(profile_json, matches)
            if dedupe_match:
                print(f"[DEDUPLICATION] Matched existing profile:\n{dedupe_match}")
                add_job_log(job_id, f"Matched existing profile:\n{dedupe_match}")
            else:
                print("[DEDUPLICATION] No matching profile found.")
                add_job_log(job_id, "No matching profile found.")

        if dedupe_match:
            for m in matches:
                if m["id"] == dedupe_match["matched_id"]:
                    qdrant_score = m["similarityScore"]
                    weighted_avg = round(
                        (DEDUPE_WEIGHT * dedupe_match["confidence"])
                        + (QDRANT_WEIGHT * qdrant_score),
                        4,
                    )
                    break
            update_job_match_results(
                job_id,
                [{"id": dedupe_match["matched_id"], "confidenceScore": weighted_avg}],
            )
        else:
            update_job_match_results(job_id, [])


        if job_id:
            update_job_status(job_id, "completed", "Client matching job completed")

        notification = Notification(
            notificationType=NotificationType.JOB,
            username=username,
            jobId=job_id,
            status=JobStatus.COMPLETED,
            type=JobType.MATCH,
            clientId=target_id,
            clientName=names[0],
            priority=Priority.LOW,
        )

        print("[MATCHING] Sending notification...")
        publish_notification(notification)
        add_job_log(job_id, "Notification sent, job completed")

        print("[MATCHING] Job completed, exiting...")

    except Exception as e:
        error_msg = f"Error while processing: {str(e)}"
        print(error_msg)

        notification = Notification(
            notificationType=NotificationType.JOB,
            username=username,
            jobId=job_id,
            status=JobStatus.FAILED,
            type=JobType.MATCH,
            clientId=target_id,
            clientName=names[0],
            priority=Priority.MEDIUM,
        )

        print("[MATCHING] Job failed, sending notification...")
        publish_notification(notification)
