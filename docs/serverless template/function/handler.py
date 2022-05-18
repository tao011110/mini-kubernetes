def handler(env):
    print("username in url is" + env.uri.get("username", ""))
    print(env.body["username"])
    return '{"test": "test"}'
