from fastapi import APIRouter
from app.models import WebTestRequest
from app.controllers import run_web_test

router=APIRouter()

@router.post("/runWebTest")
async def handle_web_test_request(req : WebTestRequest):
    return await run_web_test(req)