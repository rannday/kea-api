# Client Constructor Examples

## HTTP-based client
```go
import (
    "github.com/rannday/kea-api/client"
)

c := client.NewHTTP("http://localhost:8000")
```

## HTTP-based client with basic auth
```go
import (
    "github.com/rannday/kea-api/client"
)

c := client.NewHTTPTransport(
    "http://localhost:8000",
    client.WithAuth(client.NewBasicAuth("admin", "secret")),
)
```

## UNIX socket client
```go
import (
    "time"
    "github.com/rannday/kea-api/client
)

t := client.NewSocket("unix", "/run/kea/kea4-ctrl-socket", 5*time.Second)
```

## TCP socket client
```go
import (
    "time"
    "github.com/rannday/kea-api/client"
)

t := client.NewSocket("tcp", "127.0.0.1:8001", 5*time.Second)
```
