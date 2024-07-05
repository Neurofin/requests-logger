import fastapi
from pydantic import BaseModel

from dotenv import load_dotenv

load_dotenv()

from utils.extractText import getContextDocumentText
from utils.extractText import getContextDocumentBytes

from engines.gemini import geminiQuerier
from vertexai.generative_models import Part

import base64


router = fastapi.APIRouter()

class Body(BaseModel):
    docPath: str
    docFormat: str | None = None
    prompt: str
    #gptQuerier-4o
    #gptAgents-4o
    #geminiQuerier-1.5-flash-001
    engine: str =  "geminiQuerier-1.5-flash-001"

@router.post("/querier/classify")
def classsify(body: Body):

    documents = []
    if body.docFormat == None or body.docFormat == "":
        text = getContextDocumentText(s3Path=body.docPath)
        encoded_string = base64.b64encode(text.encode("utf-8"))
        document = Part.from_data(
                mime_type="text/plain",
                data=base64.b64decode(encoded_string.decode()))
        documents.append(document)
    else: 
        bytes = getContextDocumentBytes(s3Path=body.docPath)
        encoded_string = base64.b64encode(bytes)
        document = Part.from_data(
                mime_type="application/pdf",
                data=base64.b64decode(encoded_string.decode()))
        documents.append(document)

    engine = body.engine  
    [engine, version] = body.engine.split("-", 1)
    
    result = {}
    if engine == 'geminiQuerier':
        result = geminiQuerier(prompt=body.prompt, documents=documents, version=version)
    
    return { 'message': "Success", 'data': result }
