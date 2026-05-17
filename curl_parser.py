"""curl 解析工具：从 session.txt 读取 curl，解析出 URL/params/headers/cookies/payload"""

import re
import json
from urllib.parse import urlparse, parse_qs


def read_curl_file(path="session.txt"):
    """从文件读取 curl 命令"""
    with open(path, "r", encoding="utf-8") as f:
        return f.read().replace("\\\n", " ")  # 处理行尾反斜杠续行


def parse_curl(curl_str):
    """解析 curl 命令，返回 (url, params, headers, cookies, payload)"""

    # 提取 URL
    url_match = re.search(r"curl\s+'([^']+)'", curl_str)
    if not url_match:
        url_match = re.search(r'curl\s+"([^"]+)"', curl_str)
    if not url_match:
        raise ValueError("无法提取 URL")
    full_url = url_match.group(1)

    # 分离 base URL 和 query params
    parsed = urlparse(full_url)
    url = f"{parsed.scheme}://{parsed.netloc}{parsed.path}"
    params = {k: v[0] for k, v in parse_qs(parsed.query).items()}

    # 提取 headers（-H 'key: value'）
    headers = {}
    for match in re.finditer(r"-H\s+'([^']+)'", curl_str):
        line = match.group(1)
        if ":" in line:
            key, val = line.split(":", 1)
            headers[key.strip()] = val.strip()

    # 提取 cookies（-b 'key=val; key=val'）
    cookies = {}
    cookie_match = re.search(r"-b\s+'([^']+)'", curl_str)
    if cookie_match:
        for pair in cookie_match.group(1).split("; "):
            if "=" in pair:
                k, v = pair.split("=", 1)
                cookies[k] = v

    # 提取 payload（--data-raw '{...}'）
    payload = {}
    data_match = re.search(r"--data-raw\s+'(\{.*\})'", curl_str)
    if not data_match:
        data_match = re.search(r'--data-raw\s+"(\{.*\})"', curl_str)
    if data_match:
        try:
            payload = json.loads(data_match.group(1))
        except json.JSONDecodeError:
            pass

    return url, params, headers, cookies, payload


def get_config(file_path="session.txt"):
    """从 session.txt 解析配置"""
    curl_str = read_curl_file(file_path)
    url, params, headers, cookies, payload = parse_curl(curl_str)

    body = payload["uplink_body"]["pull_singe_chain_uplink_body"]
    body["direction"] = 0  # 向新消息方向
    body["anchor_index"] = 0  # 每次都从最新消息开始拉取

    return url, params, headers, cookies, payload


if __name__ == "__main__":
    url, params, headers, cookies, payload = get_config("session.txt")

    print(f"URL: {url}")
    print(f"PARAMS: {json.dumps(params, indent=2)}")
    print(f"HEADERS: {json.dumps(headers, indent=2)}")
    print(f"COOKIES: {json.dumps(cookies, indent=2)}")
    print(f"PAYLOAD: {json.dumps(payload, indent=2, ensure_ascii=False)}")
