# genet-microservice-LLM_client
# 🧠 Genet Microservice - LLM-Client Service

This service acts as an Interface between external clients requests and the Adapter API, enabling interaction with models like `llama3:8b` via a simple HTTP endpoint.

---

## 🚀 Features

- Accepts structured model requests over HTTP
- Converts request into valid format and forwards it to the Adapter service API
- Returns the model's response or error in a structured format
- Built with **Go**, **Fiber**, and supports **CORS**


---

## 🧱 Project Structure
.
├── controllers/ # Fiber route handlers
├── factory/ # Model factory & adapter implementations
├── models/ # Request/Response/Error models
├── routes/ # Route registration
├── main.go # Entry point
├── go.mod
└── README.md # You're here


---

## 📦 Dependencies

- [Fiber v2](https://docs.gofiber.io/)
- [json-iterator/go](https://github.com/json-iterator/go)

---

## ⚙️ Setup & Run

1. **Install Go** (≥ 1.18): https://golang.org/dl/

2. **Install Ollama** and pull the model:

   ```bash
   ollama pull llama3:8b


3.  **Start Ollama Server**:

    ```bash
    ollama serve
    ```

4.  **Set Environment Variable for the App Port**:

    On Linux/macOS:

    ```bash
    export APP_PORT=<YOUR APP PORT>
    ```

    On Windows PowerShell:

    ```powershell
    $env:APP_PORT=<YOUR APP PORT>
    ```

5.  **Run the Service**:

    ```bash
    go run main.go
    ```

---

## 🔌 API Usage

### `POST /sendRequest`

This endpoint allows you to send prompts to an Adapter service.

**Request Body (JSON):**

```json
{
  "prompt": "Say any three words"
}
```
prompt: (String, required) The text prompt to send to the model.

Successful Response (JSON):
```json
  {
    "response": "\"Here are three words:\\n\\n```\\n{\\n\\\"words\\\": [\\\"Hello\\\", \\\"World\\\", \\\"Json\\\"]\\n}\\n```\""

 }
```
Error Response (JSON):
```json

{
  "service_name":"LLM_CLIENT",
  "error": "Invalid input",
  "description": "Prompt is required"
}

Important Notes
Ollama must be running on localhost:11434 for the service to communicate with it.
The model name in the request body must precisely match the name of the model as it appears when you run ollama list. For example, if you have llama3:8b installed, you must use llama3:8b in your request.