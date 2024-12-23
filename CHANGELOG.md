# 0.x.x

## 0.26.x

### 0.26.2

- TEST: add benchmark by @poteto0 in #164
- Build(deps): bump github.com/goccy/go-yaml from 1.15.10 to 1.15.13 in #163

### 0.26.1

- EXAMPLE: add example on jsonrpc by @poteto0 in #160
- EXAMPLE: add example on fast-api by @poteto0 in #160
- EXAMPLE: add example on api by @poteto0 in #160
- BUG: fix `PotetoJSONRPCAdapter` dosen't check class by @poteto0 in #158

## 0.26.0

- FEATURE: `PotetoJSONAdapter` provides json rpc by @poteto0 in #154

KeyNote: You can serve JSONRPC server easily.

```go
type (
  Calculator struct{}
  AdditionArgs   struct {
    Add, Added int
  }
)

func (tc *TestCalculator) Add(r *http.Request, args *AdditionArgs) int {
 return args.Add + args.Added
}

func main() {
  p := poteto.New()

  rpc := TestCalculator{}
  // you can access "/add/Calculator.Add"
  p.POST("/add", func(ctx poteto.Context) error {
    return poteto.PotetoJsonRPCAdapter[Calculator, AdditionArgs](ctx, &rpc)
  })

  p.Run("8080")
}
```

- FEATURE: `Poteto.RunTLS` serve https server provided cert & key file by @poteto0 in #144

## 0.25.x

### 0.25.3

- Bump github.com/goccy/go-json from 0.10.3 to 0.10.4 in #152
- Bump github.com/goccy/go-yaml from 1.15.7 to 1.15.10 in #151

### 0.25.2

- BUG: fix ut ignore path by @poteto0 in #147
- OP: test on go@1.21.x, go@1.22.x, go@1.23.x by @poteto0 in #147
- OP: only go@1.23.x upload to codecov by @poteto0 in #148

### 0.25.1

- TEST: ut progress by @poteto0 in #141
- CHANGE: appropriate error messages by @poteto0 in #140
- TEST: ut progress by @poteto0 in #136

### 0.25.0

- FEATURE: poteto-cli released by @poteto0 in #133
- DEPENDENCY: Bump github.com/goccy/go-yaml from 1.15.5 to 1.15.7 in #134

## 0.24.x

### 0.24.0

- FEATURE: mid param router ex /users/:id/name by @poteto0 in #122
- REFACTOR: some switch case by @poteto0 in #122
- FEATURE: ctx.DebugParam by @poteto0 in #125
- OPTIMIZE: middlewareTree by @eaggle23 in #131

## 0.23.x

### 0.23.4

- OPTIMIZE: performance tuning by @poteto0 in #116
- OPTIMIZE: performance tuning by @poteto0 in #117

### 0.23.3

- BUG: fix "/" routes nothing by @poteto0 in #112

### 0.23.2

- OPTIMIZE: optimize router's structure & faster by @poteto0 in #109
- FEATURE: Now, poteto follows patch, head, options, trace, and connect by @poteto0 in #109
- DOCUMENT: update some document by @poteto0 in #109

### 0.23.1

- DOCUMENT: add example app by @poteto0 #104

### 0.23.0

- BUG: fix not allocated Server by @poteto0 #101

## 0.22.x: has critical bug

### 0.22.0

- FEATURE: `Context.RealIP()` return realIp
- CHANGE: `Context.GetIPFromXFFHeader()` return just X-Forwarded-For
- DOCUMENT: update some document

## 0.21.x: has critical bug

### 0.21.0

- FEATURE: `Poteto.Leaf(path, handler)` make router great
- DOCUMENT: Update some document

## 0.20.x: has critical bug

### 0.20.0

- CHANGE: `Poteto.Run()` internal call http.Server#Serve instead of http.ListenAndServe
  You can use your protocol such as udp
- CHANGE: `Poteto.Stop(stdContext)` stop server

## 0.19.x

### 0.19.1

- `PotetoOption`: you can make WithRequestId false
  Because it is slowly With RequestId. If you don't need this, you can make app faster
- fix bug
- refactor something of private func

### 0.19.0

- `Context.Get(key)` get value by key from store.
- `Context.RequestId()` get requestId from Header or store or generate uuid@v4
- `Poteto.ServeHTTP(r, w)` call requestId and set to Header.
  - It may be to become middle ware

## 0.18.x

### 0.18.1

- Fix bug of first msg
- optimize bit

### 0.18.0

- `Poteto` has SetLogger
- You can call ctx.Logger().XXX() from context

## 0.17.x

## 0.17.2

- `Poteto.Run()` will now also accept mere numbers. For example, `8080` is converted to `:8080` and processed.
- Poteto logged "http://localhost:<port>"

### 0.17.1

- warning handler collision detect

### 0.17.0

- timeout middleware
- poteto.Response.writer became public member

## 0.16.x

### 0.16.1

- fix bug
  - become: `Context.QueryParam()` & `Context.PathParam()` only return string
