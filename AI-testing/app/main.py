from fastapi import FastAPI
from app.routes import router
from app.database import connect_to_mongo, close_mongo_connection


app=FastAPI()
app.include_router(router)

@app.on_event("startup")
async def startup_db_client():
    await connect_to_mongo()

@app.on_event("shutdown")
async def shutdown_db_client():
    await close_mongo_connection()
    

if __name__ == '__main__':
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=5000, reload=True)
