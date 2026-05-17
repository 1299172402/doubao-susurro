"""主程序入口"""

import time
import pyperclip
from datetime import datetime
from listener import poll_new_message

if __name__ == "__main__":
    while True:
        msg = poll_new_message()
        if msg:
            now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
            print(f"[{now}] {msg}")
            pyperclip.copy(msg)

        time.sleep(1)
