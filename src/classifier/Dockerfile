FROM python:latest

MAINTAINER "dipankar das<dipankardas0115@gmail.com>"

WORKDIR /app

COPY requirements.txt .

RUN pip -r requirements.txt

COPY . .

CMD ["python", "main.py"]