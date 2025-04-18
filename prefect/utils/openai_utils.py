import json
import os
import re
from openai import OpenAIError, OpenAI
from dotenv import load_dotenv

load_dotenv()
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")

client = OpenAI(api_key=OPENAI_API_KEY)


def query_gpt4o(text: str) -> str:
    """Query GPT-4o with a structured schema-based prompt."""
    schema_path = os.path.join(os.path.dirname(__file__), "schema.json")

    if not os.path.exists(schema_path):
        raise FileNotFoundError(f"Schema file not found at: {schema_path}")

    try:
        with open(schema_path, "r") as f:
            schema = json.load(f)
    except json.JSONDecodeError as e:
        raise ValueError(f"Invalid JSON in schema file: {e}")

    try:
        response = client.chat.completions.create(
            model="gpt-4o",
            messages=[{"role": "user", "content": text}],
            response_format={"type": "json_schema", "json_schema": schema},
        )
        raw_content = response.choices[0].message.content
        return clean_response(raw_content)
    except OpenAIError as e:
        raise RuntimeError(f"OpenAI API request failed: {e}")
    except (KeyError, AttributeError, IndexError) as e:
        raise ValueError(f"Unexpected response structure: {e}")


def clean_response(response: str) -> str:
    """Clean GPT response for safe JSON parsing."""
    if not response or not isinstance(response, str):
        raise ValueError("Empty or invalid response from model.")

    cleaned = re.sub(r"^```(?:json)?\n?|\n?```$", "", response.strip())
    cleaned = re.sub(r",\s*([}\]])", r"\1", cleaned)

    try:
        json.loads(cleaned)
    except json.JSONDecodeError as e:
        raise ValueError(f"Cleaned response is not valid JSON: {e}")

    return cleaned
