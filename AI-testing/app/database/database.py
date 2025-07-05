from motor.motor_asyncio import AsyncIOMotorClient

MONGO_URL = "mongodb://localhost:27017"
MONGO_DB_NAME = "your_database"

client: AsyncIOMotorClient = None
db = None

async def connect_to_mongo():
    global client, db
    client = AsyncIOMotorClient(MONGO_URL)
    db = client[MONGO_DB_NAME]

    await ensure_collections_exist(["requests", "responses"])

async def close_mongo_connection():
    client.close()

async def ensure_collections_exist(required_collections: list[str]):
    existing_collections = await db.list_collection_names()

    for name in required_collections:
        if name not in existing_collections:
            await db.create_collection(name)
            print(f"Created missing collection: {name}")
        else:
            print(f"Collection already exists: {name}")
