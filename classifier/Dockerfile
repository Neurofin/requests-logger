# start from the official python base image
FROM python:3.10

# set the current working dir to code
WORKDIR /code

# copy the file with requirements to code directory
COPY ./requirements.txt /code/requirements.txt

# install the package dependencies in the requirements file
RUN pip install --no-cache-dir --upgrade -r /code/requirements.txt

# copy the app inside the code directory
COPY ./app /code/app
COPY ./models /code/models

# 
CMD ["fastapi", "run", "app/main.py", "--port", "80"]