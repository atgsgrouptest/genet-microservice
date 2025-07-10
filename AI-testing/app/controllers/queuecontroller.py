import json
import os
import asyncio

import httpx
from aio_pika import connect_robust, IncomingMessage
from bson import ObjectId
from fastapi.responses import JSONResponse
from fastapi.encoders import jsonable_encoder
import base64

from app.models.schema import RequestFromGo
from app.controllers.browserusecontroller import run_web_test
RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
QUEUE_NAME = "request_ids"

async def handle_message(message: IncomingMessage,db):   
    from app.controllers.datacontroller import test     
    try:
        body = message.body.decode()
        data = json.loads(body)

        requestId = data.get("requestId")

        print(f"📨 Received message - requestID: {requestId} (type: {type(requestId)})")
        # Debug all requestIds from DB

        # Convert to ObjectId and query
        object_id = ObjectId(requestId)
        print("db id",db)
        doc = await db["requests"].find_one({"requestId": ObjectId(requestId)})
        #print(doc)
        if doc is None:
            print(f"❌ No document found for requestID: {requestId}, rejecting message.")
            await message.reject(requeue=False)
            return

        print(f"🔍 Found document for requestID: {requestId}")
        
        positiveCases = doc.get("positiveCases", [])
        negativeCases = doc.get("negativeCases", [[]])
        automationCode = doc.get("automationCode",bytearray)

        request_model = RequestFromGo(
            companyId=doc.get("companyId", ""),
            requestId=requestId,
            automationCode=automationCode,
            positiveCases=positiveCases,
            negativeCases=negativeCases
            )
        print(automationCode)
        
        done=await run_web_test(request_model,db)

        print("PositiveCases --->", request_model.positiveCases[0])
        print("NegativeCases --->", request_model.negativeCases[0][0],request_model.negativeCases[1][0])
        
        
        if not done:
            print("Test not done")
        # Success
        await message.ack()

    except Exception as e:
        print(f"❌ Error: {e}")
        await message.reject(requeue=True)
    
    print("All cases ran without issue")


async def start_consumer(db):
    RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@192.168.0.100:5672/")
    
    connection = await connect_robust(RABBITMQ_URL)
    channel = await connection.channel()
    queue = await channel.declare_queue(QUEUE_NAME, durable=True)

    print(f"🚀 Listening on RabbitMQ queue '{QUEUE_NAME}'...")
    await queue.consume(lambda msg: handle_message(msg, db), no_ack=False)