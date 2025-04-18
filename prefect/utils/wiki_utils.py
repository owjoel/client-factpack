import requests
import re
from bs4 import BeautifulSoup


def get_wikipedia_url(query: str) -> str:
    """Search Wikipedia for a person and return the top result URL."""
    if not query or not query.strip():
        raise ValueError("Query cannot be empty")

    query_params = {
        "action": "query",
        "list": "search",
        "srsearch": query,
        "format": "json",
    }

    try:
        url = "https://en.wikipedia.org/w/api.php?" + "&".join(
            [f"{k}={v}" for k, v in query_params.items()]
        )
        response = requests.get(url, timeout=10)
        response.raise_for_status()
        data = response.json()

        if (
            "query" in data
            and "search" in data["query"]
            and len(data["query"]["search"]) > 0
        ):
            top_result = data["query"]["search"][0]
            page_title = top_result["title"].replace(" ", "_")
            return f"https://en.wikipedia.org/wiki/{page_title}"
        else:
            raise LookupError(f"No Wikipedia page found for query: '{query}'")

    except requests.RequestException as e:
        raise ConnectionError(f"Failed to fetch Wikipedia search results: {e}")


def extract_clean_wikipedia_page(page_url: str) -> str:
    """Extract clean text from a Wikipedia page."""
    try:
        response = requests.get(page_url, timeout=10)
        response.raise_for_status()
    except requests.RequestException as e:
        raise ConnectionError(f"Failed to fetch Wikipedia page: {e}")

    soup = BeautifulSoup(response.content, "html.parser")
    content_div = soup.find(id="mw-content-text")
    if not content_div:
        raise ValueError("Could not locate main content on the Wikipedia page.")

    extracted_text = []
    stop_sections = {"See also", "Notes", "References", "Works cited"}

    for element in content_div.find_all(["h1", "h2", "h3", "h4", "h5", "h6", "p"]):
        if (
            element.name in ["h2", "h3"]
            and element.get_text(strip=True) in stop_sections
        ):
            break

        text = element.get_text(separator=" ", strip=True)

        text = re.sub(r"\[.*?\]", "", text)
        text = re.sub(r"\[\s*citation needed\s*\]", "", text, flags=re.IGNORECASE)

        if text:
            extracted_text.append(text)

    clean_text = "\n".join(extracted_text).strip()
    if not clean_text:
        raise ValueError("Wikipedia content extracted is empty.")

    return clean_text


if __name__ == "__main__":
    url = get_wikipedia_url("Sergey Brin")
    print(url)
    text = extract_clean_wikipedia_page(url)
    print(text)
