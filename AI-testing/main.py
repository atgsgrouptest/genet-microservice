# main.py

import sys
import asyncio
from fastapi import FastAPI
import uvicorn
from motor.motor_asyncio import AsyncIOMotorClient
from app.routes.routes import router
from app.controllers.queuecontroller import start_consumer
MONGO_URL = "mongodb://root:example@192.168.0.100:27017/?authSource=admin"
MONGO_DB_NAME = "genet"

client: AsyncIOMotorClient = None
db = None

app = FastAPI()
app.include_router(router)

async def connect_to_mongo():
    global client, db
    if not client:
        client = AsyncIOMotorClient(MONGO_URL)
        db = client[MONGO_DB_NAME]
        print("✅ Connected to MongoDB using Motor")

async def close_mongo_connection():
    global client
    if client:
        client.close()
        print("❎ MongoDB connection closed")

@app.on_event("startup")
async def startup_db_client():
    await connect_to_mongo()
    app.state.db = db
    asyncio.create_task(start_consumer(app.state.db))

@app.on_event("shutdown")
async def shutdown_db_client():
    await close_mongo_connection()

if __name__ == '__main__':
    uvicorn.run("main:app", host="0.0.0.0", port=5000, reload=True)
