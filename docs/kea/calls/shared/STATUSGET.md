# StatusGet
## Control Agent
```bash
curl -X POST http://192.168.66.2:8000/ -H "Content-Type: application/json" -H "Authorization: Basic a2VhLWFwaTprZWE=" -d "{\"command\": \"status-get\"}"
```
### Output
```bash
[ { "arguments": { "pid": 14013, "reload": 595868, "uptime": 595868 }, "result": 0 } ]
```
## DHCP4
```bash
curl -X POST http://192.168.66.2:8000/ -H "Content-Type: application/json" -H "Authorization: Basic a2VhLWFwaTprZWE=" -d "{\"command\": \"status-get\", \"service\": [\"dhcp4\"]}"
```
### Output
```bash
[ { "arguments": { "dhcp-state": { "disabled-by-db-connection": [  ], "disabled-by-local-command": [  ], "disabled-by-remote-command": [  ], "disabled-by-user": false, "globally-disabled": false }, "multi-threading-enabled": true, "packet-queue-size": 64, "packet-queue-statistics": [ 0.0, 0.0, 0.0 ], "pid": 9708, "reload": 1041346, "sockets": { "status": "ready" }, "thread-pool-size": 4, "uptime": 1041346 }, "result": 0 } ]
```
## DHCP6
```bash
curl -X POST http://192.168.66.2:8000/ -H "Content-Type: application/json" -H "Authorization: Basic a2VhLWFwaTprZWE=" -d "{\"command\": \"status-get\", \"service\": [\"dhcp6\"]}"
```
### Output
```bash
[ { "arguments": { "dhcp-state": { "disabled-by-db-connection": [  ], "disabled-by-local-command": [  ], "disabled-by-remote-command": [  ], "disabled-by-user": false, "globally-disabled": false }, "extended-info-tables": false, "multi-threading-enabled": true, "packet-queue-size": 64, "packet-queue-statistics": [ 0.0, 0.0, 0.0 ], "pid": 14093, "reload": 595775, "sockets": { "status": "ready" }, "thread-pool-size": 4, "uptime": 595775 }, "result": 0 } ]
```
## DDNS
```bash
curl -X POST http://192.168.66.2:8000/ -H "Content-Type: application/json" -H "Authorization: Basic a2VhLWFwaTprZWE=" -d "{\"command\": \"status-get\", \"service\": [\"d2\"]}"
```
### Output
```bash
[ { "arguments": { "pid": 20312, "reload": 2420, "uptime": 2420 }, "result": 0 } ]
```