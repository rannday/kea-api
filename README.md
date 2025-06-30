# isc-kea

ISC Kea API Golang Library

Provides the means to interact with the [ISC Kea DHCP](https://kea.readthedocs.io/en/latest/) API.

## Features

- Native Go structs and types for Kea API commands
- Easy-to-use functions for common operations like `lease4-get`, `lease4-del`, `config-get`, etc.

## Examples

See the [`examples/`](./examples) directory for usage samples:
- [`examples/basic`](./examples/basic): minimal usage demo
- [`examples/custom_client`](./examples/custom_client): more advanced/customized client usage

## Tests

The project includes unit tests. To run tests:

```bash
make test
```
```bash
test.bat
```

## TODO
- Add coverage for more Kea API commands

## License

This project is licensed under [The Unlicense](https://unlicense.org/)