from bs4 import BeautifulSoup
import cloudscraper
from prefect import task


@task
def fetch_webpage(url):
    """
    Fetch the webpage content using cloudscraper.
    """
    scraper = cloudscraper.create_scraper(
        browser={
            "custom": "Chrome/120.0.0.0",
            "platform": "Windows",
        }
    )

    response = scraper.get(url)

    if response.status_code != 200:
        raise ValueError(
            f"Failed to retrieve the webpage. Status code: {response.status_code}"
        )

    return response.text


@task
def parse_billionaires_data(html_content: str) -> list[dict]:
    soup = BeautifulSoup(html_content, "html.parser")
    rows = soup.find_all("div", class_="table-row")

    data = []
    for row in rows:
        cells = row.find_all("div", class_="table-cell")
        if len(cells) < 7:
            continue

        rank = cells[0].get_text(strip=True)
        name_cell = cells[1]
        name = name_cell.find("a").get_text(strip=True) if name_cell.find("a") else ""
        net_worth = cells[2].get_text(strip=True)
        lcd = cells[3].get_text(strip=True)
        ycd = cells[4].get_text(strip=True)
        country = cells[5].get_text(strip=True)
        industry = cells[6].get_text(strip=True)

        data.append(
            {
                "Rank": rank,
                "Name": name,
                "Net Worth": net_worth,
                "LCD": lcd,
                "YCD": ycd,
                "Country": country,
                "Industry": industry,
            }
        )

    return data
