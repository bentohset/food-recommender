# go-eats
A food recommender system

## Scope

Telegram bot - Get user preferences for food, recommends places
  - based on mealtime (breakfast, brunch, lunch or dinner)
  - don't-wants, cuisines the user would not want
  - budget
  - mood (comfort, energy, indulgent, healthy, adventurous(out-of-user prefs))
- get random restaurant

ChatGPT API - get dish recommendations

DB 
- store restaurants
  - budget range
  - mood
  - cuisine type
  - location

Recommendation Algorithm
- Data collection
  - Manual entry
  - Web scraper
  - collect user data as its entered for past restaurant preferences
- Data preprocessing
- Feature engineering
  - Popularity of restaurants
  - average ratings
- Model selection
- Model training
- Model evaluation
- Recommendation
  - incorporate model into recommendation system
- Improvements
  - periodically retrain model

Web form - add restaurant reviews
- Name of place (auto complete with google maps api) only in SG
- Chain?
- Budget
- Mood dropdown single select
- Cuisine dropdown multiple
- Mealtime dropdown multiple
- Rating stars

Web form - approval from admin
- pending table
  - able to edit, delete and approve
  - sort by column
- approved table
  - able to edit and delete
  - pagination
  - sort by column

## User stories

- As an indecisive user I want to get food recommendation based on my mood and budget
- As an adventurous person I want to get food outside my comfort zone
- As a meticulous person I want to see which dish is recommended so that I can maximise my budget

## Planning
Version 1.0 [DONE]
  - Telegram bot <-> Server
  - OpenAI API
  - Suggest food places through prompt engineering

Version 2.0
  - Telegram bot <-> Server <-> DB
  - Store restaurants
  - Suggest food based on DB
  - Allow users to add restaurants to a website (mobile based) and wait for pending approval
  - Admin accepts/rejects/edits restaurant requests to the DB

Version 3.0
  - Telegram bot <-> Server <-> DB
  - ML model Recommendation Algorithm

## Features
Commands
1. telegram input menu
  - mealtime (brunch, lunch, dinner, snack)
  - cuisine dont wants
  - budget ()
  - mood (comfort, energy, healthy, indulgent, exotic)

2. generative output
  - chatGPT enabled recommendation
  - generate based on machine learning model

3. user recommendations
  - submit personal recommendations through a website
  - pending approval from admin


mealtime -> budget -> mood (option menu)
dont wants (text input)

## Deployment
Convert telegram bot to webhooks from update polling
Deploy each service to K8s

.yaml files for both services in a `root/deployment` folder
- declare deployment manifests for kubernetes

Bot --- Recommender: gRPC
Form --- DB: REST
Recommender --- DB: REST

DB: Azure Database for PostgreSQL - Flexible

Host Kubernetes cluster on cloud free (1 node):
- Azure
