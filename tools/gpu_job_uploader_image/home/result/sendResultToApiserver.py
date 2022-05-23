import os
import requests
from requests.adapters import HTTPAdapter, DEFAULT_POOLBLOCK
from requests.packages.urllib3.poolmanager import PoolManager


class SourcePortAdapter(HTTPAdapter):
    """"Transport adapter" that allows us to set the source port."""

    def __init__(self, port, *args, **kwargs):
        self.poolmanager = None
        self._source_port = port
        super().__init__(*args, **kwargs)

    def init_poolmanager(self, connections, maxsize, block=DEFAULT_POOLBLOCK, **pool_kwargs):
        self.poolmanager = PoolManager(
            num_pools=connections, maxsize=maxsize,
            block=block, source_address=('', self._source_port))


def main():
    result = 'EMPTY'
    if os.path.exists('/home/result/result.out'):
        result = open('/home/result/result.out').read()
    error = 'EMPTY'
    if os.path.exists('/home/result/error.err'):
        error = open('/home/result/error.err', 'r').read()
    print(result)
    print(error)
    apiserverURL = open('/home/result/apiserver_ip_and_port', 'r').read()
    print(apiserverURL)
    jobName = open('/home/job_name', 'r').read()
    print(jobName)
    payload = {'jobName': jobName, 'result': result, 'error': error}
    print(payload)
    s = requests.Session()
    s.mount('http://', SourcePortAdapter(80))
    r = s.post('http://'+apiserverURL+'/gpu_job_result', data=payload)

if __name__ == '__main__':
    main()
