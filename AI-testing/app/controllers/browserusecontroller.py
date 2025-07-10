# run_web_test.py
import asyncio
from browser_use import Agent, BrowserSession, BrowserProfile
from browser_use.llm import ChatOpenAI
from dotenv import load_dotenv
from app.models.schema import RequestFromGo, PositiveResponseFromBrowserUse, NegativeResponseFromBrowserUse
from pathlib import Path
import os
import subprocess
import httpx
from app.utils.playwright import patch_js_script
from bson import ObjectId
import time
import shutil

load_dotenv()

# Use Playwright's Chromium instead of system Chromium
PLAYWRIGHT_CHROMIUM_PATH = os.path.expanduser("~/.cache/ms-playwright/chromium-1179/chrome-linux/chrome")

import signal
user_data_dir = f"/tmp/chrome-profile-{int(time.time())}-{os.getpid()}"
def start_chrome_debugging(port=9222):
    # Create unique user data directory for each instance
  
    
    chrome_command = [
        PLAYWRIGHT_CHROMIUM_PATH,  # Use Playwright's Chromium
        f"--remote-debugging-port={port}",
        "--no-first-run",
        "--no-default-browser-check",
        "--disable-popup-blocking",
        "--disable-extensions",
        "--disable-gpu",
        "--no-sandbox",
        f"--user-data-dir={user_data_dir}",
        "--disable-dev-shm-usage",  # Helps with permission issues
        "--disable-background-timer-throttling",
        "--disable-backgrounding-occluded-windows",
        "--disable-renderer-backgrounding",
    ]

    env = os.environ.copy()
    if "DISPLAY" not in env:
        env["DISPLAY"] = ":0"

    return subprocess.Popen(
        chrome_command,
        env=env,
        preexec_fn=os.setsid  # 🔑 Start new process group
    )


async def run_node_playwright_script(script_path: str, cdp_url: str):
    env = os.environ.copy()
    env["CDP_WS_URL"] = cdp_url

    process = subprocess.Popen(["node", script_path], env=env)
    process.wait()
    return process.returncode

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

def cleanup_chrome_process(chrome_proc):
    try:
        if chrome_proc and chrome_proc.poll() is None:
            try:
                os.killpg(os.getpgid(chrome_proc.pid), signal.SIGTERM)
                chrome_proc.wait(timeout=5)
            except (subprocess.TimeoutExpired, PermissionError):
                print("⚠️ Timeout or permission error. Trying SIGKILL...")
                try:
                    os.killpg(os.getpgid(chrome_proc.pid), signal.SIGKILL)
                    chrome_proc.wait()
                except PermissionError as e:
                    print(f"⚠️ Final cleanup failed: {e}")
    except Exception as e:
        print(f"Warning: General error cleaning up Chrome process: {e}")


async def run_web_test(req: RequestFromGo, db):
    js_code = req.automationCode.decode("utf-8")
    patched_code = patch_js_script(js_code)
    
    script_filename = f"playscript_{req.requestId}.js"
    script_path = os.path.join("playwright_scripts", script_filename)
    os.makedirs("playwright_scripts", exist_ok=True)

    with open(script_path, "w", encoding="utf-8") as f:
        f.write(patched_code)
    
    positive_responses = []
    negative_responses = []

    # 👍 Positive Cases
    print(f"Processing {len(req.positiveCases)} positive cases...")
    for idx, case in enumerate(req.positiveCases):
        print(f"🚀 Launching Chrome for positive case {idx + 1}/{len(req.positiveCases)}...")
        chrome_proc = None
        try:
            chrome_proc = start_chrome_debugging()
            await asyncio.sleep(2)  # Give Chrome time to start
            
            cdp_url = await get_cdp_websocket_url()
            await run_node_playwright_script(script_path, cdp_url)
            
            task = f"You have been given a chatbot to test . You can click on bubble or type in chat.After each action wait 3 seconds . Test it with the flow: {case}"
            
            browser_session = BrowserSession(
                browser_profile=BrowserProfile(
                    record_video_dir="./recordings",
                    record_video_size={"width": 1280, "height": 720},
                    headless=False,
                ),
                cdp_url=cdp_url,
            )
            
            agent = Agent(
                task=task,
                llm=ChatOpenAI(model="gpt-4.1-nano", api_key=os.Getenv("OPEN_API_KEY")),
                use_vision=True,
                browser_session=browser_session,
            )
            
            result = await agent.run(max_steps=5)
            
            video_filename = f"{req.requestId}_pos_{idx}.webm"
            video_path = os.path.join("recordings", video_filename)
            
            # Fix extractedcontent to be a string
            extracted_content = result.extracted_content()
            if isinstance(extracted_content, list):
                extracted_content = " ".join(str(item) for item in extracted_content)
            elif extracted_content is None:
                extracted_content = ""
            
            positive_responses.append(PositiveResponseFromBrowserUse(
                requestId=req.requestId,
                requestNo=idx,
                positiveCase=case,
                positiveVideoUrl=video_path,
                extractedcontent=str(extracted_content),
                success=result.is_done()
            ))
            
            print(f"✅ Positive case {idx + 1} completed successfully")
            
        except Exception as e:
            print(f"❌ Error in positive case {idx + 1}: {e}")
            break
        finally:
            cleanup_chrome_process(chrome_proc)
            break
            await asyncio.sleep(1)  # Brief pause between cases

    print(f"✅ Completed all positive cases. Results: {len(positive_responses)}")
    
    # 👎 Negative Cases
    print(f"Processing {len(req.negativeCases)} negative case groups...")
    for idx, case_group in enumerate(req.negativeCases):
        print(f"🚀 Launching Chrome for negative case group {idx + 1}/{len(req.negativeCases)}...")
        chrome_proc = None
        try:
            chrome_proc = start_chrome_debugging()
            await asyncio.sleep(2)  # Give Chrome time to start
            
            cdp_url = await get_cdp_websocket_url()
            await run_node_playwright_script(script_path, cdp_url)
            
            browser_session = BrowserSession(
                browser_profile=BrowserProfile(
                    record_video_dir="./recordings",
                    record_video_size={"width": 1280, "height": 720},
                    headless=False,
                ),
                cdp_url=cdp_url,
            )
            
            for ind, indvCase in enumerate(case_group):
                try:
                    # Use idx to get the corresponding positive case, but handle index bounds
                    positive_case_ref = req.positiveCases[idx] if idx < len(req.positiveCases) else req.positiveCases[0]
                    
                    task = f"You have been given a chatbot to test .After each action wait 3 seconds.You can click on bubble or type in chat for the negative case: {indvCase}"
                    
                    agent = Agent(
                        task=task,
                        llm=ChatOpenAI(model="gpt-4.1-nano", api_key=os.Getenv("OPEN_API_KEY")),
                        use_vision=True,
                        browser_session=browser_session,
                    )
                    
                    result = await agent.run(max_steps=5)
                    
                    video_filename = f"{req.requestId}_neg_{idx}_{ind}.webm"
                    video_path = os.path.join("recordings", video_filename)
                    
                    # Fix extractedcontent to be a string
                    extracted_content = result.extracted_content()
                    if isinstance(extracted_content, list):
                        extracted_content = " ".join(str(item) for item in extracted_content)
                    elif extracted_content is None:
                        extracted_content = ""
                    
                    negative_responses.append(NegativeResponseFromBrowserUse(
                        requestId=req.requestId,
                        requestNo=idx * 100 + ind,
                        negativeCase=indvCase,
                        negativeVideoUrl=video_path,
                        extractedcontent=str(extracted_content),
                        success=result.is_done()
                    ))
                    
                    print(f"✅ Negative case {idx + 1}.{ind + 1} completed successfully")
                    break
                except Exception as e:
                    print(f"❌ Error in negative case {idx + 1}.{ind + 1}: {e}")
                    break
                
        except Exception as e:
            print(f"❌ Error in negative case group {idx + 1}: {e}")
            break
        finally:
            cleanup_chrome_process(chrome_proc)
            break
            await asyncio.sleep(1)  # Brief pause between case groups

    print(f"✅ Completed all negative cases. Results: {len(negative_responses)}")
    
    # Clean up the script file
    try:
        os.remove(script_path)
    except OSError as e:
        print(f"Warning: Could not remove script file: {e}")
    
    print(f"Final responses - Positive: {len(positive_responses)}, Negative: {len(negative_responses)}")
    
    # 🧾 Insert into DB
    try:
        if positive_responses:
            await db["positive_responses"].insert_many([resp.dict() for resp in positive_responses])
            print(f"✅ Inserted {len(positive_responses)} positive responses into DB")
        if negative_responses:
            await db["negative_responses"].insert_many([resp.dict() for resp in negative_responses])
            print(f"✅ Inserted {len(negative_responses)} negative responses into DB")
            
        await db["completed_requests"].insert_one({"requestId": ObjectId(req.requestId)})
        print(f"✅ Marked request {req.requestId} as completed")
        
    except Exception as e:
        print(f"❌ Error inserting into database: {e}")
        raise
    shutil.rmtree(user_data_dir, ignore_errors=True)

    return True
