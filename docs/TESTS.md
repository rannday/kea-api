# Tests

The project includes unit tests. To run tests:

```bash
make test
```
```bash
test.bat
```
## Integration Tests
To run integration tests against a docker container  
```bash
test.bat -i
```
https://github.com/rannday/kea-docker  

Very basic for now. Want to eventually have it setup a mock dhcp server with dummy subnets, all configured through SQL. Maybe even spin up two and play with HA.

## Integration Tests
These tests require a working Kea server on the same locahost.  

I've made https://github.com/rannday/kea-docker to spin up a Docker server with a very basic configuration that uses MySQL for the lease file, hosts storage, and configuration, something I need to figure out. ISC offers a paid hook for all the CRUD commands when using a SQL backend. I'm not paying for that.

## TODO
- TLS tests - Would have to create the certs via docker, or maybe before hand and just include them. Seems annoying.
- DDNS - No idea where to begin on this one, but I'd like to figure it out.