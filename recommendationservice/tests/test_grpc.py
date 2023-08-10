import server

def test_echo(grpc_stub):
    value = 'test-data'
    request = server.generated.route_pb2.EchoRequest(message=value)
    response = grpc_stub.Echo(request)

    assert response.message == f'you said: {value}'

def test_recommendation(grpc_stub):
    request = server.generated.route_pb2.RecommendationRequest(
        mealtime="lunch", 
        mood="casual", 
        cuisine_dont_wants=["Italian", "Chinese"], 
        budget=20
    )
    response = grpc_stub.GetRecommendation(request)

    assert len(response.recommendations) == 5