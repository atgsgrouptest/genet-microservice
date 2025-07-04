package models

//This is the error reponse if prompt is not valid or Reponse has a problem

type Error struct{
	ServiceName string `json:"service_name"`
	Message string `json:"message"`
	Description string `json:"description"`
}


//An array of SwaggerObject that contains the host URL and the link to the Swagger JSON file
type Request struct{
    SwaggerRequest []SwaggerObject `json:"swaggerRequest"`
}

//An object that contains the host URL and the link to the Swagger JSON file
type SwaggerObject struct{
	HostURL string `json:"hostUrl"`
	SwaggerJSONLink string `json:"swaggerJsonLink"`
}

type OpenAPI struct {
    OpenAPI string `json:"openapi"`
    Info    struct {
        Title   string `json:"title"`
        Version string `json:"version"`
    } `json:"info"`
    Paths map[string]interface{} `json:"paths"`
}

type OpenAPIRequest struct {
    Request []HttpRequest `json:"request"`
}

type HttpRequest struct {
	URL	 string            `json:"url"`
	HTTPMethod string         `json:"httpMethod"`
	Description string     `json:"description"`
    Headers map[string]string `json:"headers"`
	Body map[string]interface{} `json:"body"`
}

type APIRequest struct {
    SequenceNumber       int               `json:"sequenceNumber"`
    Description          string            `json:"description"`
    URL                  string            `json:"url"`
    Path                 string            `json:"path"`
    HTTPMethod           string            `json:"httpMethod"`
    ContentType          string            `json:"contentType"`
    Headers              map[string]string `json:"headers"`
    RequestBody          map[string]any    `json:"requestBody"`
    ExpectedResponseCode string            `json:"expectedResponseCode"`
}

type APIWrapper struct {
    APIs []APIRequest `json:"apis"`
}

type LLMResponse struct {
	Response string `json:"response"`
}
