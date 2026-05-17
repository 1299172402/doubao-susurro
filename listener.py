"""豆包消息监听"""

import uuid
import time
import requests

from config import URL, PARAMS, HEADERS, COOKIES, PAYLOAD

_last_message_id = ""


def poll_new_message():
    """轮询一次，有新用户消息返回 tts_content，无则返回 None"""
    global _last_message_id

    payload = PAYLOAD.copy()
    payload["sequence_id"] = str(uuid.uuid4())

    response = requests.post(URL, params=PARAMS, headers=HEADERS, cookies=COOKIES, json=payload)

    if response.status_code != 200:
        return None

    data = response.json()
    chain = data["downlink_body"]["pull_singe_chain_downlink_body"]
    messages = chain["messages"]

    user_messages = [m for m in messages if m.get("user_type") == 1]

    if user_messages:
        latest = user_messages[0]
        if latest["message_id"] != _last_message_id:
            _last_message_id = latest["message_id"]
            return latest["tts_content"]

    return None


def listen():
    """持续监听，有新消息时 yield 返回"""
    while True:
        msg = poll_new_message()
        if msg:
            yield msg
        time.sleep(1)
