package http

import (
	"GoFeed/internal/generated/api/go_feed"
	"fmt"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClientPool struct {
	mu       sync.Mutex
	clients  []*grpc.ClientConn
	currIdx  int
	poolSize int
}

// NewClientPool creates a new gRPC connection pool
func NewGrpcClientPool(serverAddress string, poolSize int) (*grpcClientPool, error) {
	if poolSize <= 0 {
		return nil, fmt.Errorf("invalid pool size: %d", poolSize)
	}

	clients := make([]*grpc.ClientConn, 0, poolSize)
	for i := 0; i < poolSize; i++ {
		conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %w", err)
		}
		clients = append(clients, conn)
	}
	return &grpcClientPool{
		clients:  clients,
		poolSize: poolSize,
	}, nil
}

// GetClient returns a gRPC client connection from the pool (round-robin)
func (p *grpcClientPool) GetClient() (go_feed.GoFeedServiceClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Ensure there are clients in the pool
	if len(p.clients) == 0 {
		return nil, fmt.Errorf("no gRPC connections available in the pool")
	}

	client := p.clients[p.currIdx]
	p.currIdx = (p.currIdx + 1) % len(p.clients) // round-robin logic

	return go_feed.NewGoFeedServiceClient(client), nil
}

// Close closes all gRPC connections in the pool
func (p *grpcClientPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, client := range p.clients {
		if err := client.Close(); err != nil {
			log.Printf("failed to close gRPC connection: %v", err)
		}
	}

	p.clients = nil
}
