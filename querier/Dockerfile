# start from the official python base image
FROM python:3.12

# set the current working dir to code
WORKDIR /code

# copy the file with requirements to code directory
COPY ./requirements.txt /code/requirements.txt

# install the package dependencies in the requirements file
RUN pip install --no-cache-dir --upgrade -r /code/requirements.txt
RUN pip install llama-index-readers-string-iterable
RUN pip install --upgrade google-cloud-aiplatform

# copy the app inside the code directory
COPY ./app /code/app
COPY ./gemini-dnb-24c09ae181db.json /code/gemini-dnb-24c09ae181db.json

EXPOSE 3006

# 
CMD ["fastapi", "run", "app/main.py", "--port", "3006"]