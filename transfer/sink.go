package transfer

import "io"

type Sink interface {
	io.Closer
	Process() (int64, error)
}
