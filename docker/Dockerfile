FROM python:3.9-alpine

EXPOSE 9118

COPY ./src /src
WORKDIR /src
RUN pip install -r requirements.txt

CMD ["python", "main.py"]