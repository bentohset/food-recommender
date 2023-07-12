# food-recommender
A food recommender system

## Scope

Telegram bot - Get user preferences for food, recommends places
  - based on mealtime (breakfast, brunch, lunch or dinner)
  - don't-wants, cuisines the user would not want
  - hungry level (hungry, ok hungry, quick snack, just a drink)
  - budget
  - mood (comfort, energy, indulgent, healthy, adventurous(out-of-user prefs))
- get random restaurant

ChatGPT API - get dish recommendations

DB - store user preferences
  - dietary requirements
  - favourite cuisines
  - previous recommendations
- store restaurant reviews?TBA

Web form - add restaurant reviews

## User stories

- As an indecisive user I want to get food recommendation based on my mood and budget
- As an adventurous person I want to get food outside my comfort zone
- As a meticulous person I want to see which dish is recommended so that I can maximise my budget

## Planning
Version 1.0 
  - Telegram bot <-> Server
  - OpenAI API
  - Suggest food places through prompt engineering

Version 2.0
  - Telegram bot <-> Server <-> DB
  - Store user preferences
