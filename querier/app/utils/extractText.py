import services.file
import requests

def get_text_from_url(url):
    """
    Fetches the content of a text file from the specified URL and assigns it to a variable.

    :param url: The URL of the text file.
    :return: The content of the text file as a string.
    """
    try:
        response = requests.get(url)
        response.raise_for_status()  # Raises HTTPError for bad responses
        return response.text
    except requests.exceptions.HTTPError as http_err:
        print(f"HTTP error occurred: {http_err}")  # Handle specific HTTP errors
    except Exception as err:
        print(f"Other error occurred: {err}")  # Handle other errors
    return None

def splitBucketAndKey(s3Path: str):
    path_parts=s3Path.replace("s3://","").split("/")
    bucket=path_parts.pop(0)
    key="/".join(path_parts)
    return bucket, key

def downloadAndReturnTexts(contextDocuments):
    texts = []
    for contextDocument in contextDocuments:
        path = contextDocument.docPath
        bucket, key = splitBucketAndKey(s3Path=path)
        response = services.file.get_download_url(bucket=bucket, key=key)
        if response == None:
            return None
        data = response["data"]
        url = data["URL"]
        text = get_text_from_url(url)
        texts.append(f"{contextDocument.docType}=={text}")
    return texts
