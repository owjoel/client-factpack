from pydantic import BaseModel

class ClientArticle(BaseModel):
    source: str
    title: str
    url: str
    summary: str