package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const DEFAULT_PORT = 80
const CONF_PATH = "/home/root/.config/tailscale-rm-webint/config"
const WEBINT_URL = "http://10.11.99.1:80"

func main() {
	port := portFromConfig()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go startHTTPServer(ctx, &wg, port)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
	fmt.Println("\nGracefully shutting down HTTP server...")

	cancel()
	wg.Wait()
	fmt.Println("Shutdown complete")
}

func startHTTPServer(ctx context.Context, wg *sync.WaitGroup, port int) {
	defer wg.Done()

	remote, err := url.Parse(WEBINT_URL)
	if err != nil {
		panic(err)
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			r.Host = remote.Host
			p.ServeHTTP(w, r)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc("/", handler(proxy))
	server := &http.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", port),
	}

	go func() {
		fmt.Println("Starting HTTP server...")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %s\n", err)
			fmt.Println("Is the daemon already running?")
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelShutdown()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		fmt.Printf("HTTP server shutdown error: %s\n", err)
	}
}

func portFromConfig() int {
	file, err := os.Open(CONF_PATH)
	if err != nil {
		fmt.Println("Config not found at:", CONF_PATH)
		fmt.Println("Using default port:", DEFAULT_PORT)
		return DEFAULT_PORT
	}
	fmt.Println("Config found at:", CONF_PATH)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key != "port" {
			continue
		}

		port, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Error parsing port value:", value)
			break
		}

		fmt.Println("Using custom port:", port)
		return port
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading config:", err)
	}

	fmt.Println("No valid port found. EX: port=8080")
	fmt.Println("Using default port:", DEFAULT_PORT)
	return DEFAULT_PORT
}
