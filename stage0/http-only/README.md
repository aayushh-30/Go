# HTTP Only

This is a deliberately small, separate program for learning only the `net/http` module. It has no quote model, store, repository, database, or business workflow.

## Run

From the project root:

```powershell
go run .\stage0\http-only
```

The server listens on `http://localhost:8081`.

## Try the endpoints

### Create a server response

```powershell
curl.exe -i http://localhost:8081/
```

### Inspect an incoming request

```powershell
curl.exe -i "http://localhost:8081/inspect?name=Ayush" `
  -H "X-Client-Name: curl"
```

This shows the method, URL, path, query parameter, request header, and content length.

### Accept a JSON request body

```powershell
curl.exe -i -X POST http://localhost:8081/message `
  -H "Content-Type: application/json" `
  -d '{\"text\":\"hello HTTP\"}'
```

This demonstrates a request body, JSON decoding, `201 Created`, and the `Location` response header.

### Accept a plain-text request body

```powershell
curl.exe -i -X POST http://localhost:8081/echo `
  -H "Content-Type: text/plain" `
  -d "hello from the request body"
```

### See status codes

```powershell
curl.exe -i "http://localhost:8081/status?code=204"
curl.exe -i "http://localhost:8081/status?code=400"
curl.exe -i "http://localhost:8081/status?code=404"
```

## What belongs in this lesson

- Creating `http.Server`
- Registering routes with `http.NewServeMux`
- Accepting requests with `ListenAndServe`
- Reading `http.Request`
- Reading methods, paths, query parameters, headers, and bodies
- Writing `http.ResponseWriter`
- Setting response headers before the status or body
- Returning status codes
- JSON and plain-text responses
- Limiting request body sizes
- Graceful shutdown with `Server.Shutdown`

Run `Ctrl+C` to stop the server gracefully.
