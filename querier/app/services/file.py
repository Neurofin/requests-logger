import requests

def get_download_url(bucket: str, key: str):
    try:
        url = "http://localhost:3001" + "/presign?bucket=" + bucket + "&key=" + key
        response = requests.get(url)
        response.raise_for_status()  # Raises HTTPError for bad responses
        return response.json()
    except requests.exceptions.HTTPError as http_err:
        print(f"HTTP error occurred: {http_err}")  # Handle specific HTTP errors
    except Exception as err:
        print(f"Other error occurred: {err}")  # Handle other errors
    return None

