# Telegram Bot Service
This folder represents the telegram bot microservice
The telegram bot acts as a gRPC clinet to the recommendation service to request for a recommendation. It queries the user for user preferences and acts as the entry point for all generative outputs.

## Table of Contents
- [Setup](#setup)
- [Bot Features](#bot-features)
    - [Query User Prefs](#query-user-preferences)
    - [gRPC Requests](#grpc-requests)
    - [OpenAI API](#openai-api-unused)
- [Deployment](#deployment)
- [Todo](#todo)

## Setup
**Setup go environment:**
```
go mod download
```

**Setup .env variables:**
```
TELEGRAM_APITOKEN=<bottoken from botfather>
OPENAI_APITOKEN=<optional>
```

**Generate protoc files:**
```
protoc *.proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --proto_path=.
```
>**Note**
>Ensure both server and client have exact same .proto file (except for options)

**Run the telegram bot:**
```
go build
./botservice
```
or 
```
go run .
```

## Bot Features
### Query User Preferences
Caches user variables of:
- userState: the stage at which the user is at currently (which question in the process)
- userChoices: the preferences the user has made stored as an array
- userFinalID: the final message of the recommendation for easy cleanup/deleting of messages

Queries the user preferences of:
1. mealtime - select menu
2. budget - select menu
3. mood - select menu (subjective)
4. cuisine dont-wants - forced reply by text

It can generate 2 types of recommendations:
1. Initial recommendation
    - Based on user preferences, filters the data accordingly
    - Generates 5 random restaurants
    - Select menu for personalised recommendation
2. Personalised recommendation
    - Based on initial recommendation, when a user selects a single restaurant
    - Generates 5 restaurants similar to the user selected one
    - Accuracy varied as it is dependent on reviews

### gRPC Requests
It sends 2 types of requests:
1. GetRecommendations
    - gets a list of restaurants based on user preferences
2. GetPersonalised
    - gets a list of restaurants based on previous chosen recommended

Note:
Requests are given 100 seconds of timeout
- Future iterations can lower this timeout period

### OpenAI API [Unused]
Given the user preferences, generate an output by creating a prompt.
GPT3.5 Turbo is used

Unused in favour of ML Model.
- TBC: use for generating summaries
- OpenAI API has abit of latency when generating output -> loading state?

## Deployment
```
kubectl apply -f bot.yaml
```
Deployed as a ClusterIP service on port 82

## Todo
- change update polling to websockets
- include summaries

