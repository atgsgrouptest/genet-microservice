# run_web_test.py
import asyncio
import sys

if sys.platform.startswith("win"):
    asyncio.set_event_loop_policy(asyncio.WindowsSelectorEventLoopPolicy())

from browser_use import Agent, BrowserSession, BrowserProfile
from browser_use.llm import ChatOpenAI
from dotenv import load_dotenv
from app.models.schema import RequestFromGo, PositiveResponseFromBrowserUse, NegativeResponseFromBrowserUse
from pathlib import Path
import os
import subprocess
import httpx

load_dotenv()
CHROME_PATH = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

def start_chrome_debugging(port=9222):
    return subprocess.Popen([
        CHROME_PATH,
        f"--remote-debugging-port={port}",
        "--no-first-run",
        "--no-default-browser-check",
        "--disable-popup-blocking",
        "--disable-extensions",
        "--disable-gpu",
        "--user-data-dir=/tmp/chrome-profile",
    ])

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

async def run_web_test(req: RequestFromGo, db):
    positive_responses = []
    negative_responses = []

    # 👍 Positive Cases
    for idx, case in enumerate(req.positiveCases):
        print("🚀 Launching Chrome for positive case...")
        chrome_proc = start_chrome_debugging()
        try:
            task = f"You have been given a chatbot to test at http://127.0.0.1:3000/index.html . Test it with the flow: {case}"
            browser_session = BrowserSession(
                browser_profile=BrowserProfile(
                    record_video_dir="./recordings",
                    record_video_size={"width": 1280, "height": 720},
                    headless=False,
                ),
                cdp_url=await get_cdp_websocket_url(),
            )
            agent = Agent(
                task=task,
                llm=ChatOpenAI(model="gpt-4o-mini", api_key=os.getenv("OPEN_API_KEY")),
                use_vision=True,
                browser_session=browser_session,
            )
            result = await agent.run(max_steps=30)

            video_filename = f"{req.requestId}_pos_{idx}.webm"
            video_path = os.path.join("recordings", video_filename)
            positive_responses.append(PositiveResponseFromBrowserUse(
                requestId=req.requestId,
                requestNo=idx,
                positiveCase=case,
                positiveVideoUrl=video_path
            ))
        finally:
            chrome_proc.terminate()
            chrome_proc.wait()

    # 👎 Negative Cases
    for idx, case_group in enumerate(req.negativeCases):
        print("🚀 Launching Chrome for negative case group...")
        chrome_proc = start_chrome_debugging()
        try:
            browser_session = BrowserSession(
                browser_profile=BrowserProfile(
                    record_video_dir="./recordings",
                    record_video_size={"width": 1280, "height": 720},
                    headless=False,
                ),
                cdp_url=await get_cdp_websocket_url(),
            )

            for ind, indvCase in enumerate(case_group):
                
                task = f"You have been given a chatbot to test in the '{req.positiveCases[idx]}' at http://127.0.0.1:3000/index.html for the negative case: {indvCase}"
                agent = Agent(
                    task=task,
                    llm=ChatOpenAI(model="gpt-4o-mini", api_key=os.getenv("OPEN_API_KEY")),
                    use_vision=True,
                    browser_session=browser_session,
                )
                result = await agent.run(max_steps=5)

                video_filename = f"{req.requestId}_neg_{idx}_{ind}.webm"
                video_path = os.path.join("recordings", video_filename)
                negative_responses.append(NegativeResponseFromBrowserUse(
                    requestId=req.requestId,
                    requestNo=idx * 100 + ind,
                    negativeCase=indvCase,
                    negativeVideoUrl=video_path
                ))
        finally:
            chrome_proc.terminate()
            chrome_proc.wait()

    # 🧾 Insert into DB
    if positive_responses:
        await db["positive_responses"].insert_many([resp.dict() for resp in positive_responses])
    if negative_responses:
        await db["negative_responses"].insert_many([resp.dict() for resp in negative_responses])

    return result.is_done()
