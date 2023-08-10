from concurrent import futures
import grpc
import time

from .generated import route_pb2_grpc
from .grpc import RecommendationServer
from .model import Model

_ONE_DAY_IN_SECONDS = 60 * 60 * 24

class Server:
    @staticmethod
    def run():
        model = Model()
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        route_pb2_grpc.add_MessageServiceServicer_to_server(RecommendationServer(model), server)
        server.add_insecure_port('[::]:50051')
        server.add_insecure_port('[::]:83')
        server.start()
        print("Server started on port 50051")
        try:    
            while True:
                time.sleep(_ONE_DAY_IN_SECONDS)
        except KeyboardInterrupt:
            server.stop(0)