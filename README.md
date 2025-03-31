# Trood Task

## Architecture Proposal

TODO: Should be an image here

### Basic Structure: 
- Web Server (Golang/Node/Kotlin) - Just the main API entry point  
- Intent recognition Server (Most models are made for Python, so probably python) - Should be listening for Kafka/RabbitMQ events that are going to be just the user queries to be processed for user intent. 
  -  The NLP model needs to map user messages to predefined intents (spaCy, Hugging Face Transformers, BERT)
- Database (PostgresQL) - Will be used as a knowledge base for responding to users and their intents. Might need more context, but a vector database sounds like overkill, a simple intent keyword -> answer should suffice. This will also be storing the user chat to evaluate responses and evaluate human agents when they come in play.
- Queue (RabbitMQ, Kafka) - Choosing the backend for a queueing system depends on the requirements. If the priority queue is just a simple FIFO - then kafka, if perhaps the wish is to evaluate user sentiment (angrier customer likely need to be addressed first, and happier ones may not be expecting a response at all), then RabbitMQ is a better fit.