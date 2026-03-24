/**
 * =============================================================================
 * Token Manager - Application JavaScript
 * 
 * This file contains all client-side JavaScript functionality for the
 * Token Manager frontend application. It handles API communication,
 * DOM manipulation, and user interactions.
 * 
 * Table of Contents:
 * 1. Configuration Constants
 * 2. API Utility Functions
 * 3. Data Formatting Functions
 * 4. Data Loading Functions
 * 5. User Action Functions
 * 6. Form Handlers
 * 7. Navigation Handlers
 * 8. Initialization
 * 
 * Author: Token Manager Team
 * Version: 1.0.0
 * =============================================================================
 */


/**
 * =============================================================================
 * 1. Configuration Constants
 * Global configuration and constants used throughout the application.
 * =============================================================================
 */

// API_BASE: The base URL for the backend REST API
// This should point to the server where the Go application is running
const API_BASE = 'http://localhost:8080/api';


/**
 * =============================================================================
 * 2. API Utility Functions
 * Helper functions for making HTTP requests to the backend API.
 * =============================================================================
 */

/**
 * fetchJSON - Fetch JSON data from the API
 * 
 * @param {string} url - The URL to fetch data from
 * @returns {Promise<Object>} - The parsed JSON response
 * @throws {Error} - Throws an error if the network request fails
 * 
 * This function performs a GET request and automatically parses
 * the JSON response. It includes error handling for network failures.
 */
async function fetchJSON(url) {
    const response = await fetch(url);
    if (!response.ok) throw new Error('Network error');
    return response.json();
}


/**
 * postJSON - Send JSON data to the API using POST method
 * 
 * @param {string} url - The URL to send data to
 * @param {Object} data - The data object to send in the request body
 * @returns {Promise<Object>} - The parsed JSON response from the server
 * @throws {Error} - Throws an error if the request fails, including server-side errors
 * 
 * This function sends a POST request with JSON-encoded data in the body.
 * It handles both network errors and server-returned errors.
 */
async function postJSON(url, data) {
    const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
    if (!response.ok) {
        const err = await response.json();
        throw new Error(err.error || 'Request failed');
    }
    return response.json();
}


/**
 * deleteResource - Delete a resource using DELETE method
 * 
 * @param {string} url - The URL of the resource to delete
 * @returns {Promise<void>} - Resolves when the deletion is successful
 * @throws {Error} - Throws an error if the deletion fails
 * 
 * This function performs a DELETE request to remove resources
 * from the server.
 */
async function deleteResource(url) {
    const response = await fetch(url, { method: 'DELETE' });
    if (!response.ok) throw new Error('Delete failed');
}


/**
 * =============================================================================
 * 3. Data Formatting Functions
 * Utility functions for formatting data for display.
 * =============================================================================
 */

/**
 * formatDate - Format a date string for display
 * 
 * @param {string|null} dateStr - ISO date string to format
 * @returns {string} - Formatted date string or '-' if input is null/empty
 * 
 * Converts ISO date strings into a human-readable format:
 * Example: "2024-01-15T10:30:00Z" -> "Jan 15, 2024, 10:30 AM"
 */
function formatDate(dateStr) {
    if (!dateStr) return '-';
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', { 
        year: 'numeric', 
        month: 'short', 
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}


/**
 * truncate - Truncate a string to a maximum length
 * 
 * @param {string|null} str - The string to truncate
 * @param {number} len - Maximum length before truncation (default: 20)
 * @returns {string} - Truncated string with '...' suffix, or '-' if input is empty
 * 
 * This function limits the display length of long strings like IDs
 * or token values while maintaining readability.
 */
function truncate(str, len = 20) {
    if (!str) return '-';
    return str.length > len ? str.substring(0, len) + '...' : str;
}


/**
 * =============================================================================
 * 4. Data Loading Functions
 * Functions that fetch and render data from the API.
 * =============================================================================
 */

/**
 * loadStats - Load and display system statistics
 * 
 * Fetches statistics from the /api/stats endpoint and updates
 * the overview dashboard with current system metrics including:
 * - Total number of employees
 * - Total number of tokens
 * - Number of active tokens
 * - Overall usage percentage
 * 
 * @returns {Promise<void>}
 */
async function loadStats() {
    try {
        const stats = await fetchJSON(`${API_BASE}/stats`);
        document.getElementById('totalEmployees').textContent = stats.total_employees;
        document.getElementById('totalTokens').textContent = stats.total_tokens;
        document.getElementById('activeTokens').textContent = stats.active_tokens;
        document.getElementById('usagePercent').textContent = stats.usage_percent.toFixed(1) + '%';
    } catch (e) {
        console.error('Failed to load stats:', e);
    }
}


/**
 * loadEmployees - Load and display employee list
 * 
 * Fetches all employees from the /api/employees endpoint and
 * renders them in a table. Each row includes:
 * - Employee ID (truncated)
 * - Full name
 * - Email address
 * - Department
 * - Position
 * - Action buttons (View, Delete)
 * 
 * Shows an empty state message if no employees exist.
 * 
 * @returns {Promise<void>}
 */
async function loadEmployees() {
    try {
        const employees = await fetchJSON(`${API_BASE}/employees`);
        const tbody = document.getElementById('employeeTable');
        
        if (employees.length === 0) {
            tbody.innerHTML = `<tr><td colspan="6" class="empty-state"><span class="empty-state-text">No employees found</span></td></tr>`;
            return;
        }

        tbody.innerHTML = employees.map(emp => `
            <tr class="animate-in">
                <td>${truncate(emp.id, 12)}</td>
                <td>${emp.name}</td>
                <td>${emp.email}</td>
                <td>${emp.department || '-'}</td>
                <td>${emp.position || '-'}</td>
                <td>
                    <button class="action-btn" onclick="viewEmployee('${emp.id}')">View</button>
                    <button class="action-btn" onclick="deleteEmployee('${emp.id}')">Delete</button>
                </td>
            </tr>
        `).join('');
    } catch (e) {
        console.error('Failed to load employees:', e);
    }
}


/**
 * loadTokens - Load and display token list
 * 
 * Fetches all tokens from the /api/tokens endpoint and
 * renders them in a table. Each row displays:
 * - Token ID (truncated)
 * - Employee ID (truncated)
 * - Token value (truncated)
 * - Total quota
 * - Used quota
 * - Status (Active/Revoked)
 * - Action button (Revoke, only for active tokens)
 * 
 * @returns {Promise<void>}
 */
async function loadTokens() {
    try {
        const tokens = await fetchJSON(`${API_BASE}/tokens`);
        const tbody = document.getElementById('tokenTable');
        
        if (tokens.length === 0) {
            tbody.innerHTML = `<tr><td colspan="7" class="empty-state"><span class="empty-state-text">No tokens found</span></td></tr>`;
            return;
        }

        tbody.innerHTML = tokens.map(tok => `
            <tr class="animate-in">
                <td>${truncate(tok.id, 12)}</td>
                <td>${truncate(tok.employee_id, 12)}</td>
                <td>${truncate(tok.token_value, 16)}</td>
                <td>${tok.total_quota}</td>
                <td>${tok.used_quota}</td>
                <td><span class="status-badge ${tok.is_active ? 'active' : 'revoked'}">${tok.is_active ? 'Active' : 'Revoked'}</span></td>
                <td>
                    ${tok.is_active ? `<button class="action-btn" onclick="revokeToken('${tok.id}')">Revoke</button>` : ''}
                </td>
            </tr>
        `).join('');
    } catch (e) {
        console.error('Failed to load tokens:', e);
    }
}


/**
 * loadMappings - Load and display token-employee mappings
 * 
 * Fetches all token-employee relationships from the /api/mappings endpoint
 * and displays them in a card grid format. Each card shows:
 * - Employee name
 * - Employee ID
 - Token value
 - Quota usage (remaining / total)
 * 
 * The mappings endpoint supports pagination via query parameters.
 * 
 * @returns {Promise<void>}
 */
async function loadMappings() {
    try {
        const data = await fetchJSON(`${API_BASE}/mappings`);
        const mappings = data.data || [];
        const grid = document.getElementById('mappingsGrid');
        
        if (mappings.length === 0) {
            grid.innerHTML = `<div class="empty-state" style="grid-column: 1/-1;"><span class="empty-state-text">No mappings found</span></div>`;
            return;
        }

        grid.innerHTML = mappings.map(m => `
            <div class="mapping-card animate-in">
                <div class="mapping-employee">${m.employee_name || 'Unknown'}</div>
                <div class="mapping-dept">${m.employee_id}</div>
                <div class="mapping-token">${truncate(m.token_value, 30)}</div>
                <div class="mapping-quota">
                    <span class="quota-label">Quota</span>
                    <span class="quota-value">${m.remaining} / ${m.total_quota}</span>
                </div>
            </div>
        `).join('');
    } catch (e) {
        console.error('Failed to load mappings:', e);
    }
}


/**
 * =============================================================================
 * 5. User Action Functions
 * Functions that handle user-initiated actions like delete, revoke, view.
 * =============================================================================
 */

/**
 * deleteEmployee - Delete an employee
 * 
 * @param {string} id - The ID of the employee to delete
 * 
 * Prompts for confirmation before sending a DELETE request to
 * /api/employees/:id. Upon success, reloads the employee list
 * and updates statistics.
 */
async function deleteEmployee(id) {
    if (!confirm('Delete this employee?')) return;
    try {
        await deleteResource(`${API_BASE}/employees/${id}`);
        loadEmployees();
        loadStats();
    } catch (e) {
        alert('Failed to delete employee: ' + e.message);
    }
}


/**
 * revokeToken - Revoke a token
 * 
 * @param {string} id - The ID of the token to revoke
 * 
 * Prompts for confirmation before sending a POST request to
 * /api/tokens/:id/revoke. Upon success, reloads the token list
 * and updates statistics. Revoked tokens cannot be used but
 * historical usage data is preserved.
 */
async function revokeToken(id) {
    if (!confirm('Revoke this token?')) return;
    try {
        await postJSON(`${API_BASE}/tokens/${id}/revoke`, {});
        loadTokens();
        loadStats();
    } catch (e) {
        alert('Failed to revoke token: ' + e.message);
    }
}


/**
 * viewEmployee - View detailed employee statistics
 * 
 * @param {string} id - The ID of the employee to view
 * 
 * Fetches detailed statistics for a specific employee from
 * /api/employees/:id/stats and displays them in an alert dialog.
 * The statistics include:
 * - Personal information (name, email, department, position)
 * - Token statistics (total tokens, active tokens)
 * - Usage statistics (used quota, total quota, usage percentage)
 */
async function viewEmployee(id) {
    try {
        const stats = await fetchJSON(`${API_BASE}/employees/${id}/stats`);
        const info = stats.employee_info;
        alert(`
Name: ${info.name}
Email: ${info.email}
Department: ${info.department}
Position: ${info.position}
Tokens: ${stats.token_count}
Active: ${stats.active_tokens}
Quota: ${stats.used_quota} / ${stats.total_quota}
Usage: ${stats.usage_percent.toFixed(1)}%
        `.trim());
    } catch (e) {
        alert('Failed to load employee details');
    }
}


/**
 * =============================================================================
 * 6. Form Handlers
 * Event handlers for form submissions.
 * =============================================================================
 */

/**
 * Employee Form Handler
 * 
 * Handles the submission of the employee creation form.
 * Collects form data and sends a POST request to /api/employees.
 * Upon success, clears the form and reloads the employee list.
 * 
 * Form fields:
 * - name (required): Employee's full name
 * - email (required): Employee's email address
 * - department: Employee's department
 * - position: Employee's job position
 */
document.getElementById('employeeForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const form = e.target;
    const data = {
        name: form.name.value,
        email: form.email.value,
        department: form.department.value,
        position: form.position.value
    };
    try {
        await postJSON(`${API_BASE}/employees`, data);
        form.reset();
        loadEmployees();
        loadStats();
    } catch (e) {
        alert('Failed to create employee: ' + e.message);
    }
});


/**
 * Token Form Handler
 * 
 * Handles the submission of the token issuance form.
 * Collects form data and sends a POST request to /api/tokens/issue.
 * Upon success, clears the form and reloads the token list.
 * 
 * Form fields:
 * - employee_id (required): ID of the employee to issue token to
 * - total_quota (required): Total token quota amount
 * - days_valid (required): Number of days the token is valid
 */
document.getElementById('tokenForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const form = e.target;
    const data = {
        employee_id: form.employee_id.value,
        total_quota: parseInt(form.total_quota.value),
        days_valid: parseInt(form.days_valid.value)
    };
    try {
        await postJSON(`${API_BASE}/tokens/issue`, data);
        form.reset();
        loadTokens();
        loadStats();
    } catch (e) {
        alert('Failed to issue token: ' + e.message);
    }
});


/**
 * =============================================================================
 * 7. Navigation Handlers
 * Event handlers for tab-based navigation.
 * =============================================================================
 */

/**
 * Tab Navigation Handler
 * 
 * Manages tab-based navigation between different views:
 * - overview: System statistics dashboard
 * - employees: Employee management (list and create)
 * - tokens: Token management (list and issue)
 * - mappings: Token-employee relationship view
 * 
 * Each tab click:
 * 1. Updates visual state (active class)
 * 2. Shows/hides appropriate content sections
 * 3. Loads relevant data from the API
 */
document.querySelectorAll('.nav-item').forEach(btn => {
    btn.addEventListener('click', () => {
        // Remove active state from all navigation buttons
        document.querySelectorAll('.nav-item').forEach(b => b.classList.remove('active'));
        
        // Hide all tab content
        document.querySelectorAll('.tab-content').forEach(t => t.classList.remove('active'));
        
        // Activate clicked button
        btn.classList.add('active');
        
        // Show corresponding tab content
        document.getElementById(btn.dataset.tab).classList.add('active');

        // Load data for the selected tab
        const tab = btn.dataset.tab;
        if (tab === 'overview') loadStats();
        else if (tab === 'employees') loadEmployees();
        else if (tab === 'tokens') loadTokens();
        else if (tab === 'mappings') loadMappings();
    });
});


/**
 * =============================================================================
 * 8. Initialization
 * Application startup and initial data loading.
 * =============================================================================
 */

// Initial load: Fetch and display system statistics when page loads
// This provides immediate feedback on the system state
loadStats();
