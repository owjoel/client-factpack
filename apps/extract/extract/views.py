from flask import Blueprint, jsonify, request, current_app
from .service import *

# Register blueprints, which are Flask api groups
# views blueprint is registered in main app
views = Blueprint("views", __name__)
extract = Blueprint("matching", __name__)
views.register_blueprint(extract, url_prefix="/extract")

# health is the basic healthcheck endpoint
@views.get("/health")
def health():
    return jsonify({"message": "Healthy"})

# extract_info calls openai service function
@extract.post("/parseInfo")
def extract_info():
    data = request.json
    text: str = data.get("text")
    client = current_app.parser.extract_client_info(text)
    return jsonify(client.model_dump())