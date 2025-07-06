# Integration Tests
These tests require a working Kea server on the same locahost.  

I've made https://github.com/rannday/kea-docker to spin up a Docker server with a very basic configuration that uses MySQL for the lease file, hosts storage, and configuration, something I need to figure out. ISC offers a paid hook for all the CRUD commands when using a SQL backend. I'm not paying for that.