#!/usr/bin/env python
import time
import requests
import json

def test1():
    url = "http://localhost:8088/api/v1"
    params = {
        "channel":"my_channel",
        "records":[{
            "key":"k1",
            "value":"hash1"
        },{
            "key":"k2",
            "value":"hash2"
        },{
            "key":"k3",
            "value":"hash2"
        },{
            "key":"k4",
            "value":"hash2"
        }]
    }
    jsondata = {"params": params, "jsonrpc": "2.0", "id": 0, "method": "source-insert-batch"}
    res = requests.post(url, json=jsondata, headers={'content-type': 'application/json', 'X-Api-Key':'key111'})
    s = res.content.decode()
    print(s)
    pass


if __name__ == "__main__":
    try:
        test1()
    except Exception as e:
        print("error")
        print(e)

