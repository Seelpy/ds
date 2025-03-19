package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)


type Reader struct {
	conn nats.Conn
}

func (r *Reader) Connect() error {
	r.conn.
}

func (r *Reader) Close() error {
	r.conn.Close()
}