import unittest
from unittest.mock import patch, MagicMock
from cronjob import get_articles, scrape_article_content, summarize_text, send_to_queue

class TestCronJob(unittest.TestCase):

    @patch('cronjob.pika.BlockingConnection')
    def test_send_to_queue(self, mock_connection):
        mock_channel = MagicMock()
        mock_conn_instance = MagicMock()
        mock_conn_instance.channel.return_value = mock_channel
        mock_connection.return_value = mock_conn_instance

        news_data = {
            "source": "Test Source",
            "title": "Test Title",
            "url": "http://example.com",
            "summary": "Short summary"
        }

        send_to_queue(news_data)
        mock_channel.basic_publish.assert_called()
        mock_connection.assert_called()


    @patch('cronjob.newsapi')
    def test_get_articles(self, mock_newsapi):
        mock_newsapi.get_everything.return_value = {
            "articles": [{"title": "Test Article", "url": "http://example.com", "source": {"name": "Test Source"}}]
        }

        articles = get_articles("Elon Musk")
        self.assertIsInstance(articles, list)
        self.assertGreater(len(articles), 0)
        self.assertIn("title", articles[0])


    @patch('cronjob.model')
    @patch('cronjob.tokenizer')
    def test_summarize_text(self, mock_tokenizer, mock_model):
        # Mock encode and decode on tokenizer
        mock_tokenizer.encode.return_value = "encoded"
        mock_tokenizer.decode.return_value = "This is a summary."

        # Mock generate on model
        mock_model.generate.return_value = ["summary_ids"]

        summary = summarize_text("This is a long article text.")
        self.assertEqual(summary, "This is a summary.")


    @patch('cronjob.requests.get')
    def test_scrape_article_content_success(self, mock_get):
        mock_response = MagicMock()
        mock_response.status_code = 200
        mock_response.text = "<html><body><p>Paragraph 1.</p><p>Paragraph 2.</p></body></html>"
        mock_get.return_value = mock_response

        content = scrape_article_content("http://example.com")
        self.assertIn("Paragraph 1.", content)


    @patch('cronjob.requests.get', side_effect=Exception("Timeout"))
    def test_scrape_article_content_failure(self, mock_get):
        content = scrape_article_content("http://example.com")
        self.assertIn("Scraping failed", content)


    def test_pytest_runs(self):
        self.assertTrue(1 + 1 == 2)


if __name__ == '__main__':
    unittest.main()
