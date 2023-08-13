import requests
import os
import csv
from dotenv import load_dotenv
import pandas as pd
import re
import emoji

def get_top_reviews(api_key, restaurant_name):
    base_url = "https://maps.googleapis.com/maps/api/place/textsearch/json"
    
    params = {
        "query": restaurant_name,
        "key": api_key
    }
    
    response = requests.get(base_url, params=params)
    data = response.json()
    
    res = []

    if "results" in data:
        for result in data["results"]:
            place_id = result["place_id"]
            name = result["name"]
            # print(f"Restaurant: {name}")
            # print("Top 5 Reviews:")
            
            review_params = {
                "place_id": place_id,
                "key": api_key
            }
            
            review_url = "https://maps.googleapis.com/maps/api/place/details/json"
            review_response = requests.get(review_url, params=review_params)
            review_data = review_response.json()
            
            if "result" in review_data:
                reviews = review_data["result"].get("reviews", [])
                for i, review in enumerate(reviews[:5], start=1):
                    author_name = review.get("author_name", "N/A")
                    rating = review.get("rating", "N/A")
                    text = review.get("text", "N/A")
                    # print(f"Review {i}:")
                    # print(f"Author: {author_name}")
                    # print(f"Rating: {rating}")
                    # print(f"Text: {text}")
                    res.append(text)
            else:
                print("No reviews found for this restaurant.")
            
    return res

def convert_csv_array(filepath):
    data = pd.read_csv(filepath)
    names = data["name"].tolist()

    return names

def normalize_text(text):
    # Remove all \ followed by characters
    normalized_text = text.replace("é", "e")
    normalized_text = re.sub(r'\n', ' ', normalized_text)
    normalized_text = emoji.demojize(normalized_text)  # Remove emojis
    normalized_text = re.sub(r'\\[uU][0-9a-fA-F]{4}', '', normalized_text)
    normalized_text = re.sub(r'\\[^\s]+', '', normalized_text)
    return normalized_text


def output_file(data):
    field_names = ["name", "review"]
    with open('server/data-collection/reviews.csv', 'a', encoding='utf-8', newline='') as file:
        csv_writer = csv.DictWriter(file, fieldnames=field_names)
        csv_writer.writeheader()
        
        for item in data:
            name = item["name"]
            reviews = item["reviews"]
            for review in reviews:
                normal_text = normalize_text(review)
                csv_writer.writerow({"name": name, "review": normal_text})

def main():
    load_dotenv()

    api_key = os.getenv('GOOGLEMAPSAPI_KEY')
    names = convert_csv_array("server/places-data.csv")
    # restaurant_name = "Chir Chir"
    res = []

    for i, n in enumerate(names):
        print(str(i) + '/' + str(len(names)))
        rev_arr = get_top_reviews(api_key, n)
        res.append({
            "name":n,
            "reviews": rev_arr
        })

    output_file(res)



if __name__ == "__main__":
    main()

