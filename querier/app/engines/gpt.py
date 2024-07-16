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

    response = query_engine.query(prompt)
    json_format = textjson_to_json(response.response)
    result = json_format
    return result

from llama_index.core import Settings, Document, SummaryIndex
from llama_index.core.node_parser import SentenceSplitter
from llama_index.core.tools import QueryEngineTool, ToolMetadata
from llama_index.agent.openai import OpenAIAgent
from llama_index.core.objects import ObjectIndex
import time

def gptAgentsSetup(contextDocs: dict[str, list[str]], version:str):

    model = f"gpt-{version}"

    # Define global settings for LLM and embeddings
    Settings.llm = OpenAI(temperature=0, model=model)

    node_parser = SentenceSplitter()
    
    # # Build agents for each document
    all_tools = []

    folder = f"{int(time.time())}"
    for docType, contents in contextDocs.items():
        
        # Initialize SimpleDirectoryReader correctly
        documents = []
        for content in contents:
            documents.append(Document(text=content))

        nodes = node_parser.get_nodes_from_documents(documents)
        # all_nodes.extend(nodes)

        # Build vector index
        vector_index = VectorStoreIndex(nodes)
        vector_index.storage_context.persist(persist_dir=f"./{folder}/{docType}")

        # Build summary index
        summary_index = SummaryIndex(nodes)

        # Define query engines
        vector_query_engine = vector_index.as_query_engine(llm=Settings.llm)
        summary_query_engine = summary_index.as_query_engine(llm=Settings.llm)

        query_engine_tools = [
            QueryEngineTool(
                query_engine=vector_query_engine,
                metadata=ToolMetadata(
                    name=f"vector_tool_{docType}",
                    description=f"Useful for questions related to {docType}."
                ),
            ),
            QueryEngineTool(
                query_engine=summary_query_engine,
                metadata=ToolMetadata(
                    name=f"summary_tool_{docType}",
                    description=f"Useful for summarizing {docType}."
                ),
            ),
        ]

        # Build agent
        function_llm = OpenAI(model=model)
        agent = OpenAIAgent.from_tools(
            query_engine_tools,
            llm=function_llm,
            verbose=True,
            system_prompt=f"You are a specialized agent designed to answer queries about {docType}. You must ALWAYS use at least one of the tools provided when answering a question; do NOT rely on prior knowledge."
        )

        wiki_summary = f"This content contains information about {docType}. Use this tool if you want to answer any questions about {docType}.\n"
        
        doc_tool = QueryEngineTool(
            query_engine=agent,
            metadata=ToolMetadata(
                name=f"tool_{docType}",
                description=wiki_summary,
            ),
        )
        all_tools.append(doc_tool)

    # Define an "object" index and retriever over these tools
    obj_index = ObjectIndex.from_objects(
        all_tools,
        index_cls=VectorStoreIndex,
    )

    # Build top-level OpenAI agent
    top_agent = OpenAIAgent.from_tools(
        tool_retriever=obj_index.as_retriever(similarity_top_k=3),
        system_prompt="You are an agent designed to answer queries about a set of documents. Please always use the tools provided to answer a question. Do not rely on prior knowledge.",
        verbose=True,
    )

    return { "agent": top_agent, "folder": folder }

from utils.cleanup import cleanupFolder

def gptAgentQuerier(prompt: str, contextDocs: dict[str, list[str]], version: str):
    topAgentSetup = gptAgentsSetup(contextDocs=contextDocs, version=version)
    topAgent = topAgentSetup["agent"]
    response = topAgent.query(prompt)
    json_format = textjson_to_json(response.response)
    result = json_format
    try:
        cleanupFolder(topAgentSetup["folder"])
    except Exception as err:
        print(err)
    return result
