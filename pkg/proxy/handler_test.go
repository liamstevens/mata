package proxy

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestBasicHandler(t *testing.T) {
	// Start a simple echo server as the target
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	
	targetAddr := listener.Addr().String()
	
	// Start echo server goroutine
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				buffer := make([]byte, 1024)
				n, _ := c.Read(buffer)
				_, _ = c.Write(buffer[:n])
			}(conn)
		}
	}()
	
	// Create BasicHandler
	handler := NewBasicHandler(targetAddr)
	
	// Create client and server connections
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()
	
	// Handle connection in goroutine
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	go func() {
		_ = handler.HandleConnection(ctx, serverConn)
	}()
	
	// Send test data
	testData := "hello world"
	_, _ = clientConn.Write([]byte(testData))
	
	// Read response
	buffer := make([]byte, 1024)
	_ = clientConn.SetReadDeadline(time.Now().Add(time.Second))
	n, err := clientConn.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}
	
	response := string(buffer[:n])
	if response != testData {
		t.Errorf("Expected %q, got %q", testData, response)
	}
}