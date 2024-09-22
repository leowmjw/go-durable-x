package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// This will represent the Zuplo API Gateway; assuming Cloudflare is up
// If issue; will use this raw Golang; need to find the correct tuning for handling
// 1MM continuous concurrent per 30 mins; with minimal latency due to this layer

// It might send calls to the Temporal Nexus Gateway as needed
// Authn + Authz happens on this layer
// Also Middleware processing; which might include a Rate limit ..

// It serves the static page out of the path static maps to folder of same name
// Static should have the Permanent cache; new files should use cache busting technique
// Unless there is a way to trigger from DiceDB static flag?

// In general should return no cache .. (what is the keyword?)

func main() {
	fmt.Println("API Gateway for Lead ..")
	Run()
}

func Run() {

	// Try to use what is available in the new routing .. --> https://go.dev/blog/routing-enhancements

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Define a handler function
	defaultHandler := func(w http.ResponseWriter, r *http.Request) {
		//http.Redirect(w, r, "/demo/debug/", http.StatusFound)
		fmt.Fprintf(w, "Hello World")
		return

	}

	// Attach handler function to the ServeMux
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Why u POST?")
		return
	})
	mux.HandleFunc("/", defaultHandler)
	// Read in Default Not Found page .. cross sell .. Search
	// Read in Default Error Page .. cross sell .. LLM + Customer Support ..
	// LLM can use Cloudflare AI or AWS Bedrock; have both available ..

	// HTTP Server Setup ..
	// Create the Server using the new ServeMux
	server := &http.Server{
		Addr:    ":8888",
		Handler: mux,
	}
	// Running the HTTP server in a go routine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("Server error:", err)
		}
	}()

	// Prepare for handling signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	interruptSignal := <-interrupt
	fmt.Printf("Received %s, shutting down.\n", interruptSignal)

	// Shutdown the server gracefully
	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Println("Error shutting down:", err)
	} else {
		fmt.Println("Server shutdown gracefully.")
	}

	return

}
