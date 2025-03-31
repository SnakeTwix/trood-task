import json
import random
from typing import Callable

import spacy
from spacy import Language
from spacy.tokens import Doc
from spacy.training import Example


def get_ai_instance() -> Language:
    nlp = spacy.blank("en")
    textcat = nlp.add_pipe("textcat")

    train_ai(nlp, textcat)

    return nlp


def train_ai(nlp: Language, textcat: Callable[[Doc], Doc]):
    converted_examples = []

    with open("training-data.json") as training_data_file:
        training_data = json.load(training_data_file)

        for intent, data in training_data.items():
            textcat.add_label(intent)
            for text_example in data:
                doc = nlp.make_doc(text_example['value'])
                example = Example.from_dict(doc, text_example['spacyExample'])
                converted_examples.append(example)

    nlp.initialize()
    optimizer = nlp.create_optimizer()

    # Training loop
    for epoch in range(10):
        random.shuffle(converted_examples)
        losses = {}
        for example in converted_examples:
            nlp.update([example], drop=0.2, losses=losses, sgd=optimizer)


def predict_intent(nlp: Language, text: str):
    doc = nlp(text)
    intents = doc.cats
    return max(intents, key=intents.get), intents
