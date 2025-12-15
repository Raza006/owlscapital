package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	host := getenv("POSTGRES_HOST", "postgres")
	port := getenv("POSTGRES_PORT", "5432")
	addr := net.JoinHostPort(host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	d := &net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck: cannot reach postgres at %s: %v\n", addr, err)
		os.Exit(1)
	}
	conn.Close()
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

