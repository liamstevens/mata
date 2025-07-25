package target

import (
	"context"
	"net"
	"testing"
)

func TestMultiTargetSelector(t *testing.T) {
	targets := []string{"localhost:8001", "localhost:8002", "localhost:8003"}
	selector := NewMultiTargetSelector(targets)
	
	// Test GetTargets
	result := selector.GetTargets()
	if len(result) != len(targets) {
		t.Errorf("Expected %d targets, got %d", len(targets), len(result))
	}
	
	for i, target := range targets {
		if result[i] != target {
			t.Errorf("Expected target %s at index %d, got %s", target, i, result[i])
		}
	}
}

func TestMultiTargetSelectorConnect_NoTargets(t *testing.T) {
	selector := NewMultiTargetSelector([]string{})
	
	ctx := context.Background()
	connections, err := selector.Connect(ctx)
	
	if err != nil {
		t.Errorf("Expected no error for empty targets, got %v", err)
	}
	
	if len(connections) != 0 {
		t.Errorf("Expected 0 connections, got %d", len(connections))
	}
}

func TestMultiTargetSelectorConnect_ValidTargets(t *testing.T) {
	// Create test servers
	listener1, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener1.Close()
	
	listener2, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener2.Close()
	
	// Start echo servers
	go startEchoServer(listener1)
	go startEchoServer(listener2)
	
	// Test connection to both servers
	targets := []string{listener1.Addr().String(), listener2.Addr().String()}
	selector := NewMultiTargetSelector(targets)
	
	ctx := context.Background()
	connections, err := selector.Connect(ctx)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if len(connections) != 2 {
		t.Errorf("Expected 2 connections, got %d", len(connections))
	}
	
	// Clean up connections
	for _, conn := range connections {
		if conn != nil {
			conn.Close()
		}
	}
}

func TestMultiTargetSelectorConnect_PartialFailure(t *testing.T) {
	// Create one valid server
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	
	go startEchoServer(listener)
	
	// Mix valid and invalid targets
	targets := []string{listener.Addr().String(), "localhost:99999"}
	selector := NewMultiTargetSelector(targets)
	
	ctx := context.Background()
	connections, err := selector.Connect(ctx)
	
	// Should fail due to second target being invalid
	if err == nil {
		t.Error("Expected error for invalid target, got none")
	}
	
	if connections != nil {
		t.Error("Expected nil connections on error")
	}
}

func TestMultiTargetSelectorClose(t *testing.T) {
	selector := NewMultiTargetSelector([]string{"localhost:8001"})
	
	err := selector.Close()
	if err != nil {
		t.Errorf("Expected no error from Close, got %v", err)
	}
}

// Helper function to start a simple echo server
func startEchoServer(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		go func(c net.Conn) {
			defer c.Close()
			// Just accept and close - we're testing connections, not data transfer
		}(conn)
	}
}