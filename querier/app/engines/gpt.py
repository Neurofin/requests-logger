import openai
from llama_index.core import VectorStoreIndex
from llama_index.readers.string_iterable import StringIterableReader

from llama_index.llms.openai import OpenAI
import os

from utils.parsers import textjson_to_json

def gptQuerier(prompt: str, contextTexts: list[str], version: str):
    reader = StringIterableReader()
    documents = reader.load_data(texts=contextTexts)
    index = VectorStoreIndex.from_documents(documents=documents)

    os.environ["OPENAI_API_KEY"] = os.getenv("OPENAI_API_KEY")
    openai.api_key = os.environ["OPENAI_API_KEY"]

    llm = OpenAI(temperature=0, model=f"gpt-{version}")
    query_engine = index.as_query_engine(
       chat_mode="context", llm=llm
    )

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
    return result