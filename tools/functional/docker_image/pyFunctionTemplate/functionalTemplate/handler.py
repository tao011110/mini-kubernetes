def handler(url, body):
    print("username in url is" + url.get("username", ""))
    print(body["username"])
    return '{"hjk": "hjk"}'
