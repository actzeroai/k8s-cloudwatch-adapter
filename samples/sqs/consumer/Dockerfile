FROM python:3.11-slim-bullseye

WORKDIR /usr/src/app

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

RUN pip install --upgrade pip
RUN pip install boto3~=1.26

COPY consumer.py .

ENTRYPOINT ["python", "consumer.py"]
