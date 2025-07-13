# Multiple Services
To make a call for multiple services  
```bash
curl -X POST http://192.168.66.2:8000/ -H "Content-Type: application/json" -H "Authorization: Basic a2VhLWFwaTprZWE=" -d "{\"command\": \"config-get\", \"service\": [\"dhcp4\",\"dhcp6\", \"d2\"]}"
```
### Output
```bash
[
  {
    "arguments": {
      "Dhcp4": {
        ...
      }
    },
    "result": 0
  },
  {
    "arguments": {
      "Dhcp6": {
        ...
      }
    },
    "result": 0
  },
  {
    "arguments": {
      "DhcpDdns": {
        ...
      }
    },
    "result": 0
  }
]

```
## Imporant!
You can't query the control agent like this. It must always be `[]` empty, so it must always be queried by itself.  