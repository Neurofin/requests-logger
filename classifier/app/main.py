import json
import pickle
from contextlib import asynccontextmanager
from pydantic import BaseModel

import joblib
from fastapi import FastAPI

ml_models = {}


def classifier(x: str):
    serialized_string = json.dumps(x.strip(), ensure_ascii=False)
    custom_text = serialized_string
    tfidf_vectorizer = joblib.load("models/tfidf_vectorizer.pkl")
    with open("models/random_forest_classifier.pkl", "rb") as file:
        loaded_classifier = pickle.load(file)
    custom_text_tfidf = tfidf_vectorizer.transform([custom_text])
    predicted_label = loaded_classifier.predict(custom_text_tfidf)[0]
    # Get the confidence score for the prediction
    confidence_scores = loaded_classifier.predict_proba(custom_text_tfidf)

    # Find the index of the predicted label in the classes
    class_labels = loaded_classifier.classes_
    predicted_label_index = list(class_labels).index(predicted_label)

    # Get the confidence score for the predicted label
    print(confidence_scores)
    confidence_score = confidence_scores[0][predicted_label_index]
    return predicted_label, confidence_score


@asynccontextmanager
async def lifespan(app: FastAPI):
    print("*" * 10, "LOADING MODELS", "*" * 10)
    ml_models["rf_classifier"] = classifier
    print("*" * 10, "LOADING MODELS COMPLETED", "*" * 10)
    yield
    print("*" * 10, "SHUTTING DOWN", "*" * 10)
    ml_models.clear()


app = FastAPI(lifespan=lifespan)

class Item(BaseModel):
    text: str

@app.post("/classify")
async def predict(item: Item):
    classname, confidence = ml_models["rf_classifier"](x=item.text)
    return {"data": [{"Name": classname, "Score": confidence}]}
