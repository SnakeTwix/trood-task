from fastapi import FastAPI
from pydantic import BaseModel

from nlp import get_ai_instance, predict_intent

nlp = get_ai_instance()
app = FastAPI()

class UserQuery(BaseModel):
    query: str
    conversation_id: int

def process_user_query(query: UserQuery):
    intent, scores = predict_intent(nlp, query.query)

    print(intent)
    print(scores)
    # TODO: Send a request to the chat microservice

@app.post("/intent")
def get_query_intent(query: UserQuery):
    # Not ideal, but this is a simple implementation of just starting a processing job for intent recognition
    # With the proposed architecture this would be a subscriber for Kafka/RabbitMQ queues instead of url point and deciding how much it wants to process
    process_user_query(query)
    return query
