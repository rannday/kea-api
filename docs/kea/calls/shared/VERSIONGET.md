# VersionGet

```bash
curl -X POST http://192.168.66.2:8000/ -H "Content-Type: application/json" -H "Authorization: Basic a2VhLWFwaTprZWE=" -d "{\"command\": \"version-get\"}"
```
### Output
```bash
[ { "arguments": { "extended": "2.6.3 (isc20250522135511 deb)\npremium: yes (isc20250522135511 deb)\nlinked with:\n- log4cplus 2.0.8\n- OpenSSL 3.0.16 11 Feb 2025" }, "result": 0, "text": "2.6.3" } ]
```
```bash
curl -X POST http://192.168.66.2:8000/ -H "Content-Type: application/json" -H "Authorization: Basic a2VhLWFwaTprZWE=" -d "{\"command\": \"version-get\", \"service\": [\"dhcp4\"]}"
```
### Output
```bash
[ { "arguments": { "extended": "2.6.3 (isc20250522135511 deb)\npremium: yes (isc20250522135511 deb)\nlinked with:\n- log4cplus 2.0.8\n- OpenSSL 3.0.16 11 Feb 2025\nbackends:\n- MySQL backend 22.2, library 3.3.14\n- PostgreSQL backend 22.2, library 150013\n- Memfile backend 3.0" }, "result": 0, "text": "2.6.3" } ]
```
```bash
curl -X POST http://192.168.66.2:8000/ -H "Content-Type: application/json" -H "Authorization: Basic a2VhLWFwaTprZWE=" -d "{\"command\": \"version-get\", \"service\": [\"dhcp6\"]}"
```
### Output
```bash
[ { "arguments": { "extended": "2.6.3 (isc20250522135511 deb)\npremium: yes (isc20250522135511 deb)\nlinked with:\n- log4cplus 2.0.8\n- OpenSSL 3.0.16 11 Feb 2025\nbackends:\n- MySQL backend 22.2, library 3.3.14\n- PostgreSQL backend 22.2, library 150013\n- Memfile backend 5.0" }, "result": 0, "text": "2.6.3" } ]
```
```bash
curl -X POST http://192.168.66.2:8000/ -H "Content-Type: application/json" -H "Authorization: Basic a2VhLWFwaTprZWE=" -d "{\"command\": \"version-get\", \"service\": [\"d2\"]}"
```
### Output
```bash
[ { "arguments": { "extended": "2.6.3 (isc20250522135511 deb)\npremium: yes (isc20250522135511 deb)\nlinked with:\n- log4cplus 2.0.8\n- OpenSSL 3.0.16 11 Feb 2025" }, "result": 0, "text": "2.6.3" } ]
```