# Recommendation Service
This folder represents the recommendation microservice.
The machine learning model is packaged as a gRPC server which receives requests from the Telegram Bot as a gRPC client

## Table of Contents
- [Setup](#setup)
- [ML](#machine-learning)
    - [Getting data](#getting-data)
    - [Model selection](#model-selection)
    - [Recommendation system](#recommendation-system)
- [gRPC Server](#grpc-server)
- [Deployment](#deployment)
- [Todo](#todo)

## Setup
**Setup python virtual env**
```
pip install virtualenv
python -m venv venv
venv/scripts/activate
pip install -r requirements.txt
```

**Generate protoc files:**
```
python -m grpc_tools.protoc -I ./server/protos --python_out=./server/generated --grpc_python_out=./server/generated ./server/protos/route.proto
```

>**Note:**
>python protoc generates wrong imports, fix *_grpc.py by adding a `from . ` import prefix
>ensure both server and client have exact same .proto file (except for options)


**Setup .env variables**
```
DB_HOST=<azure db uri>
DB_USERNAME=postgres
DB_PASSWORD=<password>
DB_NAME=eats
GOOGLEMAPSAPI_KEY=
```

**Run tests**
```
pytest
```

**Run the server**
```
python server.py
```

## Machine Learning
### Getting data
cd into /data-collection

Scraped from eatbook.sg/category/food-reviews
- to discover urls dynamically: web crawling/spidering with scrapy
- only explored top 250+ urls
- outputs into discovered_urls.json
```
scrapy runspider crawler.py
```
- explore urls and scrape data: web scraping with Selenium
- outputs into output.csv
```
python scraper.py
```

Getting reviews
- using Google Maps API and getting the top 5 relevant reviews
- removes all forms of emojis, special e and other unicode
- outputs into review.csv
- reviews are partitioned into different rows with the same name due to space constraints
```
python reviews.py
```

### Model selection
Content-based filtering for generating similar outputs based on previous inputs

TF-IDF(Term Frequency - Inverse Document Frequency) used in common NLP text processing algorithms.
How it works:
- measures how important a term is within the field relative to other fields
- words are vectorized into importance numbers by TfidfVectorizer from scikit-learn library
- cosine similarity is used to identify closest matches



### Recommendation system
Given user preferences of
- mealtime
- cuisine dont-wants
- budget
- mood

Filter the data according to the preferences, randomly shuffle and output the top 5 as choices.

The user can then select any of the 5 for a more personalized recommendation through the ML model. The model finds similarities between reviews and outputs 5 of the highest similarity

## gRPC Server
This microservices acts as a gRPC server to package the ML model.

routes.MessageService
| RPC Method        | Request Message       | Response Message       | Action                                                        |
|-------------------|-----------------------|------------------------|---------------------------------------------------------------|
| Echo              | EchoRequest           | EchoReply              | Echoes the message sent                                       |
| GetRecommendation | RecommendationRequest | RecommendationResponse | Filters the data according to user preferences and returns it |
| GetPersonalised   | PersonalRequest       | RecommendationResponse | Finds similar restaurants                                     |

messages:
| Message                | Fields                                                                      |
|------------------------|-----------------------------------------------------------------------------|
| EchoRequest/EchoReply  | message string                                                              |
| RecommendationRequest  | mealtime string, mood string, cuisine_dont_wants string array, budget int   |
| RecommendationResponse | array of objects{name string, budget int, cuisine string array, rating int} |
| PersonalRequest        | name string                                                                 |

## Deployment
```
kubectl apply -f recommendation.yaml
```
Deployed as a ClusterIP service as it is meant for internal use by telegram bot


## Todo
- Query reviews and place data from PostgreSQL
- Setup methods to update model periodically

