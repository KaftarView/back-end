import grpc
import sys
import os
sys.path.append('../ProtoBuffers')
import chat_pb2
import chat_pb2_grpc

def run():
    try:
        with grpc.insecure_channel('localhost:50051') as channel:
            stub = chat_pb2_grpc.ChatServiceStub(channel)
            
            response = stub.Chat(chat_pb2.ChatRequest(message="Hello, how can you assist me?"))
            print("ChatService replied: " + response.reply)
    except grpc.RpcError as e:
        print(f"RPC error: {e.code()} - {e.details()}")

if __name__ == '__main__':
    run()