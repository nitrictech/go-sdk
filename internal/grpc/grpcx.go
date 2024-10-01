package grpcx

import (
	"sync"

	"github.com/nitrictech/go-sdk/constants"
	"google.golang.org/grpc"
)

type grpcManager struct {
	conn      grpc.ClientConnInterface
	connMutex sync.Mutex
}

var m = grpcManager{
	conn:      nil,
	connMutex: sync.Mutex{},
}

func GetConnection() (grpc.ClientConnInterface, error) {
	m.connMutex.Lock()
	defer m.connMutex.Unlock()

	if m.conn == nil {
		conn, err := grpc.NewClient(constants.NitricAddress(), constants.DefaultOptions()...)
		if err != nil {
			return nil, err
		}
		m.conn = conn
	}

	return m.conn, nil
}
