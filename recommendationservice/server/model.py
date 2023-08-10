import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity

dfplace = pd.read_csv('./server/places-data.csv', encoding='utf-8')
dfreview = pd.read_csv('./server/reviews.csv', encoding='unicode_escape')

class Model():
    def __init__(self):
        self.data = dfplace
        self.data['mood'] = self.data['mood'].apply(lambda x: x.lower().split(', '))
        self.data['mealtime'] = self.data['mealtime'].apply(lambda x: x.lower().split(', '))
        self.data['cuisine'] = self.data['cuisine'].apply(lambda x: x.lower().split(', '))
        self.reviews = dfreview
        print("Model loaded and running")


    def generate_recommendations(self, mealtime, mood, cuisine_dont_wants, budget):
        # receive user input (to be received thru gRPC in prod)
        print("test")
        user_mealtime = mealtime
        user_mood = mood
        user_cuisine_dont_wants = cuisine_dont_wants
        user_budget = budget

        user_cuisine_dont_wants_lower = [x.lower() for x in user_cuisine_dont_wants]

        # initial filtering of data based on user input
        filtered_restaurants = self.data[
            (self.data['budget'] <= user_budget) &
            (~self.data['cuisine'].apply(lambda x: any(item in user_cuisine_dont_wants_lower for item in x))) &
            (self.data['mood'].apply(lambda arr: user_mood.lower() in arr)) &
            (self.data['mealtime'].apply(lambda arr: user_mealtime.lower() in arr))
        ]

        shuffled_restaurants = filtered_restaurants.sample(frac=1)

        # Output top 3 restaurant names as recommendations
        recommendations = shuffled_restaurants.head(5)
        

        return recommendations.to_dict(orient="records")
    
    def process_text(self,text):
        # replace multiple spaces with one
        text = str(text)
        text = ' '.join(text.split())
        # lowercase
        text = text.lower()

        return text

    def index_from_title(self,df,title):
        return df[df['name']==title].index.values[0]


    # function that returns the title of the movie from its index
    def title_from_index(self,df,index):
        return df[df.index==index].name.values[0]
    
   
    def recommendations(self, name, df,cosine_similarity_matrix,number_of_recommendations):
        index = self.index_from_title(df,name)
        similarity_scores = list(enumerate(cosine_similarity_matrix[index]))
        similarity_scores_sorted = sorted(similarity_scores, key=lambda x: x[1], reverse=True)
        recommendations_indices = [t[0] for t in similarity_scores_sorted[1:]]
        recommendations_indices = list(dict.fromkeys(recommendations_indices))[:number_of_recommendations]
        # recommendations_indices = [t[0] for t in similarity_scores_sorted[1:(number_of_recommendations+1)]]
        # return df['name'].iloc[recommendations_indices]

        rec_df = df.iloc[recommendations_indices][['name']].drop_duplicates(subset='name').head(5)
        rec = rec_df.to_dict(orient="records")

        for place in rec:
            entry = self.data[self.data['name'] == place['name']].iloc[0]
            place['budget'] = entry['budget']
            place['cuisine'] = entry['cuisine']
            place['rating'] = entry['rating']

        return rec

    def generate_personalized(self, user_choice):
        self.reviews['review'] = self.reviews.apply(lambda x: self.process_text(x.review),axis=1)
        tf_idf = TfidfVectorizer(stop_words='english')

        tf_idf_matrix = tf_idf.fit_transform(self.reviews['review'])

        cosine_similarity_matrix = cosine_similarity(tf_idf_matrix, tf_idf_matrix)

        # Output the top similar restaurants
        recommendations = self.recommendations(user_choice, self.reviews, cosine_similarity_matrix, 20)

        return recommendations

    

    