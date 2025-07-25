package proxy

import (
	"context"
	"io"
)

// CopyWithContext copies from src to dst with context cancellation support
func CopyWithContext(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	// Use a channel to signal completion
	type result struct {
		n   int64
		err error
	}
	
	resultChan := make(chan result, 1)
	
	go func() {
		n, err := io.Copy(dst, src)
		resultChan <- result{n, err}
	}()
	
	select {
	case res := <-resultChan:
		return res.n, res.err
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}