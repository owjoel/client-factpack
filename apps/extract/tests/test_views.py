import pytest
from flask import Flask
from apps.extract.extract.models import Client, Profile, Metadata, Residence, NetWorth, Contact
from apps.extract.extract import create_app

@pytest.fixture
def client(mocker):
    # Create the full fake completion call chain
    mock_completions = mocker.Mock()
    mock_completions.parse.return_value = Client(
        profile=Profile(
            name="Jane",
            age=30,
            nationality="SG",
            currentResidence=Residence(city="Singapore", country="SG"),
            netWorth=NetWorth(estimatedValue=500000, currency="SGD", source="source"),
            industries=["Tech"],
            occupations=["Engineer"],
            socials=[],
            contact=Contact(phone="123456", workAddress="123 Orchard Rd")
        ),
        investments=[],
        associates=[],
        metadata=Metadata(sources=["mock"])
    )

    mock_client = mocker.Mock()
    mock_client.beta.chat.completions = mock_completions

    # Patch __init__ to do nothing
    mocker.patch("apps.extract.extract.service.OpenAIClientParser.__init__", lambda self: None)
    # Patch the class to return a plain object
    mocker.patch("apps.extract.extract.service.OpenAIClientParser", autospec=True)

    app = create_app()
    app.config["TESTING"] = True

    # Manually attach mock client to parser instance
    from apps.extract.extract.service import OpenAIClientParser
    parser_instance = OpenAIClientParser()
    parser_instance.client = mock_client
    app.parser = parser_instance  # manually attach to the Flask app

    with app.test_client() as test_client:
        yield test_client, parser_instance


def test_health_endpoint(client):
    test_client, _ = client
    response = test_client.get("/api/v1/health")
    assert response.status_code == 200
    assert response.get_json() == {"message": "Healthy"}


def test_extract_info_endpoint(client):
    test_client, mock_parser = client

    # Setup fake client object
    fake_client = Client(
        profile=Profile(
            name="Jane",
            age=30,
            nationality="SG",
            currentResidence=Residence(city="Singapore", country="SG"),
            netWorth=NetWorth(estimatedValue=500000, currency="SGD", source="source"),
            industries=["Tech"],
            occupations=["Engineer"],
            socials=[],
            contact=Contact(phone="123456", workAddress="123 Orchard Rd")
        ),
        investments=[],
        associates=[],
        metadata=Metadata(sources=["mock"])
    )

    mock_parser.extract_client_info.return_value = fake_client

    response = test_client.post("/api/v1/extract/parseInfo", json={"text": "Jane is a software engineer..."})
    assert response.status_code == 200
    assert response.get_json()["profile"]["name"] == "Jane"
