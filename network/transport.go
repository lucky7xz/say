package network

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

type TransportKind string

const (
	TransportYggdrasil TransportKind = "ygg"
	TransportPlain     TransportKind = "plain"
)

func (k TransportKind) String() string {
	return string(k)
}

func ParseTransportKind(raw string) (TransportKind, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "ygg", "yggdrasil":
		return TransportYggdrasil, nil
	case "plain", "ip", "tcp", "tcpudp", "udp":
		return TransportPlain, nil
	default:
		return "", fmt.Errorf("unsupported transport %q", raw)
	}
}

type Transport interface {
	Kind() TransportKind
	LocalAddr() string
	WaitReady(ctx context.Context, timeout time.Duration) error
	ListenTCP(port int) (net.Listener, error)
	ListenUDP(port int) (net.PacketConn, error)
	DialTCP(host string, port int) (net.Conn, error)
	Close() error
}

func SetupTransport(kind TransportKind, verbose bool, cfgPath string) (Transport, error) {
	switch kind {
	case "", TransportYggdrasil:
		return newYggTransport(verbose, cfgPath)
	case TransportPlain:
		return newPlainTransport(verbose, cfgPath)
	default:
		return nil, fmt.Errorf("unsupported transport %q", kind)
	}
}
