from __future__ import print_function
import logging

import grpc

from .generated import route_pb2, route_pb2_grpc

def run():
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = route_pb2_grpc.MessageServiceStub(channel)

        request = route_pb2.RecommendationRequest(
            mealtime="lunch", 
            mood="casual", 
            cuisine_dont_wants=["Italian", "Chinese"], 
            budget=20
        )

        response = stub.GetRecommendation(request)
        for item in response.recommendations:
            print(item)

if __name__ == '__main__':
    logging.basicConfig()
    run()