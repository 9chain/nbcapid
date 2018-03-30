#!/usr/bin/env python
import time
import requests
import json

def create():
    url = "http://localhost:8088/api/v1"
    params = {
        "channel":"my_channel",
        "records":[{
            "key":"d961824324fb98f7b1b2b8e7278d960f1",
            "value":"d961824324fb98f7b1b2b8e7278d960f3712df7d8411fd36527a29fb5e64914b"
        },{
            "key":"d961824324fb98f7b1b2b8e7278d960f2",
            "value":"d961824324fb98f7b1b2b8e7278d960f3712df7d8411fd36527a29fb5e64914b"
        },{
            "key":"d961824324fb98f7b1b2b8e7278d960f3",
            "value":"d961824324fb98f7b1b2b8e7278d960f3712df7d8411fd36527a29fb5e64914b"
        },{
            "key":"d961824324fb98f7b1b2b8e7278d960f4",
            "value":"d961824324fb98f7b1b2b8e7278d960f3712df7d8411fd36527a29fb5e64914b"
        }]
    }

    jsondata = {"params": params, "jsonrpc": "2.0", "id": 0, "method": "source-insert-batch"}
    res = requests.post(url, json=jsondata, headers={'content-type': 'application/json', 'X-Api-Key':'key111'})
    s = res.content.decode()
    print(s)

def query_transactions():
    url = "http://localhost:8088/api/v1"
    params = {
        "channel":"my_channel",
        "key":"d961824324fb98f7b1b2b8e7278d960f3",
    }

    jsondata = {"params": params, "jsonrpc": "2.0", "id": 0, "method": "source-transactions"}
    res = requests.post(url, json=jsondata, headers={'content-type': 'application/json', 'X-Api-Key':'key111'})
    s = res.content.decode()
    print(s)

def query_state():
    url = "http://localhost:8088/api/v1"
    params = {
        "channel": "my_channel",
        "key": "d961824324fb98f7b1b2b8e7278d960f3",
    }

    jsondata = {"params": params, "jsonrpc": "2.0", "id": 0, "method": "source-state"}
    res = requests.post(url, json=jsondata, headers={'content-type': 'application/json', 'X-Api-Key': 'key111'})
    s = res.content.decode()
    print(s)

def sdk_create():
    url = "http://192.168.100.125:8082/v1/test"
    params = {
        "rid": "tx_786631_1",
        "channel": "businesschannel",
        "records": [
            {
                "key": "my_channel_d961824324fb98f7b1b2b8e7278d960f1",
                "value": "d961824324fb98f7b1b2b8e7278d960f3712df7d8411fd36527a29fb5e64914b"
            },
            {
                "key": "my_channel_d961824324fb98f7b1b2b8e7278d960f2",
                "value": "d961824324fb98f7b1b2b8e7278d960f3712df7d8411fd36527a29fb5e64914b"
            },
            {
                "key": "my_channel_d961824324fb98f7b1b2b8e7278d960f3",
                "value": "d961824324fb98f7b1b2b8e7278d960f3712df7d8411fd36527a29fb5e64914b"
            }
        ]
    }

    jsondata = {"params": params, "jsonrpc": "2.0", "id": 0, "method": "create"}
    res = requests.post(url, json=jsondata, headers={'content-type': 'application/json', 'X-Api-Key': '1234567890'})
    s = res.content.decode()
    print(s)

def sdk_transactions():
    url = "http://192.168.100.125:8082/v1/test"
    params = {
        "channel": "businesschannel",
        "key": "my_channel_d961824324fb98f7b1b2b8e7278d960f3",
    }


    jsondata = {"params": params, "jsonrpc": "2.0", "id": 0, "method": "transactions"}
    print(json.dumps(jsondata))
    res = requests.post(url, json=jsondata, headers={'content-type': 'application/json', 'X-Api-Key': '1234567890'})
    s = res.content.decode()
    print(s)


def sdk_state():
    url = "http://192.168.100.125:8082/v1/test"
    params = {
        "channel": "businesschannel",
        "key": "my_channel_d961824324fb98f7b1b2b8e7278d960f3",
    }

    jsondata = {"params": params, "jsonrpc": "2.0", "id": 0, "method": "state"}
    res = requests.post(url, json=jsondata, headers={'content-type': 'application/json', 'X-Api-Key': '1234567890'})
    s = res.content.decode()
    print(s)

# my_channel_d961824324fb98f7b1b2b8e7278d960f3
if __name__ == "__main__":
    try:
        create()
        query_transactions()
        query_state()

        # sdk_create()
        # sdk_transactions()
        # sdk_state()
    except Exception as e:
        print(e)

