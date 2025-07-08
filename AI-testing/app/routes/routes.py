from fastapi import APIRouter, Request
from app.controllers.datacontroller import test
from app.models.schema import RequestFromGo

router = APIRouter()

@router.post("/test")
async def test_handler(request_model: RequestFromGo, request: Request):
    db = request.app.state.db  # ✅ access db via request
    return await test(request_model, db)
