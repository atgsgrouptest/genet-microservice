package prompt

import ("github.com/atgsgrouptest/genet-microservice/Processor-service/models"
"fmt"

	jsoniter "github.com/json-iterator/go"
)

func GeneratePromptAPISequence(swaggerjsonobject string) (string, error) {
    prompt := "You are an expert API documentation processor. Your task is to analyze a given OpenAPI 3.0 JSON specification and convert it into a structured JSON array of API request details.\n\nYou must be extremely thorough and process **every single endpoint** defined in the input. Do not stop until all paths and their corresponding HTTP methods have been processed.\n\n**KEY INSTRUCTIONS:**\n\n1.  **Process ALL Endpoints:** Iterate through every path (e.g., `/pet`, `/user/login`) and every HTTP method within each path (`get`, `post`, `put`, `delete`). Do not miss any.\n\n2.  **Extract Key Details:** For each endpoint, you must identify:\n    * The API path.\n    * The HTTP method.\n    * A clear description from the `summary` or `description` field.\n    * Any parameters (path, query, header).\n    * The structure of the `requestBody`, if one exists.\n    * The primary success response code (e.g., `200`).\n\n3.  **Create Example Request Bodies:** The most important step is handling `requestBody`. When you see a `$ref` like `#/components/schemas/Pet`, you must find that schema in the `components` section of the input and use its `properties` and `example` values to create a realistic JSON example for the request body. If there is no `requestBody`, the value should be `null`.\n\n4.  **Logically Sequence:** You must order the final JSON objects in a logical sequence of use. Follow this order strictly:\n    * **1. User Management:** Start with creating users, logging in, getting user details, updating, and deleting.\n    * **2. Pet Management:** Then, list all pet-related operations (add pet, update, find, upload image, delete).\n    * **3. Store/Order Management:** Follow with store operations (place order, get inventory, get order, delete order).\n    * **4. Session Logout:** The very last step should be user logout.\n\n5.  **Strict Output Format:**\n    * The final output MUST be **ONLY a single, valid JSON array** and nothing else.\n    * Do not include any explanations, apologies, or introductory sentences like \"Here is the JSON output:\".\n    * Each object in the array must conform **exactly** to the structure provided below.\n    * For the `\"path\"` key, do not include the leading slash (e.g., use `\"user/login\"` instead of `\"/user/login\"`).\n\n**OUTPUT STRUCTURE FOR EACH API OBJECT:**\n\n```json\n{\n    \"sequenceNumber\": 0,\n    \"description\": \"\",\n    \"url\": \"\",\n    \"path\": \"\",\n    \"httpMethod\": \"\",\n    \"contentType\": \"\",\n    \"headers\": {},\n    \"requestBody\": {},\n    \"expectedResponseCode\": \"\"\n}\n```\n\n**EXAMPLE OF THE REQUIRED OUTPUT:**\n\nHere are two examples showing exactly how to format the first two objects in the sequence:\n\n```json\n[\n  {\n    \"sequenceNumber\": 1,\n    \"description\": \"Create a new user.\",\n    \"url\": \"/api/v3/user\",\n    \"path\": \"user\",\n    \"httpMethod\": \"POST\",\n    \"contentType\": \"application/json\",\n    \"headers\": {\n      \"Content-Type\": \"application/json\"\n    },\n    \"requestBody\": {\n      \"id\": 10,\n      \"username\": \"theUser\",\n      \"firstName\": \"John\",\n      \"lastName\": \"James\",\n      \"email\": \"john@email.com\",\n      \"password\": \"12345\",\n      \"phone\": \"123-456-7890\",\n      \"userStatus\": 1\n    },\n    \"expectedResponseCode\": \"200\"\n  },\n  {\n    \"sequenceNumber\": 2,\n    \"description\": \"Logs user into the system.\",\n    \"url\": \"/api/v3/user/login?username=theUser&password=12345\",\n    \"path\": \"user/login\",\n    \"httpMethod\": \"GET\",\n    \"contentType\": \"\",\n    \"headers\": {},\n    \"requestBody\": null,\n    \"expectedResponseCode\": \"200\"\n  }\n]\n```\n\nNow, process the following OpenAPI JSON and generate the complete JSON array output, following all instructions perfectly.\n\n"

    return prompt, nil
}


func PositiveCasePrompt(SequencedApi models.APIWrapper) (string, error) {
    json := jsoniter.ConfigCompatibleWithStandardLibrary
    sequenceJSON, err := json.MarshalIndent(SequencedApi, "", "  ")
    if err != nil {
        return "", fmt.Errorf("failed to marshal SequencedApi: %w", err)
    }

    prompt := 
        "I will provide a request specifications (JSON/YAML). Analyze it exhaustively to generate a complete and dependency-respecting sequence of API calls that covers all cases and HTTP methods .\n" +
        "- Include all required parameters, headers, and bodies as defined in the spec.\n" +
        "Give at least 5 of test cases that cover the positive reponses per api call\n\n" +
        "-Return only JSON and no other sentence using the structure below\n" + 
      "There should be no \"{\\}\" in the response " +



        "type APIRequest struct {\n" +
        "    SequenceNumber       int               `json:\"sequenceNumber\"`\n" +
        "    Description          string            `json:\"description\"`\n" +
        "    URL                  string            `json:\"url\"`\n" +
        "    Path                 string            `json:\"path\"`\n" +
        "    HTTPMethod           string            `json:\"httpMethod\"`\n" +
        "    ContentType          string            `json:\"contentType\"`\n" +
        "    Headers              map[string]string `json:\"headers\"`\n" +
        "    RequestBody          map[string]any    `json:\"requestBody\"`\n" +
        "    ExpectedResponseCode string            `json:\"expectedResponseCode\"`\n" +
        "}\n\n" +
        "no extra context or explanation.\n\n" +
        "Here is the request sequence:\n" +
        string(sequenceJSON)

    return prompt, nil
}
