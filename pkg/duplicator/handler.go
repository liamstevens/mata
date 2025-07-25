package duplicator

import (
	"context"
	"io"
	"net"
	"sync"
	
	"github.com/liamstevens/mata/pkg/proxy"
	"github.com/liamstevens/mata/pkg/target"
)

// DuplicatingHandler decorates a ConnectionHandler to duplicate traffic to multiple targets
type DuplicatingHandler struct {
	selector target.TargetSelector
}

// NewDuplicatingHandler creates a new DuplicatingHandler
func NewDuplicatingHandler(selector target.TargetSelector) *DuplicatingHandler {
	return &DuplicatingHandler{
		selector: selector,
	}
}

// HandleConnection duplicates the connection to all targets
func (h *DuplicatingHandler) HandleConnection(ctx context.Context, clientConn net.Conn) error {
	defer clientConn.Close()
	
	// Connect to all targets
	targetConns, err := h.selector.Connect(ctx)
	if err != nil {
		return err
	}
	defer h.closeTargetConnections(targetConns)
	
	if len(targetConns) == 0 {
		return nil
	}
	
	// Create a MultiWriter to duplicate data to all targets
	writers := make([]io.Writer, len(targetConns))
	for i, conn := range targetConns {
		writers[i] = conn
	}
	multiWriter := io.MultiWriter(writers...)
	
	// Start goroutines for bidirectional copying
	var wg sync.WaitGroup
	errChan := make(chan error, len(targetConns)+1)
	
	// Copy from client to all targets
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := proxy.CopyWithContext(ctx, multiWriter, clientConn)
		errChan <- err
	}()
	
	// Copy from each target back to client (use first target as primary)
	if len(targetConns) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := proxy.CopyWithContext(ctx, clientConn, targetConns[0])
			errChan <- err
		}()
	}
	
	// Wait for completion or error
	go func() {
		wg.Wait()
		close(errChan)
	}()
	
	// Return first error encountered
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	
	return nil
}

// closeTargetConnections closes all target connections
func (h *DuplicatingHandler) closeTargetConnections(connections []net.Conn) {
	for _, conn := range connections {
		if conn != nil {
			conn.Close()
		}
	}
}