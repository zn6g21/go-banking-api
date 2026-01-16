package pkg

import (
	"net"
	"net/url"
	"os"
	"time"
)

func GetEndpoint(path string) string {
	var baseURL string
	baseURL = "http://localhost:8080"
	env := os.Getenv("APP_ENV")
	if env == "stage" {
		baseURL = "https://stage.localhost:8080"
	}
	p, _ := url.Parse(path)
	b, _ := url.Parse(baseURL)
	return b.ResolveReference(p).String()
}

func CheckPort(host string, port string) bool {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return true
	}
	if closeErr := conn.Close(); closeErr != nil {
		return false
	}
	return false
}

func WaitForPort(host string, port string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if CheckPort(host, port) {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}
