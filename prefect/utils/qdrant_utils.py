from qdrant_client import QdrantClient
from dotenv import load_dotenv
import os
from qdrant_client.http.models import VectorParams, Distance

load_dotenv()
QDRANT_URL = os.getenv("QDRANT_URL")
QDRANT_API_KEY = os.getenv("QDRANT_API_KEY")

qdrant_client = QdrantClient(
    url=QDRANT_URL,
    api_key=QDRANT_API_KEY,
)


def create_qdrant_collection(collection_name: str, vector_size: int, distance: str):
    qdrant_client.create_collection(
        collection_name=collection_name,
        vectors_config=VectorParams(size=vector_size, distance=distance),
    )


if __name__ == "__main__":
    create_qdrant_collection(
        collection_name="clients",
        vector_size=1024,
        distance=Distance.DOT,
    )
