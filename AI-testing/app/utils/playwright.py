

def patch_js_script(original_js: str) -> str:
    # Replace chromium.launch with connectOverCDP
    patched_js = original_js.replace(
        "const browser = await chromium.launch({ headless: false });",
        """const wsUrl = process.env.CDP_WS_URL;
  if (!wsUrl) {
    console.error("CDP WebSocket URL not provided");
    process.exit(1);
  }
  const browser = await chromium.connectOverCDP(wsUrl);"""
    )
    patched_js = patched_js.replace("  const context = await browser.newContext();", "const context = browser.contexts()[0];")
    # Optional: Remove browser.close() to avoid closing shared instance
    patched_js = patched_js.replace("http://127.0.0.1:3000/index.html", " http://127.0.0.1:3000/AI-testing/index.html")
    patched_js = patched_js.replace(" const page = await context.newPage();", " const page = context.pages()[0];")
    return patched_js
