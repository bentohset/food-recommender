import pytest
from server.model import Model


@pytest.fixture(scope='module')
def grpc_add_to_server():
    from server.generated.route_pb2_grpc import add_MessageServiceServicer_to_server

    return add_MessageServiceServicer_to_server


@pytest.fixture(scope='module')
def grpc_servicer():
    from server.grpc import RecommendationServer
    model = Model()
    return RecommendationServer(model)


@pytest.fixture(scope='module')
def grpc_stub(grpc_channel):
    from server.generated.route_pb2_grpc import MessageServiceStub

    return MessageServiceStub(grpc_channel)