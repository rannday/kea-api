# Development in Windows
## Requirements
- Install Go - `winget install GoLang.Go`
- Install Docker - `winget install Docker.DockerCLI`
- Maybe Install Desktop App - `winget install Docker.DockerDesktop`
#### Docker Setup
Daemon configuration
```
{
  "builder": {
    "gc": {
      "defaultKeepStorage": "20GB",
      "enabled": true
    }
  },
  "default-address-pools": [
    {
      "base": "192.168.68.0/22",
      "size": 26
    }
  ],
  "experimental": false
}
```