# Project Setup

## Prerequisites

Make sure you have Python installed on your system.

## Installation

1. **Create a `.env` file** in the project root and add the following environment variables:

   ```ini
   OPENAI_API_KEY=<your_openai_api_key>
   OPENAI_ASSISTANT_ID=<your_openai_assistant_id>
   OPENAI_MODEL=<your_openai_model>
   ```
2. **Create python environment**
   
   ```sh
   python -m venv env

   # windows
   env/Scripts/activate

   # mac
   source env/bin/activate
   ```

3. **Install dependencies** by running:

   ```sh
   pip install -r requirements.txt
   ```

3. **Run the application** using Flask:

   ```sh
   flask --app extract run
   ```

## Usage

After running the Flask server, you can interact with the application via the provided API endpoints, namely:

1. `/api/v1/extract/parseInfo` </br>
POST method with text from scraped document