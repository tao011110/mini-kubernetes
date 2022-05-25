import json


# body:{'type_state1': 1/2}

def handler(env):
    var = env.body["type_state1"]
    ret = {'type': var, 'state': 3}
    return json.dumps(ret)
