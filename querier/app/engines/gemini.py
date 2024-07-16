import vertexai
from vertexai.generative_models import GenerativeModel
import vertexai.preview.generative_models as generative_models

from google.oauth2 import service_account

# Load credentials
credentials = service_account.Credentials.from_service_account_file('./gemini-dnb-24c09ae181db.json')

# Configuration settings
generation_config = {
    "max_output_tokens": 8192,
    "temperature": 0,
    "top_p": 0.95,
}

safety_settings = {
    generative_models.HarmCategory.HARM_CATEGORY_HATE_SPEECH: generative_models.HarmBlockThreshold.BLOCK_MEDIUM_AND_ABOVE,
    generative_models.HarmCategory.HARM_CATEGORY_DANGEROUS_CONTENT: generative_models.HarmBlockThreshold.BLOCK_MEDIUM_AND_ABOVE,
    generative_models.HarmCategory.HARM_CATEGORY_SEXUALLY_EXPLICIT: generative_models.HarmBlockThreshold.BLOCK_MEDIUM_AND_ABOVE,
    generative_models.HarmCategory.HARM_CATEGORY_HARASSMENT: generative_models.HarmBlockThreshold.BLOCK_MEDIUM_AND_ABOVE,
}

def generate(input_array: list[str], version: str):
    vertexai.init(project="gemini-dnb", location="asia-south1", credentials=credentials)
    model = GenerativeModel(
        f"gemini-{version}",
    )
    #1.5-flash-001
    responseText = ""
    try:
        response = model.generate_content(
            input_array,
            generation_config=generation_config,
            safety_settings=safety_settings,
            stream=False,
        )
        responseText = response.text
    except Exception as err:
        raise Exception(str(err))
    return responseText

from vertexai.generative_models import Part
from utils.parsers import textjson_to_json

def geminiQuerier(documents: list[Part], prompt: str, version: str):
    response = generate(input_array=[*documents, prompt], version=version)
    json_format = textjson_to_json(response)
    result = json_format
    return result
