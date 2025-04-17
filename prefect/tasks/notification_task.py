from pydantic import BaseModel
from enum import Enum
import pika
from dotenv import load_dotenv
from prefect import task
import os
from typing import Optional


load_dotenv()
RABBITMQ_URL = os.getenv("RABBITMQ_URL")


class NotificationType(str, Enum):
    JOB = "job"
    CLIENT = "client"


class JobStatus(str, Enum):
    PENDING = "pending"
    PROCESSING = "processing"
    COMPLETED = "completed"
    FAILED = "failed"


class JobType(str, Enum):
    SCRAPE = "scrape"
    MATCH = "match"


class Priority(str, Enum):
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"


class Notification(BaseModel):
    notificationType: NotificationType
    username: Optional[str] = None
    jobId: Optional[str] = None
    status: Optional[JobStatus] = None
    type: Optional[JobType] = None
    clientId: Optional[str] = None
    clientName: Optional[str] = None
    priority: Optional[Priority] = None


@task
def publish_notification(notification: Notification):
    connection = pika.BlockingConnection(pika.URLParameters(RABBITMQ_URL))
    channel = connection.channel()
    channel.queue_declare(queue="notifications", durable=True)
    payload = notification.json()

    channel.basic_publish(
        exchange="",
        routing_key="notifications",
        body=payload,
        properties=pika.BasicProperties(delivery_mode=2),
    )

    channel.close()
