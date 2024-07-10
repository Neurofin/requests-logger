# start from the official python base image
FROM alpine:latest
# FROM winamd64/python:3.11-windowsservercore-ltsc2022
# FROM winamd64/python
# :3.11.9-windowsservercore-1809

# RUN apk --update --no-cache add python3~3.11.4 && ln -sf python3 /usr/bin/python
# RUN python3 -m ensurepip
# RUN pip3 install --no-cache --upgrade pip setuptools

# you can specify python version during image build
ARG PYTHON_VERSION=3.11.4

# install build dependencies and needed tools
RUN apk add \
    wget \
    gcc \
    make \
    zlib-dev \
    libffi-dev \
    openssl-dev \
    musl-dev

# download and extract python sources
RUN cd /opt \
    && wget https://www.python.org/ftp/python/${PYTHON_VERSION}/Python-${PYTHON_VERSION}.tgz \                                              
    && tar xzf Python-${PYTHON_VERSION}.tgz

# build python and remove left-over sources
RUN cd /opt/Python-${PYTHON_VERSION} \ 
    && ./configure --prefix=/usr --enable-optimizations --with-ensurepip=install \
    && make install \
    && rm /opt/Python-${PYTHON_VERSION}.tgz /opt/Python-${PYTHON_VERSION} -rf

RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools

RUN python3 --version

# set the current working dir to code
WORKDIR /code

# RUN python -m venv .venv

# ENV PATH="/.venv/bin:$PATH"

# copy the file with requirements to code directory
COPY ./requirements.txt /requirements.txt

# install the package dependencies in the requirements file
RUN pip install --no-cache-dir --upgrade -r /requirements.txt


# copy the app inside the code directory
COPY . .
# RUN .venv/Scripts/activate


EXPOSE 3008

# 
CMD ["fastapi", "run", "main.py", "--port", "3008"]