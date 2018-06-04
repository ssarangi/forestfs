package forestfs

import (
	"context"
	"io"
)

type Request struct {
	Ctx     context.Context
	Conn    io.ReadWriter
	Header  *protocol.RequestHeader
	Request interface{}
}

type Response struct {
	Ctx      context.Context
	Conn     io.ReadWriter
	Header   *protocol.RequestHeader
	Response interface{}
}
