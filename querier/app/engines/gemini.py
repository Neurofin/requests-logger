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
    vertexai.init(project="gemini-dnb", location="us-central1", credentials=credentials)
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
        print(err)
    return responseText

from vertexai.generative_models import Part
from utils.parsers import textjson_to_json

def geminiQuerier(texts: list[str], prompt: str, version: str):
    # documents = []
    # for text in texts:
    #     document = Part.from_data(
    #         mime_type="text/plain",
    #         data=text.encode('utf-8'))
    #     documents.append(document)

    response = generate(input_array=[*texts, prompt], version=version)
    print(f"Query processing for texts complete")
    try:
        json_format = textjson_to_json(response)
        result = json_format
    except Exception as err:
        print(
            "Error in deserialzing Gemini response. Cause: ",
            prompt,
            "Response: ",
            response,
            err,
        )
    return result
