FROM python:3.9-slim

WORKDIR /recommendation-service
COPY requirements.txt /recommendation-service/

RUN python -m pip install --upgrade pip
RUN pip install -r requirements.txt

COPY . /recommendation-service

EXPOSE 50051
EXPOSE 83


ENTRYPOINT ["python", "server.py"]