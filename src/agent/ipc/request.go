package ipc

const (
	UNKNOWN = iota
	USERINFO_REQUEST
)

type RequestType struct {
	Code int16			// tos
	CH chan interface{} // service-oriented data channel
}
