from app.models.schema import RequestFromGo
from browser_use.llm import ChatOpenAI
from browser_use import ChatOllama, Agent, BrowserSession
import subprocess
import httpx
import asyncio
import os
from dotenv import load_dotenv

load_dotenv()

CHROME_PATH = r"/usr/bin/chromium"

def start_chrome_debugging(port=9222):
    chrome_command = [
        CHROME_PATH,
        f"--remote-debugging-port={port}",
        "--no-first-run",
        "--no-default-browser-check",
        "--disable-popup-blocking",
        "--disable-extensions",
        "--disable-gpu",
        "--no-sandbox",
        "--headless=true",
        "--user-data-dir=/tmp/chrome-profile",
    ]
    return subprocess.Popen(chrome_command)

async def get_cdp_websocket_url(host="localhost", port=9222):
    async with httpx.AsyncClient() as client:
        for _ in range(10):
            try:
                resp = await client.get(f"http://{host}:{port}/json/version")
                resp.raise_for_status()
                return resp.json()["webSocketDebuggerUrl"]
            except httpx.RequestError:
                await asyncio.sleep(0.5)
        raise RuntimeError("Chrome CDP not responding on port", port)

async def test(request_model: RequestFromGo, db):
    print("🚀 Launching Chrome...")
    chrome_proc = start_chrome_debugging()

    try:
        print("⏳ Fetching WebSocket Debugger URL...")
        cdp_ws_url = await get_cdp_websocket_url()

        llm = ChatOllama(model="gemma3:4b")

        browser_session = BrowserSession(cdp_url=cdp_ws_url)

        agent = Agent(
            task=(
                "get me to the http://127.0.0.1:3000/index.html, "
                "click on chatbot then conversation then 'new conversation' and type 'Hello'."
            ),
            llm=ChatOpenAI(
                model="gpt-4o-mini",
                api_key=os.getenv("OPENAI_API_KEY"),
            ),
            use_vision=True,
            browser_session=browser_session,
        )

        result = await agent.run()
        print("✅ Agent finished:")
        print(result)

        return True  # ✅ don't forget to return success

    finally:
        print("🧹 Closing Chrome...")
        chrome_proc.kill()
