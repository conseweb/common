## coin
Hyperledger fabric chaincode account model digital currency. Maybe the coin will be called Lepuscoin.

### coin unit
Lepuscoin now has 4 units:
+ tinycoin: the smallest unit of Lepuscoin(tc)
+ minicoin: tinycoin * 1000(mc)
+ smallcoin: minicoin * 1000(sc)
+ coin: smallcoin * 1000(cc)

***Notice***chaincode request(invoke/query) args coin unit always be cc, such as "1.999", means 1.999cc

### developer test
```
CORE_CHAINCODE_ID_NAME=mycc CORE_PEER_ADDRESS=0.0.0.0:7051 ./lepuscoin
```
```
docker exec -it xxxx bash
```
### deploy request

REST
```
{
    "jsonrpc": "2.0",
    "method": "deploy",
    "params": {
        "type": 1,
        "chaincodeID": {
            "path": "https://github.com/mintzhao/common/assets/coin"
        },
        "ctorMsg": {
            "function": "deploy",
            "args": [
            ]
        },
        "secureContext": ""
    },
    "id": 1
}
```
CLI
```
peer chaincode deploy -n mycc -l golang -c '{"Function":"deploy","Args":[]}'
```
### invoke request
#### awardMiner

REST
```
{
    "jsonrpc": "2.0",
    "method": "invoke",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "CHAINCODE_HASH_HERE"
        },
        "ctorMsg": {
            "function": "awardMiner",
            "args": [
                "addr1",
                "100"
            ]
        },
        "secureContext": ""
    },
    "id": 1
}
```
CLI
```
peer chaincode invoke -n mycc -l golang -c '{"Function":"awardMiner", "Args":["addr1", "100"]}'
```

### query request
#### queryAccount

REST
```
{
    "jsonrpc": "2.0",
    "method": "query",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "CHAINCODE_HASH_HERE"
        },
        "ctorMsg": {
            "function": "queryAccount",
            "args": [
                "addr1"
            ]
        },
        "secureContext": ""
    },
    "id": 1
}
```

CLI
```
peer chaincode query -n mycc -l golang -c '{"Function":"queryAccount", "Args":["addr1"]}'
```
#### query_cb

REST
```
{
    "jsonrpc": "2.0",
    "method": "query",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "CHAINCODE_HASH_HERE"
        },
        "ctorMsg": {
            "function": "query_cb",
            "args": [
            ]
        },
        "secureContext": ""
    },
    "id": 1
}
```

CLI
```
peer chaincode query -n mycc -l golang -c '{"Function":"query_cb", "Args":[]}'
```