from fastapi import FastAPI
from utils import generate_description

app = FastAPI()

@app.get("/")
def read_root():
    return {"message": "Welcome to FastAPI!"}

@app.get("/items/{item_id}")
def read_item(item_id: int, q: str = None):
    return {"item_id": item_id, "q": q}
@app.get("/users/{user_id}")
def get_user(user_id: int):
    return {"user_id": user_id}

@app.post("/users/{user_id}")
def create_user(user_id: int):
    return {"user_id": user_id}

@app.post("/chatbot")
async def chatbot(input: str):
    response = generate_description(input)
    return {"message": response}