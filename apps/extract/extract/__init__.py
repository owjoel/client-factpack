from typing import Optional, cast
from flask import Flask
from .views import views
from .service import ClientParser, OpenAIClientParser

import os

class App(Flask):
    parser: ClientParser

def create_app():
    app = App(__name__)
    app.config.from_prefixed_env()
    app.register_blueprint(views, url_prefix="/api/v1")
    
    # llm dependency injection
    app.parser = OpenAIClientParser()

    return app