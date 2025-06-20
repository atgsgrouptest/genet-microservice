# genet-microservice
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

### `POST /sendRequestImages`

This endpoint allows you to send prompts to an Adapter service with images 

**Request Body (JSON):**

#### Remember request for this API endpoint should be in multipart/form-data
```Powershell
 curl.exe -X POST http://127.0.0.1:8002/sendRequestImages -F "prompt=Say any three words" -F "image=@C:\Users\lokes\OneDrive\Pictures\God of War Ragnarök\ScreenShot-2024-10-3_19-23-12.png"
```

prompt: (String, required) The text prompt to send to the model.

Successful Response (JSON):
``
{"response":"\"```json\\n{\\n  \\\"image_description\\\": {\\n    \\\"scene\\\": \\\"Dark, stone corridor within a crumbling, ancient structure.\\\",\\n    \\\"subject\\\": {\\n      \\\"type\\\": \\\"Character\\\",\\n      \\\"appearance\\\": \\\"A warrior or adventurer, wearing elaborate armor with intricate gold and dark designs.  They are holding a bow and arrow, looking ahead with a focused expression.\\\",\\n      \\\"pose\\\": \\\"Standing, facing forward, appearing to be contemplating something in the distance.\\\"\\n    },\\n    \\\"environment\\\": {\\n      \\\"walls\\\": \\\"Rough, uneven stone walls with a significant amount of yellow-green mineral deposits or corrosion. The walls are hewn with vertical, jagged edges.\\\",\\n      \\\"lighting\\\": \\\"Dim, likely from a single torch or light source casting strong shadows. The green tint suggests an otherworldly or subterranean environment.\\\",\\n      \\\"architecture\\\": \\\"Corridor-like, with a gradual narrowing of space.\\\",\\n      \\\"atmosphere\\\": \\\"Foreboding, mysterious, and potentially dangerous.\\\"\\n    },\\n    \\\"color_palette\\\": [\\n      \\\"Dark green\\\",\\n      \\\"Yellow-green\\\",\\n      \\\"Gold (metallic)\\\",\\n      \\\"Grey (stone)\\\"\\n    ],\\n    \\\"overall_impression\\\": \\\"The image evokes a sense of adventure, exploration, and the potential for danger within a forgotten, ancient place.\\\"\\n  }\\n}\\n```\""}
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