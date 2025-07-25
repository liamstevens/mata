package proxy

import (
	"context"
	"net"
)

// ConnectionHandler defines the interface for handling network connections
type ConnectionHandler interface {
	// HandleConnection processes an incoming network connection
	HandleConnection(ctx context.Context, conn net.Conn) error
}

// BasicHandler implements a simple pass-through connection handler
type BasicHandler struct {
	target string
}

// NewBasicHandler creates a new BasicHandler that forwards to the specified target
func NewBasicHandler(target string) *BasicHandler {
	return &BasicHandler{
		target: target,
	}
}

// HandleConnection forwards the connection to the target
func (h *BasicHandler) HandleConnection(ctx context.Context, conn net.Conn) error {
	defer conn.Close()

	// Connect to target
	targetConn, err := net.Dial("tcp", h.target)
	if err != nil {
		return err
	}
	defer targetConn.Close()

	// Start bidirectional copying
	errChan := make(chan error, 2)
	
	// Copy from client to target
	go func() {
		_, err := CopyWithContext(ctx, targetConn, conn)
		errChan <- err
	}()
	
	// Copy from target to client
	go func() {
		_, err := CopyWithContext(ctx, conn, targetConn)
		errChan <- err
	}()

	// Wait for first error or context cancellation
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}