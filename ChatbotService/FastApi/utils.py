from openai import OpenAI
from dotenv import load_dotenv
import openai
import os

load_dotenv()
client = OpenAI()
client.api_key = os.getenv("OPENAI_API_KEY")
client.base_url = "https://api.avalai.ir/v1"
print(client.api_key)
def generate_description(input):
    messages = [
        {
            "role" : "system",
            "content" : """As an assistant provide user 
            with information on the task needed """
        }

    ]
    messages.append({"role" : "user", "content" : f"{input}"})
    completion = client.chat.completions.create(model="gpt-3.5-turbo",
    messages=messages)
    reply = completion.choices[0].message.content
    return reply