# 0.x.x

## 0.18.x

### 0.18.0

- `Poteto` has SetLogger
- You can call logger.XXX() from context

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
