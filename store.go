/**
 * =============================================================================
 * Token Manager - Data Store
 *
 * This file contains the Store type which manages all data operations.
 * It provides an in-memory data store with JSON file persistence.
 *
 * Table of Contents:
 * 1. Package Declaration
 * 2. Import Statements
 * 3. Store Type Definition
 * 4. Store Initialization
 * 5. Data Persistence Methods
 * 6. ID Generation Methods
 * 7. Employee Business Logic
 * 8. Token Business Logic
 * 9. Statistics Methods
 * 10. Mapping Methods
 *
 * Thread Safety:
 * - All data access is protected by sync.RWMutex
 * - Read operations use RLock (concurrent reads allowed)
 * - Write operations use Lock (exclusive access required)
 *
 * Design Patterns:
 * - Repository Pattern: Abstracts data access logic
 * - Factory Pattern: NewStore function acts as factory
 * - Unit of Work: saveToFileLocked groups related changes
 *
 * Author: Token Manager Team
 * Version: 1.0.0
 * =============================================================================
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

/**
 * =============================================================================
 * Store Type Definition
 * The Store manages all application data with thread-safe access.
 * =============================================================================
 */

/**
 * Store represents the data store for the application.
 * It holds all employees, tokens, and usage records in memory
 * and persists them to a JSON file.
 *
 * Fields:
 * - filePath: Path to the JSON file for persistence
 * - mu: Read-write mutex for thread-safe data access
 * - data: In-memory DataStore containing all application data
 *
 * Thread Safety:
 * - The mutex must be held during all read and write operations
 * - Always use defer to ensure mutex is unlocked
 * - Use RLock for reads, Lock for writes
 */
type Store struct {
	filePath string
	mu       sync.RWMutex
	data     DataStore
}

/**
 * =============================================================================
 * Store Initialization
 * Functions for creating and initializing a new Store.
 * =============================================================================
 */

/**
 * NewStore creates a new Store instance and loads existing data.
 *
 * This function:
 * 1. Creates a new Store with empty data
 * 2. Attempts to load existing data from the file
 * 3. If file doesn't exist, creates a new one with empty data
 *
 * @param filePath string - Path to the JSON data file
 * @return *Store - New Store instance
 * @return error - Any error that occurred during initialization
 *
 * Note:
 * - If the file doesn't exist, it will be created automatically
 * - If loading fails for other reasons, the error is returned
 */
func NewStore(filePath string) (*Store, error) {
	// Create new store with empty data and specified file path
	store := &Store{
		filePath: filePath,
		data:     DataStore{},
	}

	// Attempt to load existing data from file
	if err := store.load(); err != nil {
		// If file doesn't exist, create a new one
		if os.IsNotExist(err) {
			return store, store.save()
		}
		// Return other errors (permission issues, corrupted file, etc.)
		return nil, err
	}

	return store, nil
}

/**
 * =============================================================================
 * Data Persistence Methods
 * Methods for loading and saving data to JSON file.
 * =============================================================================
 */

/**
 * load reads the JSON data file into memory.
 *
 * This method acquires a write lock since it modifies the internal data.
 *
 * Error Handling:
 * - Returns error if file cannot be read
 * - Returns error if JSON cannot be parsed
 */
func (s *Store) load() error {
	// Acquire write lock for thread safety
	s.mu.Lock()
	defer s.mu.Unlock()

	// Read file contents
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	// Parse JSON into internal data structure
	return json.Unmarshal(data, &s.data)
}

/**
 * save writes the current in-memory data to the JSON file.
 *
 * This method acquires a write lock since it modifies the file.
 * Data is formatted with indentation for readability.
 *
 * File permissions: 0644 (rw-r--r--)
 * - Owner: read and write
 * - Group: read
 * - Others: read
 */
func (s *Store) save() error {
	// Acquire write lock for thread safety
	s.mu.Lock()
	defer s.mu.Unlock()

	// Serialize data to JSON with indentation
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(s.filePath, data, 0644)
}

/**
 * saveToFileLocked saves data to file when already holding write lock.
 *
 * This is an optimization to avoid acquiring the lock twice when
 * the caller already holds the write lock.
 *
 * Precondition: Caller must hold s.mu.Lock()
 * Postcondition: File is written, lock is still held
 */
func (s *Store) saveToFileLocked() error {
	// Serialize data to JSON with indentation
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	// Write to file (lock already held by caller)
	return os.WriteFile(s.filePath, data, 0644)
}

/**
 * =============================================================================
 * ID Generation Methods
 * Helper methods for generating unique identifiers.
 * =============================================================================
 */

/**
 * generateID creates a unique ID based on current timestamp.
 *
 * Uses UnixNano for high-resolution timestamp to ensure uniqueness.
 * The ID is returned as a string representation of the timestamp.
 *
 * @return string - Unique identifier
 */
func (s *Store) generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

/**
 * generateTokenValue creates a unique token value string.
 *
 * Token format: "tok_<timestamp>_<id_prefix>"
 * Example: "tok_1705312200000000000_abc12345"
 *
 * This format:
 * - Starts with "tok_" prefix for identification
 * - Includes timestamp for uniqueness
 * - Includes first 8 characters of generated ID for additional entropy
 *
 * @return string - Unique token value
 */
func (s *Store) generateTokenValue() string {
	return fmt.Sprintf("tok_%d_%s", time.Now().UnixNano(), s.generateID()[:8])
}

/**
 * =============================================================================
 * Employee Business Logic
 * CRUD operations for employees.
 * =============================================================================
 */

/**
 * CreateEmployee adds a new employee to the system.
 *
 * This method:
 * 1. Acquires write lock for thread safety
 * 2. Creates new Employee with generated ID and timestamps
 * 3. Appends to employees slice
 * 4. Persists changes to file
 *
 * @param req CreateEmployeeRequest - Employee creation data
 * @return *Employee - Newly created employee
 * @return error - Any error that occurred
 */
func (s *Store) CreateEmployee(req CreateEmployeeRequest) (*Employee, error) {
	// Acquire write lock
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create employee with generated ID and current timestamps
	emp := Employee{
		ID:         s.generateID(),
		Name:       req.Name,
		Department: req.Department,
		Email:      req.Email,
		Position:   req.Position,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Add employee to data store
	s.data.Employees = append(s.data.Employees, emp)

	// Persist changes to file
	if err := s.saveToFileLocked(); err != nil {
		return nil, err
	}

	return &emp, nil
}

/**
 * GetAllEmployees returns all employees in the system.
 *
 * Uses read lock to allow concurrent reads.
 *
 * @return []Employee - Copy of all employees
 */
func (s *Store) GetAllEmployees() []Employee {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.Employees
}

/**
 * GetEmployee retrieves a single employee by ID.
 *
 * Uses read lock for thread-safe access.
 * Returns nil if employee not found.
 *
 * @param id string - Employee ID to search for
 * @return *Employee - Pointer to employee if found, nil otherwise
 */
func (s *Store) GetEmployee(id string) *Employee {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Linear search through employees
	for _, emp := range s.data.Employees {
		if emp.ID == id {
			return &emp
		}
	}
	return nil
}

/**
 * UpdateEmployee modifies an existing employee's information.
 *
 * Uses write lock for thread safety.
 * Updates the UpdatedAt timestamp to track modifications.
 *
 * @param id string - Employee ID to update
 * @param req CreateEmployeeRequest - New employee data
 * @return *Employee - Updated employee
 * @return error - Error if employee not found
 */
func (s *Store) UpdateEmployee(id string, req CreateEmployeeRequest) (*Employee, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find and update employee
	for i, emp := range s.data.Employees {
		if emp.ID == id {
			// Update employee fields
			s.data.Employees[i].Name = req.Name
			s.data.Employees[i].Department = req.Department
			s.data.Employees[i].Email = req.Email
			s.data.Employees[i].Position = req.Position
			s.data.Employees[i].UpdatedAt = time.Now()

			// Persist changes
			if err := s.saveToFileLocked(); err != nil {
				return nil, err
			}

			return &s.data.Employees[i], nil
		}
	}
	return nil, fmt.Errorf("employee not found")
}

/**
 * DeleteEmployee removes an employee from the system.
 *
 * Uses write lock for thread safety.
 * Note: This does NOT delete associated tokens; they remain
 * in the system but are orphaned.
 *
 * @param id string - Employee ID to delete
 * @return error - Error if employee not found
 */
func (s *Store) DeleteEmployee(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find and remove employee using slice manipulation
	for i, emp := range s.data.Employees {
		if emp.ID == id {
			// Remove element by concatenating slices
			s.data.Employees = append(s.data.Employees[:i], s.data.Employees[i+1:]...)
			return s.saveToFileLocked()
		}
	}
	return fmt.Errorf("employee not found")
}

/**
 * =============================================================================
 * Token Business Logic
 * CRUD operations for tokens.
 * =============================================================================
 */

/**
 * IssueToken creates a new token for an employee.
 *
 * This method:
 * 1. Validates that the employee exists
 * 2. Creates token with generated ID and value
 * 3. Sets expiration based on days_valid parameter
 * 4. Persists to file
 *
 * @param req IssueTokenRequest - Token issuance data
 * @return *Token - Newly created token
 * @return error - Error if employee not found or other issue
 */
func (s *Store) IssueToken(req IssueTokenRequest) (*Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify employee exists
	empFound := false
	for _, emp := range s.data.Employees {
		if emp.ID == req.EmployeeID {
			empFound = true
			break
		}
	}
	if !empFound {
		return nil, fmt.Errorf("employee not found")
	}

	// Create new token
	token := Token{
		ID:         s.generateID(),
		EmployeeID: req.EmployeeID,
		TokenValue: s.generateTokenValue(),
		TotalQuota: req.TotalQuota,
		UsedQuota:  0,
		IsActive:   true,
		IssuedAt:   time.Now(),
		ExpiredAt:  time.Now().AddDate(0, 0, req.DaysValid),
	}

	// Add token to data store
	s.data.Tokens = append(s.data.Tokens, token)

	// Persist changes
	if err := s.saveToFileLocked(); err != nil {
		return nil, err
	}

	return &token, nil
}

/**
 * GetAllTokens returns all tokens in the system.
 *
 * Uses read lock for thread-safe concurrent access.
 *
 * @return []Token - Copy of all tokens
 */
func (s *Store) GetAllTokens() []Token {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.Tokens
}

/**
 * GetEmployeeTokens returns all tokens assigned to an employee.
 *
 * Uses read lock for thread-safe access.
 *
 * @param employeeID string - Employee ID to filter tokens
 * @return []Token - Slice of employee's tokens
 */
func (s *Store) GetEmployeeTokens(employeeID string) []Token {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Collect matching tokens
	var tokens []Token
	for _, t := range s.data.Tokens {
		if t.EmployeeID == employeeID {
			tokens = append(tokens, t)
		}
	}
	return tokens
}

/**
 * GetToken retrieves a single token by ID.
 *
 * Uses read lock for thread-safe access.
 *
 * @param id string - Token ID to search for
 * @return *Token - Pointer to token if found, nil otherwise
 */
func (s *Store) GetToken(id string) *Token {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, t := range s.data.Tokens {
		if t.ID == id {
			return &t
		}
	}
	return nil
}

/**
 * RevokeToken marks a token as inactive.
 *
 * This sets IsActive to false and records the revocation timestamp.
 * The token remains in the system for historical purposes.
 *
 * @param id string - Token ID to revoke
 * @return error - Error if token not found
 */
func (s *Store) RevokeToken(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	// Find and revoke token
	for i, t := range s.data.Tokens {
		if t.ID == id {
			s.data.Tokens[i].IsActive = false
			s.data.Tokens[i].RevokedAt = &now
			return s.saveToFileLocked()
		}
	}
	return fmt.Errorf("token not found")
}

/**
 * UseToken records a token usage event.
 *
 * This method:
 * 1. Validates token exists and is active
 * 2. Checks token hasn't expired
 * 3. Verifies sufficient quota remaining
 * 4. Updates used quota
 * 5. Creates usage record
 * 6. Persists changes
 *
 * @param req UseTokenRequest - Usage details
 * @return *TokenUsage - Created usage record
 * @return error - Error for validation failures
 */
func (s *Store) UseToken(req UseTokenRequest) (*TokenUsage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find token
	var token *Token
	var tokenIdx int
	for i, t := range s.data.Tokens {
		if t.ID == req.TokenID {
			token = &t
			tokenIdx = i
			break
		}
	}

	// Validate token exists
	if token == nil {
		return nil, fmt.Errorf("token not found")
	}

	// Validate token is active
	if !token.IsActive {
		return nil, fmt.Errorf("token is not active")
	}

	// Validate token hasn't expired
	if time.Now().After(token.ExpiredAt) {
		return nil, fmt.Errorf("token has expired")
	}

	// Validate sufficient quota
	remaining := token.TotalQuota - token.UsedQuota
	if req.Amount > remaining {
		return nil, fmt.Errorf("insufficient token quota, remaining: %d", remaining)
	}

	// Update used quota
	s.data.Tokens[tokenIdx].UsedQuota += req.Amount

	// Create usage record
	usage := TokenUsage{
		ID:          s.generateID(),
		TokenID:     req.TokenID,
		UsedAt:      time.Now(),
		Amount:      req.Amount,
		Model:       req.Model,
		Description: req.Description,
	}

	// Add to usage records
	s.data.UsageRecords = append(s.data.UsageRecords, usage)

	// Persist changes
	if err := s.saveToFileLocked(); err != nil {
		return nil, err
	}

	return &usage, nil
}

/**
 * GetTokenUsage returns usage history for a token.
 *
 * Uses read lock for thread-safe access.
 *
 * @param tokenID string - Token ID to get usage for
 * @return []TokenUsage - Slice of usage records
 */
func (s *Store) GetTokenUsage(tokenID string) []TokenUsage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Collect matching usage records
	var usages []TokenUsage
	for _, u := range s.data.UsageRecords {
		if u.TokenID == tokenID {
			usages = append(usages, u)
		}
	}
	return usages
}

/**
 * =============================================================================
 * Statistics Methods
 * Methods for calculating system and employee statistics.
 * =============================================================================
 */

/**
 * GetStats calculates aggregated system statistics.
 *
 * Uses read lock for thread-safe access.
 * Calculates totals across all tokens and employees.
 *
 * @return StatsResponse - Aggregated statistics
 */
func (s *Store) GetStats() StatsResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Initialize counters
	var totalTokens, activeTokens, revokedTokens, totalQuota, usedQuota int

	// Aggregate token statistics
	for _, t := range s.data.Tokens {
		totalTokens++
		totalQuota += t.TotalQuota
		usedQuota += t.UsedQuota
		if t.IsActive {
			activeTokens++
		} else {
			revokedTokens++
		}
	}

	// Calculate usage percentage
	usagePercent := 0.0
	if totalQuota > 0 {
		usagePercent = float64(usedQuota) / float64(totalQuota) * 100
	}

	// Return aggregated statistics
	return StatsResponse{
		TotalEmployees: len(s.data.Employees),
		TotalTokens:    totalTokens,
		ActiveTokens:   activeTokens,
		RevokedTokens:  revokedTokens,
		TotalQuota:     totalQuota,
		UsedQuota:      usedQuota,
		RemainingQuota: totalQuota - usedQuota,
		UsagePercent:   usagePercent,
	}
}

/**
 * GetEmployeeStats calculates statistics for a specific employee.
 *
 * Uses read lock for thread-safe access.
 * Includes employee info and token usage statistics.
 *
 * @param employeeID string - Employee ID to get stats for
 * @return *EmployeeStatsResponse - Employee-specific statistics
 * @return error - Error if employee not found
 */
func (s *Store) GetEmployeeStats(employeeID string) (*EmployeeStatsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Find employee
	var emp *Employee
	for _, e := range s.data.Employees {
		if e.ID == employeeID {
			emp = &e
			break
		}
	}

	// Validate employee exists
	if emp == nil {
		return nil, fmt.Errorf("employee not found")
	}

	// Initialize counters
	var tokenCount, activeTokens, totalQuota, usedQuota int

	// Aggregate token statistics for this employee
	for _, t := range s.data.Tokens {
		if t.EmployeeID == employeeID {
			tokenCount++
			totalQuota += t.TotalQuota
			usedQuota += t.UsedQuota
			if t.IsActive {
				activeTokens++
			}
		}
	}

	// Calculate usage percentage
	usagePercent := 0.0
	if totalQuota > 0 {
		usagePercent = float64(usedQuota) / float64(totalQuota) * 100
	}

	// Return employee statistics
	return &EmployeeStatsResponse{
		EmployeeInfo:   emp,
		TokenCount:     tokenCount,
		ActiveTokens:   activeTokens,
		TotalQuota:     totalQuota,
		UsedQuota:      usedQuota,
		RemainingQuota: totalQuota - usedQuota,
		UsagePercent:   usagePercent,
	}, nil
}

/**
 * =============================================================================
 * Mapping Methods
 * Methods for retrieving token-employee relationships.
 * =============================================================================
 */

/**
 * GetMappings returns all token-employee relationships.
 *
 * Uses read lock for thread-safe access.
 * Builds a mapping by combining employee and token data.
 * Creates an in-memory map for quick employee name lookup.
 *
 * @return []TokenMapping - All token-employee relationships
 */
func (s *Store) GetMappings() []TokenMapping {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Build employee ID to name mapping
	empMap := make(map[string]string)
	for _, e := range s.data.Employees {
		empMap[e.ID] = e.Name
	}

	// Build token-employee mappings
	var mappings []TokenMapping
	for _, t := range s.data.Tokens {
		mapping := TokenMapping{
			EmployeeID:   t.EmployeeID,
			EmployeeName: empMap[t.EmployeeID],
			TokenID:      t.ID,
			TokenValue:   t.TokenValue,
			TotalQuota:   t.TotalQuota,
			UsedQuota:    t.UsedQuota,
			Remaining:    t.TotalQuota - t.UsedQuota,
			IsActive:     t.IsActive,
		}
		mappings = append(mappings, mapping)
	}

	return mappings
}
