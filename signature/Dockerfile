# start from the official python base image
FROM python:3.11

# set the current working dir to code
WORKDIR /code

# copy the file with requirements to code directory
COPY ./requirements.txt /code/requirements.txt

# install the package dependencies in the requirements file
RUN pip install --no-cache-dir --upgrade -r /code/requirements.txt
RUN pip3 install opencv-python

# copy the app inside the code directory
COPY ./ /code

EXPOSE 3008

# 
CMD ["fastapi", "run", "main.py", "--port", "3008"]