# genet-microservice
# 🧠 Genet Microservice - Adapter Service

This service acts as an adapter to [Ollama](https://ollama.com) API, enabling interaction with models like `llama3:8b` via a simple HTTP endpoint.

---

## 🚀 Features

- Accepts structured model requests over HTTP
- Converts request into valid format and forwards it to the Ollama model API
- Returns the model's response or error in a structured format
- Built with **Go**, **Fiber**, and supports **CORS**
- Supports easy addition of new models via `ModelAdapter` interface

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
- [Ollama](https://ollama.com/)
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
    export APP_PORT=3000
    ```

    On Windows PowerShell:

    ```powershell
    $env:APP_PORT="3000"
    ```

5.  **Run the Service**:

    ```bash
    go run main.go
    ```

---

## 🔌 API Usage

### `POST /modelRequest`

This endpoint allows you to send prompts to an Ollama model.

**Request Body (JSON):**

```json
{
  "model": "llama3:8b",
  "prompt": "Say any three words"
}
model: (String, required) The name of the Ollama model to use. This must exactly match a model listed by ollama list (e.g., llama3:8b).
prompt: (String, required) The text prompt to send to the model.
Successful Response (JSON):
```

```json
{
  "message": "Here are three words:\n\n1. Sunshine\n2. Butterflies\n3. Laughter"
}
```

Error Response (JSON):

```json
{
  "service_name":"<NAME of SERVICE>",
  "error": "Invalid input",
  "description": "Prompt is required"
}
```
Adding a New Model
To extend the service with support for new models, follow these steps:

Implement the ModelAdapter interface. Create a new Go struct for your adapter (e.g., bertAdapter).
Extend the GetModelType() function in factory.go. Modify this function to return an instance of your new adapter when its corresponding model type is requested.
Important Notes
Ollama must be running on localhost:11434 for the service to communicate with it.
The model name in the request body must precisely match the name of the model as it appears when you run ollama list. For example, if you have llama3:8b installed, you must use llama3:8b in your request.