FROM python:latest

MAINTAINER "dipankar das<dipankardas0115@gmail.com>"

WORKDIR /app

COPY requirements.txt .

RUN pip install -r requirements.txt

COPY . .

RUN mkdir models

ENTRYPOINT ["python"]
