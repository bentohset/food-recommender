from .generated import route_pb2_grpc, route_pb2
from . import model

class RecommendationServer(route_pb2_grpc.MessageServiceServicer):
    def __init__(self, model):
        self.model = model

    def GetRecommendation(self, request, context):
        recommendations = self.model.generate_recommendations(
            mealtime=request.mealtime,
            mood=request.mood,
            cuisine_dont_wants=request.cuisine_dont_wants,
            budget=request.budget
        )
        
        restaurants = [route_pb2.Place(name=r["name"], budget=r["budget"], cuisine=r["cuisine"], rating=r["rating"]) for r in recommendations]
        response = route_pb2.RecommendationResponse(recommendations=restaurants)

        return response
    
    def GetPersonalised(self, request, context):
        name = request.name
        personalised = self.model.generate_personalized(name)

        restaurants = [route_pb2.Place(name=r["name"], budget=r["budget"], cuisine=r["cuisine"], rating=r["rating"]) for r in personalised]
        response = route_pb2.RecommendationResponse(recommendations=restaurants)

        return response
    
    def Echo(self, request, context):
        return route_pb2.EchoReply(message=f'you said: {request.message}')