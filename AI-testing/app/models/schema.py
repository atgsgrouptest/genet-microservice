from pydantic import BaseModel

class RequestFromGo(BaseModel):
    company_id: str
    request_id: str
    positive_cases: list[str]
    negative_cases: list[list[str]]
    
class PositiveResponseFromBrowserUse(BaseModel):
    request_id: str
    request_no: int
    positive_case: str
    positive_video_url: str

class NegativeResponseFromBrowserUse(BaseModel):
    request_id: str
    request_no: int
    negative_case: str
    negative_video_url: str