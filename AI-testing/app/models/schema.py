from pydantic import BaseModel

class RequestFromGo(BaseModel):
    companyId: str
    requestId: str
    #automationcode
    #requestMaterial: str
    positiveCases: list[str]
    negativeCases: list[list[str]]
    
class PositiveResponseFromBrowserUse(BaseModel):
    requestId: str
    requestNo: int
    positiveCase: str
    positiveVideoUrl: str

class NegativeResponseFromBrowserUse(BaseModel):
    requestId: str
    requestNo: int
    negativeCase: str
    negativeVideoUrl: str