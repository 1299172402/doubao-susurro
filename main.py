"""主程序入口"""

import time
import pyperclip
from listener import poll_new_message

if __name__ == "__main__":
    while True:
        msg = poll_new_message()
        if msg:
            print(f"[已复制到剪贴板] {msg}")
            pyperclip.copy(msg)

        time.sleep(1)
