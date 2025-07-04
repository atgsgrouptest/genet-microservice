from pydantic import BaseModel

class WebTestRequest(BaseModel):
    prompt:str
    language:str 

class RequestFromGo(BaseModel):
    request_id: str
    positive_cases: list[str]
    negative_cases: list[list[str]]
    
class ResponseFromBrowserUser(BaseModel):
    request_id: str
    request_no: int
    positive_case: str
    negative_case: list[str]
    video_url: str