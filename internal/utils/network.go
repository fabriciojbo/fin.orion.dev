package utils

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// CheckPort verifica se uma porta está aberta
func CheckPort(port int) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", fmt.Sprintf("%d", port)), 5*time.Second)
	if err != nil {
		return false
	}
	defer func() { _ = conn.Close() }()
	return true
}

// CheckPortString verifica se uma porta está aberta usando string
func CheckPortString(port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", port), 5*time.Second)
	if err != nil {
		return false
	}
	defer func() { _ = conn.Close() }()
	return true
}

// CheckHTTPEndpoint verifica se um endpoint HTTP está respondendo
func CheckHTTPEndpoint(url string) bool {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()
	return resp.StatusCode < 500
}

// CheckHTTPEndpointWithTimeout verifica se um endpoint HTTP está respondendo com timeout customizado
func CheckHTTPEndpointWithTimeout(url string, timeout time.Duration) bool {
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()
	return resp.StatusCode < 500
}
