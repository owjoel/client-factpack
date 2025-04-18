from flows.search_client_articles import update_clients
from dotenv import load_dotenv

load_dotenv()


if __name__ == "__main__":
    update_clients.serve(name="news-client-update", cron="0 9 * * *")
