package ctxrouter


//ErrorInterface You can custom any error structure you want
type ErrorInterface interface {
	StatusCode() int
}
