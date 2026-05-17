"""主程序入口"""

import pyperclip
from listener import listen

if __name__ == "__main__":
    for msg in listen():
        pyperclip.copy(msg)
        print(f"[已复制到剪贴板] {msg}")
