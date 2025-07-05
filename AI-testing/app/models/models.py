from pydantic import BaseModel

class WebTestRequest(BaseModel):
    prompt:str
    language:str 