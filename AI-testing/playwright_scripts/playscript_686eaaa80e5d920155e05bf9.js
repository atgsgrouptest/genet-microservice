const { chromium } = require('playwright'); // or 'firefox' or 'webkit'

(async () => {
  // Launch browser
  const wsUrl = process.env.CDP_WS_URL;
  if (!wsUrl) {
    console.error("CDP WebSocket URL not provided");
    process.exit(1);
  }
  const browser = await chromium.connectOverCDP(wsUrl); // set to true to run headless
const context = browser.contexts()[0];
  const page = context.pages()[0];

  // Go to your local page
  await page.goto(' http://127.0.0.1:3000/AI-testing/index.html');

  // Optional: wait for body to be visible
  await page.waitForSelector('body');
  
  await browser.close();
})();
