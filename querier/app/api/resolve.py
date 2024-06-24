import fastapi
from pydantic import BaseModel

import openai
from langchain_community.chat_models.openai import ChatOpenAI
from llama_index.core import VectorStoreIndex
from llama_index.readers.string_iterable import StringIterableReader

from llama_index.llms.openai import OpenAI

import json
import os
from dotenv import load_dotenv

load_dotenv()

from utils.extractText import downloadAndReturnTexts

router = fastapi.APIRouter()

class ContextDocument(BaseModel):
    docType: str
    docPath: str

class Query(BaseModel):
    contextDocuments: list[ContextDocument]
    prompt: str


def textjson_to_json(json_text):
  json_start_index = json_text.find('{')
  json_end_index = json_text.rfind('}')

  json_text = json_text[json_start_index:json_end_index+1]
  jsonData = json.loads(json_text)
  return jsonData

@router.post("/resolve")
def resolve(query: Query):
    texts = downloadAndReturnTexts(query.contextDocuments)

    reader = StringIterableReader()
    documents = reader.load_data(texts=texts)
    index = VectorStoreIndex.from_documents(documents=documents)

    os.environ["OPENAI_API_KEY"] = os.getenv("OPENAI_API_KEY")
    openai.api_key = os.environ["OPENAI_API_KEY"]

    llm = OpenAI(temperature=0, model="gpt-3.5-turbo-16k")
    query_engine = index.as_query_engine(
       chat_mode="context", llm=llm
    )

    prompt = query.prompt

    result = {}
    response = query_engine.query(prompt)
    # print(response.response)
    try:
        json_format = textjson_to_json(response.response)
        result = json_format
    except Exception as err:
        print(
            "Error in deserialzing GPT response. Cause: ",
            prompt,
            "Response: ",
            response.response,
            err,
        )
    return { 'message': "Success", 'data': result }



