from pydantic import BaseModel

class RequestFromGo(BaseModel):
    companyId: str
    requestId: str
    automationCode: bytes
    #requestMaterial: bytes
    positiveCases: list[str]
    negativeCases: list[list[str]]
    
class PositiveResponseFromBrowserUse(BaseModel):
    requestId: str
    requestNo: int
    extractedcontent: str
    positiveCase: str
    positiveVideoUrl: str
    success : bool

class NegativeResponseFromBrowserUse(BaseModel):
    requestId: str
    requestNo: int
    extractedcontent: str
    negativeCase: str
    negativeVideoUrl: str
    success : bool