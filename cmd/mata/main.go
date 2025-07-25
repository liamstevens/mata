package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	
	"github.com/liamstevens/mata/pkg/proxy"
	"github.com/liamstevens/mata/pkg/duplicator"
	"github.com/liamstevens/mata/pkg/target"
)

type Config struct {
	Source  string
	Targets []string
	Echo    bool
}

func main() {
	config := parseArgs()
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()
	
	// Create handler based on configuration
	var handler proxy.ConnectionHandler
	if config.Echo {
		handler = &EchoHandler{}
	} else if len(config.Targets) == 1 {
		handler = proxy.NewBasicHandler(config.Targets[0])
	} else {
		// Multiple targets - use duplicating handler
		selector := target.NewMultiTargetSelector(config.Targets)
		handler = duplicator.NewDuplicatingHandler(selector)
	}
	
	// Start proxy server
	if err := startProxy(ctx, config.Source, handler); err != nil {
		log.Fatal(err)
	}
}

func parseArgs() *Config {
	var (
		source  = flag.String("source", "", "Source port to listen on (e.g., :8080)")
		targets = flag.String("targets", "", "Comma-separated list of target addresses")
		echo    = flag.Bool("echo", false, "Echo mode - return received data back to client")
	)
	flag.Parse()
	
	if *source == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -source PORT [options]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	
	config := &Config{
		Source: *source,
		Echo:   *echo,
	}
	
	if *targets != "" {
		config.Targets = strings.Split(*targets, ",")
		for i, target := range config.Targets {
			config.Targets[i] = strings.TrimSpace(target)
		}
	}
	
	if !config.Echo && len(config.Targets) == 0 {
		log.Fatal("Must specify either -echo or -targets")
	}
	
	return config
}

func startProxy(ctx context.Context, address string, handler proxy.ConnectionHandler) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}
	defer listener.Close()
	
	log.Printf("Mata proxy listening on %s", address)
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("Accept error: %v", err)
			continue
		}
		
		go func(c net.Conn) {
			if err := handler.HandleConnection(ctx, c); err != nil {
				log.Printf("Connection error: %v", err)
			}
		}(conn)
	}
}

// EchoHandler implements a simple echo server for testing
type EchoHandler struct{}

func (h *EchoHandler) HandleConnection(ctx context.Context, conn net.Conn) error {
	defer conn.Close()
	
	_, err := proxy.CopyWithContext(ctx, conn, conn)
	return err
}