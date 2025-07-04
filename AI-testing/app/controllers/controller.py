from browser_use import Agent, BrowserSession, BrowserProfile
from langchain_openai import ChatOpenAI
from langchain_ollama import ChatOllama
from dotenv import load_dotenv
from app.utils.playwright import run_playwright_and_get_pids
from pathlib import Path
import os

load_dotenv()


async def run_web_test(req):
    browser_session = BrowserSession(cdp_url="http://localhost:9222")
    if not req.prompt:
        raise ValueError("Task prompt is required")

    llm = ChatOllama(model="llama3.1:8b", num_ctx=18000)
    task = req.prompt

    browser_profile = BrowserProfile(
        record_video_dir="./recordings",  # Directory to save .webm video recordings
        record_video_size={"width": 1280, "height": 720},  # Optional: set video size
        headless=False,  # Set to False to see the browser while recording
    )

    browser_session = BrowserSession(browser_profile=browser_profile,cdp_url="http://localhost:9222")

    agent = Agent(
        task=task,
        llm=ChatOpenAI(
            model="gpt-4o-mini",
            api_key=os.getenv("OPEN_API_KEY"),  # Replace with your real key
        ),
        use_vision=True,
        browser_session=browser_session,
    )

    result = await agent.run(max_steps=5)

    print("🔗 Visited URLs:", result.urls())
    print("📄 Extracted Content:", result.extracted_content())
    print("❌ Errors:", result.errors())
    print("🧠 Model Actions:", result.model_actions())
    print("✅ Finished:", result.is_done())

    steps = list(result)
    final_step = steps[-1] if steps else None

    if final_step and getattr(final_step, "success", False):
        print("✅ Task completed successfully.")
    else:
        print("❌ Task completed without success.")

    return {
        "urls": result.urls(),
        "extracted_content": result.extracted_content(),
        "errors": result.errors(),
        "model_actions": result.model_actions(),
        "is_done": result.is_done(),
    }