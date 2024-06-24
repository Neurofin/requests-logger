from contextlib import asynccontextmanager
from pydantic import BaseModel

import joblib
import numpy as np
from fastapi import FastAPI

ml_models = {}


def classifier(text: str):
    best_pipeline = joblib.load("models/best_pipeline.pkl")
    predicted_probabilities = best_pipeline.predict_proba([text])[0]
    predicted_label = best_pipeline.predict([text])[0]
    confidence_score = np.max(predicted_probabilities)
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
    classname, confidence = ml_models["rf_classifier"](text=item.text)
    return {"data": [{"Name": classname, "Score": confidence}]}
