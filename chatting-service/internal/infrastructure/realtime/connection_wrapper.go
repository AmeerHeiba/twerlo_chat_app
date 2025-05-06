package realtime

import (
	"errors"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

// ConnectionWrapper bridges fiber/websocket and our notifier
type ConnectionWrapper struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (w *ConnectionWrapper) WriteJSON(v interface{}) error {
	if w == nil || w.conn == nil {
		return errors.New("nil connection")
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteJSON(v)
}

func (w *ConnectionWrapper) ReadMessage() (int, []byte, error) {
	return w.conn.ReadMessage()
}

func (w *ConnectionWrapper) Close() error {
	return w.conn.Close()
}
