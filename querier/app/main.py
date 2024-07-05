from fastapi import FastAPI
from api import resolve
from api import classify
from dotenv import load_dotenv

load_dotenv()

app = FastAPI()

@app.get('/querier')
async def root():
    return { 'message': "Hello World!" }


app.include_router(resolve.router)
app.include_router(classify.router)
