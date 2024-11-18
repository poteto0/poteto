# Middleware

## middleware list

| name   | summary              |
| ------ | -------------------- |
| CORS   | CORS Policy          |
| Camara | Some Security Header |

## use middleware

```go
package main

import (
  "net/http"

  "github.com/poteto0/poteto"
  "github.com/poteto0/poteto/middleware"
)

func main() {
  p := poteto.New()

  p.Middleware(middleware.CORSWithConfig(CORSConfig{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete}}
  ))
  // handler
  ...
}
```
