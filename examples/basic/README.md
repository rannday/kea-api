# Usage Examples

```go
// HTTP-based client (recommended)
import (
    "github.com/rannday/kea-api/client"
    "github.com/rannday/kea-api/agent"
)

c := client.NewHTTPTransport("http://localhost:8000")
a := agent.NewWithClient(client.NewClient(c))
```

```go
// HTTP-based client with basic auth
import (
    "github.com/rannday/kea-api/client"
    "github.com/rannday/kea-api/agent"
)

c := client.NewHTTPTransport(
    "http://localhost:8000",
    client.WithAuth(client.NewBasicAuth("admin", "secret")),
)
a := agent.NewWithClient(client.NewClient(c))
```

```go
// UNIX socket client
import (
    "time"
    "github.com/rannday/kea-api/client"
    "github.com/rannday/kea-api/agent"
)

t := client.NewSocketTransport("unix", "/run/kea/kea4-ctrl-socket", 5*time.Second)
a := agent.NewWithClient(client.NewClient(t))
```

```go
// TCP socket client
import (
    "time"
    "github.com/rannday/kea-api/client"
    "github.com/rannday/kea-api/agent"
)

t := client.NewSocketTransport("tcp", "127.0.0.1:8001", 5*time.Second)
a := agent.NewWithClient(client.NewClient(t))
```

```go
// Example: get Kea status
status, err := a.StatusGet()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Kea status: %+v\n", status)
```
