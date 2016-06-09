package worker

import (
	natsc "github.com/nats-io/nats"
)

var nats *natsc.Conn

func MountNATS(n *natsc.Conn) {
	nats = n
}
