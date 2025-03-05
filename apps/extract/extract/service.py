from openai import OpenAI
from .models import Client
from abc import ABC, abstractmethod

import os

class ClientParser(ABC):
    """Port for LLM models"""

    @abstractmethod
    def extract_client_info(self, text: str) -> Client:
        pass



class OpenAIClientParser(ClientParser):
    """Adapter for OpenAI"""

    def __init__(self):
        self.client = OpenAI()
        self.model=os.getenv("OPENAI_MODEL")
        self.assistant_id=os.getenv("OPENAI_ASSISTANT_ID")
    
    def extract_client_info(self, text: str) -> Client:
        print("extracting client info")
        completion = self.client.beta.chat.completions.parse(
            model=self.model,
            messages=[
                {"role": "system", "content": "Extract relevant infomation out of the user raw text in a JSON object."},
                {"role": "user", "content": text}
            ],
            response_format=Client
        )
        client_data = completion.choices[0].message.parsed
        return client_data