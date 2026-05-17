"""豆包获取最新一次发送的消息"""

import requests
from curl_parser import get_config

_last_message_id = ""


def poll_new_message():
    """轮询一次，有新用户消息返回 tts_content，无则返回 None"""
    global _last_message_id

    url, params, headers, cookies, payload = get_config()

    response = requests.post(url, params=params, headers=headers, cookies=cookies, json=payload)

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


if __name__ == "__main__":
    msg = poll_new_message()
    print(msg)
