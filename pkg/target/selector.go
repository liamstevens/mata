package target

import (
	"context"
	"net"
)

// TargetSelector defines the interface for selecting and managing target connections
type TargetSelector interface {
	// GetTargets returns the list of target addresses
	GetTargets() []string
	
	// Connect establishes connections to all targets
	Connect(ctx context.Context) ([]net.Conn, error)
	
	// Close closes all managed connections
	Close() error
}

// MultiTargetSelector implements TargetSelector for multiple targets
type MultiTargetSelector struct {
	targets []string
}

// NewMultiTargetSelector creates a new MultiTargetSelector
func NewMultiTargetSelector(targets []string) *MultiTargetSelector {
	return &MultiTargetSelector{
		targets: targets,
	}
}

// GetTargets returns the list of target addresses
func (s *MultiTargetSelector) GetTargets() []string {
	return s.targets
}

// Connect establishes connections to all targets
func (s *MultiTargetSelector) Connect(ctx context.Context) ([]net.Conn, error) {
	connections := make([]net.Conn, 0, len(s.targets))
	
	for _, target := range s.targets {
		conn, err := net.Dial("tcp", target)
		if err != nil {
			// Close any existing connections on error
			s.closeConnections(connections)
			return nil, err
		}
		connections = append(connections, conn)
	}
	
	return connections, nil
}

// Close closes all managed connections
func (s *MultiTargetSelector) Close() error {
	return nil
}

// closeConnections is a helper to close a slice of connections
func (s *MultiTargetSelector) closeConnections(connections []net.Conn) {
	for _, conn := range connections {
		if conn != nil {
			conn.Close()
		}
	}
}