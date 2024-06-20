import fastapi

from pydantic import BaseModel

from langchain_community.chat_models.openai import ChatOpenAI
from llama_index.core import VectorStoreIndex
from llama_index.readers.string_iterable import StringIterableReader

import json

router = fastapi.APIRouter()

class Query(BaseModel):
    contextDocumentPaths: list[str]
    prompts: list[str]


@router.post("/resolve")
def resolve(query: Query):
    # TODO: Download context files
    reader = StringIterableReader()
    documents = reader.load_data(texts=texts)
    index = VectorStoreIndex.from_documents(documents=documents)
    query_engine = index.as_query_engine(
        llm=ChatOpenAI(temperature=0, model="gpt-3.5-turbo-16k")
    )

    prompts = query.prompts
    results = []
    for prompt in prompts:
        response = query_engine.query(prompt)
        # print(response.response)
        try:
            results.append(json.loads(response.response))
        except Exception as err:
            print(
                "Error in deserialzing GPT response. Cause: ",
                prompt,
                "Response: ",
                response.response,
                err,
            )
    return { 'message': "Success", 'data': results }
