import requests
import json
import time

url = ""

params = {}

headers = {}

cookies = {}

payload = {}

last_message_id = ""


def poll_new_message():
    """轮询一次，有新用户消息返回 tts_content，无则返回 None"""
    global last_message_id

    payload["sequence_id"] = str(__import__("uuid").uuid4())

    response = requests.post(url, params=params, headers=headers, cookies=cookies, json=payload)

    if response.status_code != 200:
        return None

    data = response.json()
    chain = data["downlink_body"]["pull_singe_chain_downlink_body"]
    messages = chain["messages"]

    user_messages = [m for m in messages if m.get("user_type") == 1]

    if user_messages:
        latest = user_messages[0]
        if latest["message_id"] != last_message_id:
            last_message_id = latest["message_id"]
            return latest["tts_content"]

    return None


def listen():
    """持续监听，有新消息时返回（生成器）"""
    while True:
        msg = poll_new_message()
        if msg:
            yield msg
        time.sleep(1)


if __name__ == "__main__":
    for new_msg in listen():
        print(f"[新消息] {new_msg}")
