from flask import Flask, request, render_template
import grpc
import chat_pb2
import chat_pb2_grpc

app = Flask(__name__)

def get_grpc_response(message):
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = chat_pb2_grpc.ChatServiceStub(channel)
        response = stub.Chat(chat_pb2.ChatRequest(message=message))
        return response.reply

@app.route('/', methods=['GET', 'POST'])
def index():
    if request.method == 'POST':
        user_message = request.form['message']
        bot_reply = get_grpc_response(user_message)
        return render_template('index.html', user_message=user_message, bot_reply=bot_reply)
    return render_template('index.html')

if __name__ == '__main__':
    app.run(debug=True)