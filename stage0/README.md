# Stage 0: HTTP and JSON Service

This folder is the first backend project from the roadmap. It uses only Go's standard library so you can understand HTTP before adding a framework.

## Run it

From the project root:

```powershell
go run .\stage0
```

The server listens on `http://localhost:8080`. Stop it with `Ctrl+C`.

To use another port:

```powershell
$env:PORT = "9090"
go run .\stage0
```

## Endpoints

### Health and response headers

```powershell
curl.exe -i http://localhost:8080/health
```

This demonstrates a `GET` request, a `200 OK` status, JSON, and the custom `X-Learning-Server` response header.

### List quotes and query parameters

```powershell
curl.exe -i "http://localhost:8080/quotes?author=Unknown"
```

The `author` value is a query parameter. Query parameters are useful for filtering, searching, sorting, and pagination.

### Read one quote using a path parameter

```powershell
curl.exe -i http://localhost:8080/quotes/1
```

The `1` is a path parameter. The route is declared as `/quotes/{id}`.

### Create a quote with a JSON request body

```powershell
curl.exe -i -X POST http://localhost:8080/quotes `
  -H "Content-Type: application/json" `
  -d '{\"text\":\"Small steps compound.\",\"author\":\"Learner\"}'
```

This demonstrates:

- `POST` for creating a resource
- `Content-Type` request header
- JSON request body
- `201 Created` response
- `Location` response header

### Delete a quote

```powershell
curl.exe -i -X DELETE http://localhost:8080/quotes/1
```

`204 No Content` means the deletion succeeded and there is no response body.

### Inspect request metadata

```powershell
curl.exe -i "http://localhost:8080/request-info?name=Ayush" `
  -H "X-Client-Name: curl"
```

This endpoint demonstrates:

- The request method and URL
- Query parameters
- Request headers
- Cookies
- Content length
- `Set-Cookie` response headers

Call it twice to observe the cookie sent by the second request when using a client that stores cookies.

### Echo a plain-text request body

```powershell
curl.exe -i -X POST http://localhost:8080/echo `
  -H "Content-Type: text/plain" `
  -d "hello from the request body"
```

### Submit a form body

```powershell
curl.exe -i -X POST http://localhost:8080/form `
  -H "Content-Type: application/x-www-form-urlencoded" `
  -d "name=Ayush&topic=HTTP"
```

## HTTP concepts covered

### Methods

- `GET` reads a resource.
- `POST` creates a resource or submits data.
- `DELETE` removes a resource.
- The standard `ServeMux` returns `405 Method Not Allowed` for a known route used with an unsupported method.

`PUT` and `PATCH` are not needed by this first project, but you will use them when adding full update behavior in Stage 1.

### Status codes

- `200 OK`: the request succeeded and has a response body.
- `201 Created`: a new resource was created.
- `204 No Content`: the request succeeded without a response body.
- `400 Bad Request`: the client sent invalid input.
- `404 Not Found`: the requested resource does not exist.
- `405 Method Not Allowed`: the resource exists, but that HTTP method is unsupported.

### Headers

Request headers describe what the client sends or accepts, such as `Content-Type` and `X-Client-Name`.

Response headers describe what the server sends, such as `Content-Type`, `Location`, `Allow`, `Set-Cookie`, and `X-Learning-Server`.

Headers must be set before `WriteHeader` or before writing a body, because writing commits the response headers.

### Request bodies

The request body is a stream. The examples demonstrate:

- JSON decoding with `json.Decoder`
- Plain-text reading with `io.ReadAll`
- Form parsing with `Request.ParseForm`
- Body-size limits with `http.MaxBytesReader`

### Middleware

`loggingMiddleware` wraps the router and records the method, path, and duration for every request. Middleware is a function that receives a handler and returns another handler.

### Graceful shutdown

The server listens in a goroutine and waits for `Ctrl+C`. `Server.Shutdown` stops new requests and gives active requests a short time to finish.

## Test it

Run the unit and HTTP handler tests from the project root:

```powershell
go test .\stage0
go test -race .\stage0
```

The tests use `httptest`, so they do not need a real network port.
