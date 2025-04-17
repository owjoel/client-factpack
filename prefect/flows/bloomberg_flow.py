from prefect import flow
from datetime import datetime
from tasks.bloomberg_task import fetch_webpage, parse_billionaires_data
from utils.file_utils import save_to_csv


@flow(name="bloomberg-scraper", log_prints=True)
def scrape_bloomberg_flow():
    url = "https://www.bloomberg.com/billionaires/"
    current_date = datetime.now().strftime("%d%m%Y")
    folder_name = "bloomberg"
    csv_file = f"{current_date}"

    html_content = fetch_webpage(url)
    data = parse_billionaires_data(html_content)

    if not data:
        print("No data parsed.")
        return

    save_to_csv(data, folder_name, csv_file)
    print(f"[LOG] Saved to {folder_name}/{csv_file}")
