/**
 * =============================================================================
 * Token Manager - Main Entry Point
 *
 * This file serves as the entry point for the Token Manager application.
 * It initializes the server, sets up the Gin router, and starts the HTTP server.
 *
 * Table of Contents:
 * 1. Package Declaration
 * 2. Import Statements
 * 3. Main Function
 *
 * Author: Token Manager Team
 * Version: 1.0.0
 * =============================================================================
 */

package main

/**
 * =============================================================================
 * Import Statements
 *
 * Standard Library:
 * - log: For logging messages and fatal errors
 *
 * External Packages:
 * - github.com/gin-gonic/gin: Web framework for building HTTP services
 * =============================================================================
 */

import (
	"log"

	"github.com/gin-gonic/gin"
)

/**
 * =============================================================================
 * Main Function
 *
 * The main entry point of the application. It performs the following:
 * 1. Creates a new Store instance for data persistence
 * 2. Initializes the Gin router with default middleware
 * 3. Registers all API routes
 * 4. Starts the HTTP server on port 8080
 *
 * Error Handling:
 * - If store creation fails, logs fatal error and exits
 * - If server fails to start, logs fatal error and exits
 * =============================================================================
 */

func main() {
	// Create a new store instance for data persistence
	// The store will manage employees, tokens, and usage records
	// It will automatically load existing data from data.json or create a new file
	store, err := NewStore("data.json")
	if err != nil {
		// Log fatal error and exit if store creation fails
		log.Fatalf("Failed to create store: %v", err)
	}

	// Initialize Gin router with default middleware
	// Default middleware includes:
	// - Logger: Logs incoming HTTP requests
	// - Recovery: Recovers from any panics during request handling
	r := gin.Default()

	// Register all API routes with the router
	// This sets up all endpoints for employees, tokens, and statistics
	store.RegisterRoutes(r)

	// Log server startup message
	log.Println("Server starting on :8080")

	// Start the HTTP server and block until it stops
	// The server listens on all network interfaces (0.0.0.0) on port 8080
	if err := r.Run(":8080"); err != nil {
		// Log fatal error and exit if server fails to start
		log.Fatalf("Failed to start server: %v", err)
	}
}
