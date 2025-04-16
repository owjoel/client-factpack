from prefect import task
import torch
import logging

from model.client_article import Sentiment
from transformers import AutoModelForSequenceClassification, AutoTokenizer, pipeline

logging.basicConfig(level=logging.INFO)

finbert_model = AutoModelForSequenceClassification.from_pretrained("ProsusAI/finbert")
finbert_tokenizer = AutoTokenizer.from_pretrained("ProsusAI/finbert")

device = 0 if torch.cuda.is_available() else -1
sentiment_analyzer = pipeline(
    "sentiment-analysis",
    model=finbert_model,
    tokenizer=finbert_tokenizer,
    device=device,
)

@task
def analyze_sentiment(text):
    """
    Analyze the sentiment of the article text using FinBERT
    """
    if not text:
        return Sentiment(label="neutral", score=0.5)

    try:
        # For long texts, analyze chunks and take the average
        if len(text) > 500:
            max_chars = 5000
            chunk_size = 500
            chunks = [text[i : i + 500] for i in range(0, min(len(text), max_chars), chunk_size)]
            results = []

            for chunk in chunks:
                result = sentiment_analyzer(chunk)[0]
                results.append(result)

            # Calculate the average sentiment
            sentiment_map = {"positive": 1, "negative": -1, "neutral": 0}
            weighted_sentiment = sum(
                sentiment_map[r["label"]] * r["score"] for r in results
            ) / len(results)

            # Map back to label
            if weighted_sentiment > 0.2:
                final_label = "positive"
            elif weighted_sentiment < -0.2:
                final_label = "negative"
            else:
                final_label = "neutral"

            return Sentiment(label=final_label, score=abs(weighted_sentiment))
        else:
            result = sentiment_analyzer(text)[0]
            return Sentiment(label=result['label'], score=result['score'])
    except Exception as e:
        logging.error("Error analyzing sentiment", exc_info=True)
        return Sentiment(label="neutral", score=0.5)

def getPriority(label: str):
    match label:
        case "positive":
            return "low"
        case "neutral":
            return "medium"
        case "negative":
            return "high"