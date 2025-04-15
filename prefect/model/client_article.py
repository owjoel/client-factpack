from pydantic import BaseModel

class Sentiment(BaseModel):
    label: str
    score: float

class ClientArticle(BaseModel):
    source: str
    title: str
    url: str
    summary: str
    sentiment: Sentiment