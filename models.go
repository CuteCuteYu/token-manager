/**
 * =============================================================================
 * Token Manager - Data Models
 *
 * This file defines all data structures (structs) used throughout the application.
 * These models represent the core entities: employees, tokens, and usage records.
 *
 * Table of Contents:
 * 1. Package Declaration
 * 2. Import Statements
 * 3. Core Entity Structs
 * 4. Request/Response DTOs
 *
 * Design Notes:
 * - All structs use JSON tags for serialization/deserialization
 * - snake_case is used for JSON field names (Go convention)
 * - Pointer types (*time.Time) are used when a field may be nil
 *
 * Author: Token Manager Team
 * Version: 1.0.0
 * =============================================================================
 */

package main

import "time"

/**
 * =============================================================================
 * Core Entity Structs
 * These structs represent the main data entities in the system.
 * =============================================================================
 */

/**
 * Employee represents a company employee in the system.
 *
 * Fields:
 * - ID: Unique identifier (generated from timestamp)
 * - Name: Employee's full name
 * - Department: Employee's department (optional)
 * - Email: Employee's email address (unique)
 * - Position: Job title (optional)
 * - CreatedAt: Timestamp when employee was created
 * - UpdatedAt: Timestamp when employee was last updated
 *
 * JSON Example:
 * {
 *     "id": "1234567890",
 *     "name": "John Doe",
 *     "department": "Engineering",
 *     "email": "john.doe@company.com",
 *     "position": "Software Engineer",
 *     "created_at": "2024-01-15T10:30:00Z",
 *     "updated_at": "2024-01-15T10:30:00Z"
 * }
 */
type Employee struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Department string    `json:"department"`
	Email      string    `json:"email"`
	Position   string    `json:"position"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

/**
 * Token represents an LLM token allocated to an employee.
 *
 * Fields:
 * - ID: Unique identifier for the token
 * - EmployeeID: Reference to the employee who owns this token
 * - TokenValue: The actual token string used for API authentication
 * - TotalQuota: Maximum number of tokens this token allows
 * - UsedQuota: Number of tokens already consumed
 * - IsActive: Whether the token can still be used
 * - IssuedAt: When the token was created
 * - ExpiredAt: When the token expires (no longer usable after this time)
 * - RevokedAt: When the token was revoked (nil if still active)
 *
 * JSON Example:
 * {
 *     "id": "token_123",
 *     "employee_id": "emp_456",
 *     "token_value": "tok_1234567890_abc12345",
 *     "total_quota": 1000000,
 *     "used_quota": 250000,
 *     "is_active": true,
 *     "issued_at": "2024-01-15T10:30:00Z",
 *     "expired_at": "2024-02-15T10:30:00Z",
 *     "revoked_at": null
 * }
 */
type Token struct {
	ID         string     `json:"id"`
	EmployeeID string     `json:"employee_id"`
	TokenValue string     `json:"token_value"`
	TotalQuota int        `json:"total_quota"`
	UsedQuota  int        `json:"used_quota"`
	IsActive   bool       `json:"is_active"`
	IssuedAt   time.Time  `json:"issued_at"`
	ExpiredAt  time.Time  `json:"expired_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
}

/**
 * TokenUsage represents a single usage event of a token.
 *
 * Fields:
 * - ID: Unique identifier for this usage record
 * - TokenID: Reference to the token that was used
 * - UsedAt: Timestamp when the token was used
 * - Amount: Number of tokens consumed in this usage
 * - Model: The LLM model that was used (e.g., "gpt-4")
 * - Description: Optional description of what this usage was for
 *
 * JSON Example:
 * {
 *     "id": "usage_789",
 *     "token_id": "token_123",
 *     "used_at": "2024-01-15T14:30:00Z",
 *     "amount": 1500,
 *     "model": "gpt-4",
 *     "description": "Code review assistance"
 * }
 */
type TokenUsage struct {
	ID          string    `json:"id"`
	TokenID     string    `json:"token_id"`
	UsedAt      time.Time `json:"used_at"`
	Amount      int       `json:"amount"`
	Model       string    `json:"model"`
	Description string    `json:"description"`
}

/**
 * DataStore is the root data structure that holds all application data.
 * This is what gets serialized to and from the JSON file.
 *
 * Fields:
 * - Employees: Slice of all employees
 * - Tokens: Slice of all tokens
 * - UsageRecords: Slice of all usage records
 *
 * This structure is persisted to data.json for data persistence.
 */
type DataStore struct {
	Employees    []Employee   `json:"employees"`
	Tokens       []Token      `json:"tokens"`
	UsageRecords []TokenUsage `json:"usage_records"`
}

/**
 * =============================================================================
 * Request DTOs (Data Transfer Objects)
 * These structs define the expected structure of incoming API requests.
 * They include validation tags for automatic request validation.
 * =============================================================================
 */

/**
 * CreateEmployeeRequest defines the expected structure for creating a new employee.
 *
 * Validation Rules (from binding tags):
 * - Name: Required
 * - Email: Required and must be a valid email format
 * - Department: Optional
 * - Position: Optional
 */
type CreateEmployeeRequest struct {
	Name       string `json:"name" binding:"required"`
	Department string `json:"department"`
	Email      string `json:"email" binding:"required,email"`
	Position   string `json:"position"`
}

/**
 * IssueTokenRequest defines the expected structure for issuing a new token.
 *
 * Validation Rules:
 * - EmployeeID: Required - must reference an existing employee
 * - TotalQuota: Required and must be greater than 0
 * - DaysValid: Required and must be greater than 0
 */
type IssueTokenRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	TotalQuota int    `json:"total_quota" binding:"required,gt=0"`
	DaysValid  int    `json:"days_valid" binding:"required,gt=0"`
}

/**
 * UseTokenRequest defines the expected structure for using a token.
 *
 * Validation Rules:
 * - TokenID: Required - must reference an existing token
 * - Amount: Required and must be greater than 0
 * - Model: Optional - the LLM model used
 * - Description: Optional - description of the usage
 */
type UseTokenRequest struct {
	TokenID     string `json:"token_id" binding:"required"`
	Amount      int    `json:"amount" binding:"required,gt=0"`
	Model       string `json:"model"`
	Description string `json:"description"`
}

/**
 * =============================================================================
 * Response DTOs
 * These structs define the structure of API responses.
 * =============================================================================
 */

/**
 * StatsResponse contains aggregated system statistics.
 *
 * Fields:
 * - TotalEmployees: Total number of employees in the system
 * - TotalTokens: Total number of tokens (active + revoked)
 * - ActiveTokens: Number of currently active tokens
 * - RevokedTokens: Number of revoked tokens
 * - TotalQuota: Sum of all token quotas
 * - UsedQuota: Sum of all used quotas
 * - RemainingQuota: Total quota minus used quota
 * - UsagePercent: Percentage of quota that has been used
 */
type StatsResponse struct {
	TotalEmployees int     `json:"total_employees"`
	TotalTokens    int     `json:"total_tokens"`
	ActiveTokens   int     `json:"active_tokens"`
	RevokedTokens  int     `json:"revoked_tokens"`
	TotalQuota     int     `json:"total_quota"`
	UsedQuota      int     `json:"used_quota"`
	RemainingQuota int     `json:"remaining_quota"`
	UsagePercent   float64 `json:"usage_percent"`
}

/**
 * EmployeeStatsResponse contains statistics specific to an employee.
 *
 * Fields:
 * - EmployeeInfo: Pointer to the employee's data
 * - TokenCount: Total tokens assigned to this employee
 * - ActiveTokens: Number of active tokens for this employee
 * - TotalQuota: Sum of all quotas across employee's tokens
 * - UsedQuota: Sum of all used quotas
 * - RemainingQuota: Remaining quota available
 * - UsagePercent: Percentage of quota used
 */
type EmployeeStatsResponse struct {
	EmployeeInfo   *Employee `json:"employee_info"`
	TokenCount     int       `json:"token_count"`
	ActiveTokens   int       `json:"active_tokens"`
	TotalQuota     int       `json:"total_quota"`
	UsedQuota      int       `json:"used_quota"`
	RemainingQuota int       `json:"remaining_quota"`
	UsagePercent   float64   `json:"usage_percent"`
}

/**
 * TokenMapping represents the relationship between a token and employee.
 * This is used for the /api/mappings endpoint to show who owns which tokens.
 *
 * Fields:
 * - EmployeeID: ID of the employee
 * - EmployeeName: Name of the employee
 * - TokenID: ID of the token
 * - TokenValue: The token value string
 * - TotalQuota: Total quota for this token
 * - UsedQuota: Used quota for this token
 * - Remaining: Remaining quota
 * - IsActive: Whether the token is active
 */
type TokenMapping struct {
	EmployeeID   string `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
	TokenID      string `json:"token_id"`
	TokenValue   string `json:"token_value"`
	TotalQuota   int    `json:"total_quota"`
	UsedQuota    int    `json:"used_quota"`
	Remaining    int    `json:"remaining"`
	IsActive     bool   `json:"is_active"`
}
