# 0.x.x

### 0.23.0

- BUG: fix not allocated Server

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
