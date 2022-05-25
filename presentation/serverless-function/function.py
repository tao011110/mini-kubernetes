import json
import numpy


# activerIP:port/function/function_test?test_param1=1&test_param2=2
# body:{'userType': 1}

def handler(env):
    var1 = env.uri.get("test_param1", 0)
    var2 = env.uri.get("test_param2", 0)
    var3 = int(numpy.abs(int(var1) - int(var2)))
    var4 = env.body["userType"]
    ret = {'var3': var3, 'var4': var4}
    return json.dumps(ret)
