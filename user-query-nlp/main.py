import requests
from fastapi import FastAPI
from pydantic import BaseModel

from nlp import get_ai_instance, intent_score

nlp = get_ai_instance()
app = FastAPI()


class UserQuery(BaseModel):
    query: str
    conversation_id: int


def process_user_query(query: UserQuery):
    intent, score = intent_score(nlp, query.query)

    print(intent)
    print(score)

    # Handle whatever error there could be, I don't have to do that
    _ = requests.post('http://customer-support-api:1323/nlp/intent', json={
        "confidence": score,
        "intent": intent,
        "conversation_id": query.conversation_id,
    })


@app.post("/intent")
def get_query_intent(query: UserQuery):
    # Not ideal, but this is a simple implementation of just starting a processing job for intent recognition
    # With the proposed architecture this would be a subscriber for Kafka/RabbitMQ queues instead of url point and deciding how much it wants to process
    process_user_query(query)
    return query
