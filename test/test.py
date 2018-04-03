#!/usr/bin/env python
import time
import requests
import json

# NBCAPID_URL="https://www.ninechain.net/api/v2"
NBCAPID_URL="http://localhost:8080/api/v2"

def create():
    url = NBCAPID_URL
    params = {
        "channel":"skynovo",
        "records":[{
            "key":"00000000000000000000000000000001",
             "value":"1001"
         } ,{
            "key":"00000000000000000000000000000002",
            "value":"2001"
        },{
            "key":"00000000000000000000000000000001",
            "value":"1002"
        }
        ]
    }

    jsondata = {"params": params, "jsonrpc": "2.0", "id": 0, "method": "source-insert-batch"}

    # s = json.dumps(jsondata)
    # sign = sha256(s+nonce+secret_key+timestamp)
    # jsondata = {"sign":sign,"nonce":nonce, "timestamp":"xx", data=s}

    res = requests.post(url, json=jsondata, headers={'content-type': 'application/json', 'X-Api-Key':'1957394652937f2fa9b10c92cba464acd98c61f0249e82046d170c767c1161d6'}, verify=False)
    s = res.content.decode()
    print(json.dumps(jsondata))
    print(s)

def query_transactions():
    url = NBCAPID_URL
    params = {
        "channel":"skynovo",
        "key":"00000000000000000000000000000001",
    }

    jsondata = {"params": params, "jsonrpc": "2.0", "id": 0, "method": "source-transactions"}
    res = requests.post(url, json=jsondata, headers={'content-type': 'application/json', 'X-Api-Key':'1957394652937f2fa9b10c92cba464acd98c61f0249e82046d170c767c1161d6'}, verify=False)
    s = res.content.decode()
    print(json.dumps(jsondata))
    print(s)

def query_transaction():
    url = NBCAPID_URL
    params = {
        "channel":"skynovo",
        "tx_id":"4ed93d61d54b14cfd4e6cfe64d5277d099940d3c8033111e69ae36a73dddf235",
        # "tx_id":"8796b627a09a64ea01c0465b6d073e16f6d0efc0b059189246ac6c4db23f7732"
    }

    jsondata = {"params": params, "jsonrpc": "2.0", "id": 0, "method": "source-transaction"}
    res = requests.post(url, json=jsondata, headers={'content-type': 'application/json', 'X-Api-Key':'1957394652937f2fa9b10c92cba464acd98c61f0249e82046d170c767c1161d6'}, verify=False)
    s = res.content.decode()
    print(json.dumps(jsondata))
    print(s)

def query_state():
    url = NBCAPID_URL
    params = {
        "channel": "skynovo",
        "key": "00000000000000000000000000000001",
    }

    jsondata = {"params": params, "jsonrpc": "2.0", "id": 0, "method": "source-state"}
    res = requests.post(url, json=jsondata, headers={'content-type': 'application/json', 'X-Api-Key': '1957394652937f2fa9b10c92cba464acd98c61f0249e82046d170c767c1161d6'}, verify=False)
    print(url, json.dumps(jsondata))
    s = res.content.decode()
    print(json.dumps(jsondata))
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
        # create()
        # query_transactions()
        # query_state()
        query_transaction()

        # sdk_create()
        # sdk_transactions()
        # sdk_state()
    except Exception as e:
        print(e)

