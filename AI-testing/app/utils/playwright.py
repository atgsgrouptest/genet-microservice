import psutil
from playwright.sync_api import sync_playwright

def find_main_browser_pid(user_data_dir):
    procs = []
    for p in psutil.process_iter(['pid', 'name', 'cmdline', 'ppid', 'create_time']):
        info = p.info
        name = info.get('name', '').lower()
        if name in ('chrome', 'chrome.exe', 'chromium', 'chromium-browser'):
            cmd = ' '.join(info.get('cmdline') or [])
            if user_data_dir in cmd:
                procs.append(info)
    # Filter out those whose parent is also in this list (likely renderer/helper)
    parent_pids = {p['ppid'] for p in procs}
    root_procs = [p for p in procs if p['pid'] not in parent_pids]
    if not root_procs:
        return None
    
    main = min(root_procs, key=lambda p: p['create_time'])
    return main['pid']

def run_playwright_and_get_pids(script_str, user_data_dir='/tmp/pw_session'):
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=False, args=[f'--user-data-dir={user_data_dir}'])
        page = browser.new_page()
        exec(script_str, {}, {'page': page})
        pid = find_main_browser_pid(user_data_dir)
        browser.close()
        return pid
