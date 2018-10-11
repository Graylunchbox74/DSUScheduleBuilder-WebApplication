FROM python:3.7-alpine

WORKDIR /usr/src/app

COPY ./requirements.txt .
COPY ./run.py           .
COPY ./server ./server
COPY ./config ./config

RUN pip install --no-cache -r ./requirements.txt

# This needs to match what's in config/default.py
EXPOSE 80

#Use development environment
ENV APP_CONFIG=/usr/src/app/config/development.py

#Use production environment
#ENV APP_CONFIG=/usr/src/app/config/production.py
#ENV IN_PRODUCTION=False

CMD [ "python", "run.py" ]