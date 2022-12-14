<h1>Error Responses</h1>

- [Overview](#overview)
- [Types and Interfaces](#types-and-interfaces)
  - [CustomError](#customerror)
  - [ErrorWriter](#errorwriter)
  - [HTTPError](#httperror)
  - [WrappedError](#wrappederror)
- [Callbacks](#callbacks)
  - [MapError](#maperror)
  - [WriteError](#writeerror)


# Overview 

By default, error responses will all look something along the lines of:
```json
{
    "status":
      {
          "code": "321",
          "description": "Error Description"
      }
}
```

There are a number of different ways to modify how the errors are returned listed below (the most common being [MapError](#MapError)).

# Types and Interfaces

## CustomError
```go
type CustomError map[string]string
```
Will map to an [HTTPError](#HTTPError) based on the keys in the map:
- http_status - will run `strconv.Atoi` on the value and use it as the `HTTPCode` member
- http_code - maps to the `Code` member
- http_message - maps to the `Description` member

Custom error types can be automatically generated by annotating error types with the `error` tag, for example, this sysl:
```sysl
    !type UnauthorizedError [~error]:
        http_status <: string [value = "401"]
        http_code <: string [value = "1003"]
        http_message <: string [value = "Unauthorized"]
```
Will generate this type:
```go
var UnauthorizedError common.CustomError = map[string]string{
	"name":         "UnauthorizedError",
	"http_code":    "1003",
	"http_message": "Unauthorized",
	"http_status":  "401",
}
```

## ErrorWriter
```go
// If the error returned is an ErrorWriter the error handling will call the writeError method before any of the
// regular error handling (no mapping).
//
// If the call returns true it means it wrote the error and will not do any more handling.
// If it returns false it will go through the normal error writing path (via both the MapError and WriteError callbacks).
type ErrorWriter interface {
	WriteError(ctx context.Context, w http.ResponseWriter) bool
}
```
If the error returned from a handler is an `ErrorWriter` then it will call the `WriteError` method directly instead of the normal error mapping path. Note this means that if `true` is returned from the call, then none of the other types or callbacks will be used or called for this response.

## HTTPError 
```go
type HTTPError struct {
	HTTPCode    int    `json:"-"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
```
The normal error response path will map all errors into an `HTTPError` which will then be converted to json (see [Overview](#Overview) for an example) and returned under the `status` attribute with the status code taken from the `HTTPCode` field. 

## WrappedError
*Note this type is deprecated, prefer using a different method, for example you could call `AddField` directly on the [HTTPError](#HTTPError) or use an [ErrorWriter](#errorwriter) or [WriteError](#writeerror) to get full access to the response writing.
```go
func WrappedError(err error, fields ...KV) error {
	return wrappedError{
		e:      err,
		fields: fields,
	}
}

type KV struct {
    K string
    V interface{}
}
```
You can return the value of `WrappedError` to add additional fields to the response json message. `K` will be used as attribute and `json.Marshal` will be called on `V` to produce the value.

# Callbacks

## MapError

```go
// MapError maps an error to an HTTPError in instances where custom error mapping is required.
// Return nil to perform default error mapping; defined as:
// 1. CustomError.HTTPError if the original error is a CustomError, otherwise
// 2. common.MapError
MapError func(ctx context.Context, err error) *common.HTTPError
```
All errors (both those produced by sysl-go library and by the handlers) will go through this mapping to map them to an [HTTPError](#HTTPError) before being written to the `http.ResponseWriter`.

This is the most common callback used for manipulating the response values before they sent back.

For example creating a new server would look something like:
```go
func newAppServer(ctx context.Context) (core.StoppableServer, error) {
	return example.NewServer(ctx, func(ctx context.Context, cfg AppConfig) (*example.ServiceInterface, *core.Hooks, error) {
		return &example.ServiceInterface{
				//...
			}, &core.Hooks{
				MapError: func(ctx context.Context, err error) *common.HTTPError {
					return &common.HTTPError{
						// fill details here based on err
					}
				},
			},
			nil
	})
}
```

## WriteError
```go
// WriteError can be used to write the error to the writer in whatever way you want.
// If not supplied it will use httpError.WriteError as the default.
WriteError func(ctx context.Context, w http.ResponseWriter, httpError *common.HTTPError)
```
If you need to modify the format that the responses are written, you can supply a `WriteError` callback that will be called instead of `httpError.WriteError`.

For example creating a new server would look something like:
```go
func newAppServer(ctx context.Context) (core.StoppableServer, error) {
	return example.NewServer(ctx, func(ctx context.Context, cfg AppConfig) (*example.ServiceInterface, *core.Hooks, error) {
		return &example.ServiceInterface{
				//...
			}, &core.Hooks{
				WriteError: func(ctx context.Context, w http.ResponseWriter, httpError *common.HTTPError) {
					// These next few lines are just an example of what could be done
					w.Header().Set("Content-Type", "application/json;charset=UTF-8")
					w.WriteHeader(httpError.HTTPCode)
					w.Write([]byte{
						// fill details here based on httpError
					})
				},
			},
			nil
	})
}
```
