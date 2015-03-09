package sockcore

import "io"

const (
	DATA_MAXLEN = 256 * 256 * 256
)

type DecodeReader interface {
	io.ByteReader
	io.Reader
}

type EncodeWriter interface {
	io.ByteWriter
	io.Writer
}
