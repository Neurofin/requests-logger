import fastapi
from pydantic import BaseModel

from dotenv import load_dotenv

load_dotenv()

from utils.extractText import downloadAndReturnTexts
from engines.gpt import gptQuerier

router = fastapi.APIRouter()

class ContextDocument(BaseModel):
    docType: str
    docPath: str

class Query(BaseModel):
    contextDocuments: list[ContextDocument]
    prompt: str
    engine: str =  "gptQuerier-4o"


@router.post("/resolve")
def resolve(query: Query):
    texts = downloadAndReturnTexts(query.contextDocuments)

    engine = query.engine  
    [engine, version] = query.engine.split("-")

    result = {}
    if engine == 'gptQuerier':
        result = gptQuerier(prompt=query.prompt, contextTexts=texts, version=version)
    
    return { 'message': "Success", 'data': result }


