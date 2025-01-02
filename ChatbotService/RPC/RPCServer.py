
import sys
sys.path.append('../ProtoBuffers')
sys.path.append('../chat_pb2_grpc.py')
sys.path.append('../chat_pb2.py')
import grpc
from concurrent import futures
import chat_pb2
import chat_pb2_grpc
import openai
import os
from dotenv import load_dotenv

load_dotenv()

openai.api_key = os.getenv("OPENAI_API_KEY")
openai.api_base = "https://api.avalai.ir/v1"

class ChatService(chat_pb2_grpc.ChatServiceServicer):
    def __init__(self, api_key):
        self.api_key = api_key

    def Chat(self, request, context):
        messages = [
            {
                "role": "system",
                "content": "As an assistant, provide the user with information on the task needed."
            },
            {
                "role": "user",
                "content": request.message
            }
        ]
        try:
            # Generate response using OpenAI API
            completion = openai.ChatCompletion.create(
                model="gpt-3.5-turbo",
                messages=messages
            )
            reply = completion.choices[0].message.content
            return chat_pb2.ChatReply(reply=reply)
        except Exception as e:
            # Log and return error response
            print(f"An error occurred: {e}")
            context.set_details(f"An error occurred: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            return chat_pb2.ChatReply(reply="An error occurred while processing your request.")

def serve():
    # Start gRPC server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    chat_pb2_grpc.add_ChatServiceServicer_to_server(
        ChatService(api_key=os.getenv("OPENAI_API_KEY")),
        server
    )
    server.add_insecure_port('[::]:50051')
    print("Server started on port 50051")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()