import json


# body:{'type': 1/2}

def handler(env):
    var = env.body["type"]
    ret = {'type_state1': var}
    return json.dumps(ret)
