import json


# body:{'type': 1/2, 'state': 1/2}

def handler(env):
    var = env.body["state"]
    ret = {'state_through': var}
    return json.dumps(ret)
