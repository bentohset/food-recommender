## Getting data
Scrape from eatbook.sg/category/food-reviews
- to discover urls dynamically: web crawling/spidering with scrapy
```
scrapy runspider crawler.py
```
- explore urls and scrape data: web scraping with Selenium
```
python scraper.py
```

## Data processing



## Feature engineering


## Model selection

## Model training

## Model evaluation

## Recommendation system

## Notes
python -m grpc_tools.protoc -I ./server/protos --python_out=./server/generated --grpc_python_out=./server/generated ./server/protos/route.proto

ensure both server and client have exact same .proto file
- except for options

python protoc generates wrong imports, fix *_grpc.py with `from . ` import prefix

## Todo
- Add a new route for getting personalised recommendation or use the current route (might need to add a description)
- create model with content-based filtering, convert summary or review to text and find similarity with cosine similarity index
https://towardsdatascience.com/hands-on-content-based-recommender-system-using-python-1d643bf314e4
- implement in client