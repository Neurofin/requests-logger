import fastapi
from pydantic import BaseModel

from dotenv import load_dotenv

load_dotenv()

from utils.extractText import getContextDocumentText

from engines.gpt import gptQuerier
from engines.gpt import gptAgentQuerier
from engines.gemini import geminiQuerier


router = fastapi.APIRouter()

class Body(BaseModel):
    docPath: str
    prompt: str
    #gptQuerier-4o
    #gptAgents-4o
    #geminiQuerier-1.5-flash-001
    engine: str =  "geminiQuerier-1.5-flash-001"

@router.post("/classify")
def resolve(body: Body):

    text = getContextDocumentText(s3Path=body.docPath)
    texts = [text]

    engine = body.engine  
    [engine, version] = body.engine.split("-", 1)
    
    result = {}
    if engine == 'geminiQuerier':
        result = geminiQuerier(prompt=body.prompt, texts=texts, version=version)
    
    return { 'message': "Success", 'data': result }
