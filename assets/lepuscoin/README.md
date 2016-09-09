## Lepuscoin
Hyperledger fabric chaincode account model & UTXO model digital currency

### developer test

run fabric peer using docker

docker-commpose.yml
```
membersrvc:
  image: hyperledger/fabric-membersrvc
  ports:
    - "7054:7054"
  command: membersrvc
vp0:
  image: hyperledger/fabric-peer
  ports:
    - "7050:7050"
    - "7051:7051"
    - "7053:7053"
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
  environment:
    - CORE_PEER_ADDRESSAUTODETECT=true
#    - CORE_VM_ENDPOINT=unix:///var/run/docker.sock
    - CORE_LOGGING_LEVEL=DEBUG
    - CORE_PEER_ID=vp0
    - CORE_PEER_PKI_ECA_PADDR=membersrvc:7054
    - CORE_PEER_PKI_TCA_PADDR=membersrvc:7054
    - CORE_PEER_PKI_TLSCA_PADDR=membersrvc:7054
    - CORE_SECURITY_ENABLED=false
    - CORE_SECURITY_ENROLLID=test_vp0
    - CORE_SECURITY_ENROLLSECRET=MwYpmSRjupbT
  links:
    - membersrvc
  command: sh -c "sleep 5; peer node start --peer-chaincodedev"
```
start docker
```
docker-compose up
```
build lepuscoin and run locally
```
CORE_CHAINCODE_ID_NAME=lepuscoin CORE_PEER_ADDRESS=0.0.0.0:7051 ./lepuscoin
```
```
docker ps -a
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
            "name": "lepuscoin"
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
peer chaincode deploy -n lepuscoin -l golang -c '{"Function":"deploy","Args":[]}'
```
### invoke request

#### invoke_coinbase
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
            "function": "invoke_coinbase",
            "args": [
            "CAEQiPrJvgUaBgj/////DyInCOgHEiIxNENXU0U4NjlpbkEzWGZ5Y29IOW1GaHl6Qmt3OVg2Yk40MAE="
            ]
        },
        "secureContext": ""
    },
    "id": 1
}
```
CLI
```
peer chaincode invoke -n lepuscoin -l golang -c '{"Function":"invoke_coinbase", "Args":["CAEQiPrJvgUaBgj/////DyInCOgHEiIxNENXU0U4NjlpbkEzWGZ5Y29IOW1GaHl6Qmt3OVg2Yk40MAE="]}'
```

#### invoke_transfer
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
            "function": "invoke_transfer",
            "args": [
            "CAEQiPrJvgUaQhJAY2M1ZjE5YTMwMTc4MjI1YTRjNjA2MTExZWQwMmYyMDBmN2IxODQ0ZWIwNzY2MDJhZDNkODUzYjIyOTRiYTFjOSIOCKAGEgkxMjM0NTY3ODkiJwjIARIiMTRDV1NFODY5aW5BM1hmeWNvSDltRmh5ekJrdzlYNmJONCoiMTRDV1NFODY5aW5BM1hmeWNvSDltRmh5ekJrdzlYNmJONA=="
            ]
        },
        "secureContext": ""
    },
    "id": 1
}
```
CLI
```
peer chaincode invoke -n lepuscoin -l golang -c '{"Function":"invoke_transfer", "Args":["CAEQiPrJvgUaQhJAY2M1ZjE5YTMwMTc4MjI1YTRjNjA2MTExZWQwMmYyMDBmN2IxODQ0ZWIwNzY2MDJhZDNkODUzYjIyOTRiYTFjOSIOCKAGEgkxMjM0NTY3ODkiJwjIARIiMTRDV1NFODY5aW5BM1hmeWNvSDltRmh5ekJrdzlYNmJONCoiMTRDV1NFODY5aW5BM1hmeWNvSDltRmh5ekJrdzlYNmJONA=="]}'
```
### query request
#### query_addr

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
            "function": "query_addr",
            "args": [
                "14CWSE869inA3XfycoH9mFhyzBkw9X6bN4"
            ]
        },
        "secureContext": ""
    },
    "id": 1
}
```

CLI
```
peer chaincode query -n lepuscoin -l golang -c '{"Function":"query_addr", "Args":["14CWSE869inA3XfycoH9mFhyzBkw9X6bN4"]}'
```
#### query_addrs
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
            "function": "query_addrs",
            "args": [
                "14CWSE869inA3XfycoH9mFhyzBkw9X6bN4",
                "123456789"
            ]
        },
        "secureContext": ""
    },
    "id": 1
}
```
CLI
```
peer chaincode query -n lepuscoin -l golang -c '{"Function":"query_addrs", "Args":["14CWSE869inA3XfycoH9mFhyzBkw9X6bN4","123456789"]}'
```
#### query_tx
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
            "function": "query_tx",
            "args": [
                "cc5f19a30178225a4c606111ed02f200f7b1844eb076602ad3d853b2294ba1c9",
            ]
        },
        "secureContext": ""
    },
    "id": 1
}
{
    "jsonrpc": "2.0",
    "method": "query",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "CHAINCODE_HASH_HERE"
        },
        "ctorMsg": {
            "function": "query_tx",
            "args": [
                "7ea4cca868a52c2dcd45a1ada8f62e586106067eb081409d977ced4b76c292a4",
            ]
        },
        "secureContext": ""
    },
    "id": 1
}
```
CLI
```
peer chaincode query -n lepuscoin -l golang -c '{"Function":"query_tx", "Args":["cc5f19a30178225a4c606111ed02f200f7b1844eb076602ad3d853b2294ba1c9"]}'
peer chaincode query -n lepuscoin -l golang -c '{"Function":"query_tx", "Args":["7ea4cca868a52c2dcd45a1ada8f62e586106067eb081409d977ced4b76c292a4"]}'
```
#### query_coin

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
            "function": "query_coin",
            "args": []
        },
        "secureContext": ""
    },
    "id": 1
}
```

CLI
```
peer chaincode query -n lepuscoin -l golang -c '{"Function":"query_coin", "Args":[]}'
```