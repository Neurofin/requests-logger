from fastapi import FastAPI
from api import resolve
from dotenv import load_dotenv

load_dotenv()

app = FastAPI()

@app.get('/')
async def root():
    return { 'message': "Hello World!" }

@app.get("/stupid")
def getResolve():
    print("I am stupid")
    return { 'message': "I am stupid" }

app.include_router(resolve.router)
