import os
import requests


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
    payload = {'result': result, 'error': error}
    print(payload)
    r = requests.post('http://'+apiserverURL, data=payload)


if __name__ == '__main__':
    main()
