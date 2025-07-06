# kea-api

Golang ISC Kea API Module

Provides the means to interact with the [ISC Kea DHCP](https://kea.readthedocs.io/en/latest/) API.

## Features

- Native Go structs and types for Kea API commands
- Easy-to-use functions for common operations like `lease4-get`, `lease4-del`, `config-get`, etc.

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

## License

This project is licensed under [The Unlicense](https://unlicense.org/)