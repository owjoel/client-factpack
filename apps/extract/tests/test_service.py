from apps.extract.extract.models import Client, Profile, Metadata, Residence, NetWorth, Contact

def test_extract_client_info_mocked(monkeypatch):
    class MockClient:
        def __init__(self):
            self.beta = self
            self.chat = self
            self.completions = self

        def parse(self, **kwargs):
            class Message:
                parsed = Client(
                    profile=Profile(
                        name="Jane",
                        age=30,
                        nationality="UK",
                        currentResidence=Residence(city="London", country="UK"),
                        netWorth=NetWorth(estimatedValue=500000, currency="GBP", source="mock"),
                        industries=["Tech"],
                        occupations=["Engineer"],
                        socials=[],
                        contact=Contact(workAddress="123 Tech St", phone="555-5555")
                    ),
                    investments=[],
                    associates=[],
                    metadata=Metadata(sources=["mock"])
                )
            class Choice:
                message = Message()
            class Completion:
                choices = [Choice()]
            return Completion()

    monkeypatch.setattr("apps.extract.extract.service.OpenAI", lambda: MockClient())

    from apps.extract.extract.service import OpenAIClientParser
    parser = OpenAIClientParser()
    client_data = parser.extract_client_info("Jane is from London")

    assert client_data.profile.name == "Jane"

