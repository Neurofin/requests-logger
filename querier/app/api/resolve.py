import fastapi
from pydantic import BaseModel

from dotenv import load_dotenv

load_dotenv()

from utils.extractText import getContextDocumentsMapping

from engines.gpt import gptQuerier
from engines.gpt import gptAgentQuerier
from engines.gemini import geminiQuerier

router = fastapi.APIRouter()

class ContextDocument(BaseModel):
    docType: str
    docPath: str

class Query(BaseModel):
    contextDocuments: list[ContextDocument]
    prompt: str
    #gptQuerier-4o
    #gptAgents-4o
    #geminiQuerier-1.5-flash-001
    engine: str =  "gptAgents-4o"


@router.post("/resolve")
def resolve(query: Query):
    documents = getContextDocumentsMapping(query.contextDocuments)

    engine = query.engine  
    [engine, version] = query.engine.split("-", 1)

    result = {}
    if engine == 'gptQuerier':
        texts = []
        for docType, contents in documents.items():
            texts.append(f"{docType}=={'\n'.join(contents)}")
        result = gptQuerier(prompt=query.prompt, contextTexts=texts, version=version)
    if engine == 'gptAgents':
        result = gptAgentQuerier(prompt=query.prompt, contextDocs=documents, version=version)
    if engine == 'geminiQuerier':
        texts = []
        for docType, contents in documents.items():
            texts.append(f"{docType}=={'\n'.join(contents)}")
        result = geminiQuerier(prompt=query.prompt, texts=texts, version=version)
    
    return { 'message': "Success", 'data': result }


