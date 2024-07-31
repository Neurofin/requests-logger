import fastapi
from pydantic import BaseModel

from dotenv import load_dotenv

load_dotenv()

from utils.extractText import getContextDocumentsMapping
from utils.extractText import getContextDocumentBytesMapping

from engines.gpt import gptQuerier
from engines.gpt import gptAgentQuerier
from engines.gemini import geminiQuerier

import base64
from vertexai.generative_models import Part


router = fastapi.APIRouter()

class ContextDocument(BaseModel):
    docType: str
    docPath: str

class Query(BaseModel):
    contextDocuments: list[ContextDocument]
    docFormat: str | None = None
    prompt: str
    #gptQuerier-4o
    #gptAgents-4o
    #geminiQuerier-1.5-flash-001
    engine: str =  "gptAgents-4o"


@router.post("/querier/resolve")
def resolve(query: Query):

    if not query.prompt.strip():
        raise fastapi.HTTPException(status_code=400, detail="Prompt should not be empty.")
        
    for doc in query.contextDocuments:
        if not doc.docPath.strip():
            raise fastapi.HTTPException(status_code=400, detail="docPath should not be empty.")

    try:
        engine = query.engine  
        [engine, version] = query.engine.split("-", 1)

        result = {}
        if query.docFormat == None or query.docFormat == "":    
            documents = getContextDocumentsMapping(query.contextDocuments)

            if engine == 'gptQuerier':
                texts = []
                for docType, contents in documents.items():
                    texts.append(f"{docType}=={'\n'.join(contents)}")
                result = gptQuerier(prompt=query.prompt, contextTexts=texts, version=version)
            if engine == 'gptAgents':
                result = gptAgentQuerier(prompt=query.prompt, contextDocs=documents, version=version)
            if engine == 'geminiQuerier':
                inputDocuments = []
                for docType, contents in documents.items():
                    text = f"{docType}=={'\n'.join(contents)}"
                    encoded_string = base64.b64encode(text.encode("utf-8"))
                    document = Part.from_data(
                            mime_type="text/plain",
                            data=base64.b64decode(encoded_string.decode()))
                    inputDocuments.append(document)
                result = geminiQuerier(prompt=query.prompt, documents=inputDocuments, version=version)
            
            return { 'message': "Success", 'data': result }
        else:

            documents = getContextDocumentBytesMapping(query.contextDocuments)

            if engine == 'geminiQuerier':
                inputDocuments = []
                for docType, contents in documents.items():
                    # text = f"{docType}=={'\n'.join(contents)}"
                    for content in contents:
                        encoded_string = base64.b64encode(content)
                        document = Part.from_data(
                                mime_type="application/pdf",
                                data=base64.b64decode(encoded_string.decode()))
                        inputDocuments.append(document)
                result = geminiQuerier(prompt=query.prompt, documents=inputDocuments, version=version)
            
            return { 'message': "Success", 'data': result }
    except Exception as e:
        raise fastapi.HTTPException(status_code=500, detail=str(e))
