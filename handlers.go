/**
 * =============================================================================
 * Token Manager - HTTP Handlers
 *
 * This file contains all HTTP request handlers for the API endpoints.
 * Each handler processes incoming requests, calls the appropriate Store
 * methods, and returns responses to clients.
 *
 * Table of Contents:
 * 1. Package Declaration
 * 2. Import Statements
 * 3. Route Registration
 * 4. Employee Handlers
 * 5. Token Handlers
 * 6. Statistics Handlers
 * 7. Mapping Handlers
 *
 * Handler Patterns:
 * - Each handler receives a Gin Context
 * - Request bodies are bound using ShouldBindJSON
 * - Responses use appropriate HTTP status codes:
 *   - 200: Successful GET/PUT
 *   - 201: Successful POST (resource created)
 *   - 400: Bad request (validation error)
 *   - 404: Resource not found
 *   - 500: Internal server error
 *
 * Author: Token Manager Team
 * Version: 1.0.0
 * =============================================================================
 */

package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/**
 * =============================================================================
 * Route Registration
 * Registers all API routes with the Gin engine.
 * =============================================================================
 */

/**
 * RegisterRoutes sets up all API endpoints for the application.
 *
 * Routes are organized under the /api prefix:
 *
 * Employee Endpoints:
 * - POST   /api/employees         - Create new employee
 * - GET    /api/employees        - List all employees
 * - GET    /api/employees/:id    - Get specific employee
 * - PUT    /api/employees/:id    - Update employee
 * - DELETE /api/employees/:id    - Delete employee
 *
 * Token Endpoints:
 * - POST   /api/tokens/issue              - Issue new token
 * - GET    /api/tokens                   - List all tokens
 * - GET    /api/employees/:id/tokens     - Get employee's tokens
 * - POST   /api/tokens/:id/revoke        - Revoke a token
 * - POST   /api/tokens/use               - Use a token
 * - GET    /api/tokens/:id/usage         - Get token usage history
 *
 * Statistics Endpoints:
 * - GET    /api/stats                    - Get system statistics
 * - GET    /api/employees/:id/stats      - Get employee statistics
 *
 * Other Endpoints:
 * - GET    /api/mappings                 - Get token-employee mappings
 *
 * @param r *gin.Engine - The Gin engine to register routes with
 */
func (s *Store) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Employee routes
		api.POST("/employees", s.createEmployee)
		api.GET("/employees", s.getAllEmployees)
		api.GET("/employees/:id", s.getEmployee)
		api.PUT("/employees/:id", s.updateEmployee)
		api.DELETE("/employees/:id", s.deleteEmployee)

		// Token routes
		api.POST("/tokens/issue", s.issueToken)
		api.GET("/tokens", s.getAllTokens)
		api.GET("/employees/:id/tokens", s.getEmployeeTokens)
		api.POST("/tokens/:id/revoke", s.revokeToken)
		api.POST("/tokens/use", s.useToken)
		api.GET("/tokens/:id/usage", s.getTokenUsage)

		// Statistics routes
		api.GET("/stats", s.getStats)
		api.GET("/employees/:id/stats", s.getEmployeeStats)

		// Mapping routes
		api.GET("/mappings", s.getMappings)
	}
}

/**
 * =============================================================================
 * Employee Handlers
 * HTTP handlers for employee-related operations.
 * =============================================================================
 */

/**
 * createEmployee handles POST /api/employees
 *
 * Creates a new employee in the system.
 *
 * Request Body: CreateEmployeeRequest (JSON)
 * - name: Employee's full name (required)
 * - email: Employee's email (required, must be valid email)
 * - department: Department name (optional)
 * - position: Job position (optional)
 *
 * Response: Employee (JSON) with HTTP 201 Created
 *
 * Error Responses:
 * - 400: Invalid request body or validation error
 * - 500: Internal server error during creation
 */
func (s *Store) createEmployee(c *gin.Context) {
	// Parse and validate request body
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create employee through store
	emp, err := s.CreateEmployee(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return created employee with 201 status
	c.JSON(http.StatusCreated, emp)
}

/**
 * getAllEmployees handles GET /api/employees
 *
 * Retrieves a list of all employees in the system.
 *
 * Response: Array of Employee objects (JSON) with HTTP 200 OK
 */
func (s *Store) getAllEmployees(c *gin.Context) {
	employees := s.GetAllEmployees()
	c.JSON(http.StatusOK, employees)
}

/**
 * getEmployee handles GET /api/employees/:id
 *
 * Retrieves a specific employee by their ID.
 *
 * URL Parameters:
 * - id: Employee ID
 *
 * Response: Employee object (JSON) with HTTP 200 OK
 *
 * Error Responses:
 * - 404: Employee not found
 */
func (s *Store) getEmployee(c *gin.Context) {
	// Extract employee ID from URL parameter
	id := c.Param("id")

	// Fetch employee from store
	emp := s.GetEmployee(id)
	if emp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}

	c.JSON(http.StatusOK, emp)
}

/**
 * updateEmployee handles PUT /api/employees/:id
 *
 * Updates an existing employee's information.
 *
 * URL Parameters:
 * - id: Employee ID
 *
 * Request Body: CreateEmployeeRequest (JSON)
 *
 * Response: Updated Employee object (JSON) with HTTP 200 OK
 *
 * Error Responses:
 * - 400: Invalid request body or validation error
 * - 404: Employee not found
 */
func (s *Store) updateEmployee(c *gin.Context) {
	// Extract employee ID from URL parameter
	id := c.Param("id")

	// Parse request body
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update employee through store
	emp, err := s.UpdateEmployee(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, emp)
}

/**
 * deleteEmployee handles DELETE /api/employees/:id
 *
 * Removes an employee from the system.
 *
 * URL Parameters:
 * - id: Employee ID
 *
 * Response: Success message (JSON) with HTTP 200 OK
 *
 * Error Responses:
 * - 404: Employee not found
 */
func (s *Store) deleteEmployee(c *gin.Context) {
	// Extract employee ID from URL parameter
	id := c.Param("id")

	// Delete employee through store
	if err := s.DeleteEmployee(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee deleted"})
}

/**
 * =============================================================================
 * Token Handlers
 * HTTP handlers for token-related operations.
 * =============================================================================
 */

/**
 * issueToken handles POST /api/tokens/issue
 *
 * Issues a new token to an employee.
 *
 * Request Body: IssueTokenRequest (JSON)
 * - employee_id: ID of the employee to receive the token (required)
 * - total_quota: Total token quota amount (required, > 0)
 * - days_valid: Number of days the token is valid (required, > 0)
 *
 * Response: Token object (JSON) with HTTP 201 Created
 *
 * Error Responses:
 * - 400: Invalid request body or validation error
 * - 500: Internal server error during issuance
 */
func (s *Store) issueToken(c *gin.Context) {
	// Parse and validate request body
	var req IssueTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Issue token through store
	token, err := s.IssueToken(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, token)
}

/**
 * getAllTokens handles GET /api/tokens
 *
 * Retrieves a list of all tokens in the system.
 *
 * Response: Array of Token objects (JSON) with HTTP 200 OK
 */
func (s *Store) getAllTokens(c *gin.Context) {
	tokens := s.GetAllTokens()
	c.JSON(http.StatusOK, tokens)
}

/**
 * getEmployeeTokens handles GET /api/employees/:id/tokens
 *
 * Retrieves all tokens assigned to a specific employee.
 *
 * URL Parameters:
 * - id: Employee ID
 *
 * Response: Array of Token objects (JSON) with HTTP 200 OK
 */
func (s *Store) getEmployeeTokens(c *gin.Context) {
	// Extract employee ID from URL parameter
	id := c.Param("id")

	// Fetch employee's tokens from store
	tokens := s.GetEmployeeTokens(id)
	c.JSON(http.StatusOK, tokens)
}

/**
 * revokeToken handles POST /api/tokens/:id/revoke
 *
 * Revokes a token, making it inactive. Revoked tokens cannot be used
 * but their usage history is preserved.
 *
 * URL Parameters:
 * - id: Token ID
 *
 * Response: Success message (JSON) with HTTP 200 OK
 *
 * Error Responses:
 * - 404: Token not found
 */
func (s *Store) revokeToken(c *gin.Context) {
	// Extract token ID from URL parameter
	id := c.Param("id")

	// Revoke token through store
	if err := s.RevokeToken(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "token revoked"})
}

/**
 * useToken handles POST /api/tokens/use
 *
 * Records a token usage event. This decrements the token's remaining quota.
 *
 * Request Body: UseTokenRequest (JSON)
 * - token_id: ID of the token being used (required)
 * - amount: Number of tokens consumed (required, > 0)
 * - model: LLM model used (optional)
 * - description: Description of the usage (optional)
 *
 * Response: TokenUsage object (JSON) with HTTP 201 Created
 *
 * Error Responses:
 * - 400: Invalid request, token not active, expired, or insufficient quota
 */
func (s *Store) useToken(c *gin.Context) {
	// Parse and validate request body
	var req UseTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Record token usage through store
	usage, err := s.UseToken(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, usage)
}

/**
 * getTokenUsage handles GET /api/tokens/:id/usage
 *
 * Retrieves the usage history for a specific token.
 *
 * URL Parameters:
 * - id: Token ID
 *
 * Response: Array of TokenUsage objects (JSON) with HTTP 200 OK
 */
func (s *Store) getTokenUsage(c *gin.Context) {
	// Extract token ID from URL parameter
	id := c.Param("id")

	// Fetch token's usage history from store
	usages := s.GetTokenUsage(id)
	c.JSON(http.StatusOK, usages)
}

/**
 * =============================================================================
 * Statistics Handlers
 * HTTP handlers for statistics and analytics endpoints.
 * =============================================================================
 */

/**
 * getStats handles GET /api/stats
 *
 * Retrieves aggregated system statistics including:
 * - Total employees, tokens, active tokens
 * - Token quotas and usage percentages
 *
 * Response: StatsResponse object (JSON) with HTTP 200 OK
 */
func (s *Store) getStats(c *gin.Context) {
	// Fetch system statistics from store
	stats := s.GetStats()
	c.JSON(http.StatusOK, stats)
}

/**
 * getEmployeeStats handles GET /api/employees/:id/stats
 *
 * Retrieves statistics specific to an employee including:
 * - Employee information
 * - Token counts and quotas
 * - Usage statistics and percentages
 *
 * URL Parameters:
 * - id: Employee ID
 *
 * Response: EmployeeStatsResponse object (JSON) with HTTP 200 OK
 *
 * Error Responses:
 * - 404: Employee not found
 */
func (s *Store) getEmployeeStats(c *gin.Context) {
	// Extract employee ID from URL parameter
	id := c.Param("id")

	// Fetch employee statistics from store
	stats, err := s.GetEmployeeStats(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

/**
 * =============================================================================
 * Mapping Handlers
 * HTTP handlers for token-employee relationship endpoints.
 * =============================================================================
 */

/**
 * getMappings handles GET /api/mappings
 *
 * Retrieves all token-employee relationships with pagination support.
 * This endpoint shows which employees own which tokens.
 *
 * Query Parameters:
 * - page: Page number (default: 1)
 * - page_size: Number of items per page (default: 10)
 *
 * Response: Object with pagination metadata and array of TokenMapping (JSON)
 *
 * Response Structure:
 * {
 *     "total": 100,
 *     "page": 1,
 *     "page_size": 10,
 *     "data": [TokenMapping...]
 * }
 */
func (s *Store) getMappings(c *gin.Context) {
	// Extract pagination parameters from query string
	// Default to page 1 with 10 items per page if not specified
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	// Convert string parameters to integers
	pageNum, _ := strconv.Atoi(page)
	size, _ := strconv.Atoi(pageSize)

	// Get all mappings from store
	mappings := s.GetMappings()

	// Calculate pagination bounds
	start := (pageNum - 1) * size
	end := start + size

	// Ensure bounds don't exceed array length
	if start > len(mappings) {
		start = len(mappings)
	}
	if end > len(mappings) {
		end = len(mappings)
	}

	// Extract page of data
	var pagedMappings []TokenMapping
	if start < end {
		pagedMappings = mappings[start:end]
	} else {
		pagedMappings = []TokenMapping{}
	}

	// Return paginated response with metadata
	c.JSON(http.StatusOK, gin.H{
		"total":     len(mappings),
		"page":      pageNum,
		"page_size": size,
		"data":      pagedMappings,
	})
}
