# kea-api
Golang [ISC Kea DHCP](https://kea.readthedocs.io/en/latest/) API Module  
## Features
- Native Go structs and types for the Kea API
- Easy-to-use functions for all API commands
## Import & Usage
Install the module:  
```bash
go get github.com/rannday/kea-api
```
[Basic Usage](examples/basic/README.md)
### Examples
See the [`examples/`](./examples) directory for usage samples:
- [`examples/basic`](./examples/basic): minimal usage demo
- [`examples/custom_client`](./examples/custom_client): more advanced/customized client usage
## TODO
- Add coverage for more Kea API commands
- Add support for calling multiple services at once for shared commands
## License
This project is licensed under [The Unlicense](https://unlicense.org/)