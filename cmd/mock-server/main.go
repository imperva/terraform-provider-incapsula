package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/terraform-providers/terraform-provider-incapsula/incapsula"
)

func main() {
	port := flag.Int("port", 19443, "Port to listen on")
	flag.Parse()

	// Create the mock server
	mock := incapsula.NewMockImpervaServer()

	// Create a new server on the specified port instead of using the httptest server
	mock.Server.Close() // Close the auto-started httptest server

	addr := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(mock.ServeHTTP),
	}

	// Handle shutdown gracefully
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down mock server...")
		done <- true
	}()

	log.Printf("Mock Imperva API server starting on http://localhost%s", addr)
	log.Println("")
	log.Println("Environment variables to use with Terraform:")
	log.Println("")
	fmt.Printf("export INCAPSULA_API_ID=mock-api-id\n")
	fmt.Printf("export INCAPSULA_API_KEY=mock-api-key\n")
	fmt.Printf("export INCAPSULA_BASE_URL=http://localhost%s\n", addr)
	fmt.Printf("export INCAPSULA_BASE_URL_REV_2=http://localhost%s\n", addr)
	fmt.Printf("export INCAPSULA_BASE_URL_REV_3=http://localhost%s\n", addr)
	fmt.Printf("export INCAPSULA_BASE_URL_API=http://localhost%s\n", addr)
	fmt.Printf("export INCAPSULA_CUSTOM_TEST_DOMAIN=.mock.incaptest.com\n")
	log.Println("")
	log.Println("Press Ctrl+C to stop the server")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-done
	log.Println("Server stopped")
}
