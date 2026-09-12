package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// writeJSONError sends all API errors in a consistent JSON format.
func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// expenseHandler sends each HTTP method to the correct handler.
func expenseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		getExpenses(w, r)

	case http.MethodPost:
		createExpense(w, r)

	case http.MethodPut:
		updateExpense(w, r)

	case http.MethodDelete:
		deleteExpense(w, r)

	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getExpenses returns all expenses.
func getExpenses(w http.ResponseWriter, r *http.Request) {
	// Ask PostgreSQL for all expenses.
	rows, err := db.Query("SELECT id, name, amount FROM expenses")
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Start with an empty slice so empty results return [] instead of null.
	expenses := []Expense{}

	// Move through each row PostgreSQL returned.
	for rows.Next() {
		var expense Expense

		err := rows.Scan(&expense.ID, &expense.Name, &expense.Amount)
		if err != nil {
			writeJSONError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Add this expense to our list.
		expenses = append(expenses, expense)
	}

	// Check whether an error occurred while iterating through rows.
	if err := rows.Err(); err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send completed response back to client.
	json.NewEncoder(w).Encode(expenses)
}

// createExpense creates a new expense.
func createExpense(w http.ResponseWriter, r *http.Request) {
	// Create an empty Expense.
	var newExpense Expense

	// Decode incoming JSON into newExpense.
	err := json.NewDecoder(r.Body).Decode(&newExpense)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate expense name.
	if newExpense.Name == "" {
		writeJSONError(w, "Expense name is required", http.StatusBadRequest)
		return
	}

	// Validate expense amount.
	if newExpense.Amount <= 0 {
		writeJSONError(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}

	// Insert expense into PostgreSQL and retrieve generated ID.
	err = db.QueryRow(
		"INSERT INTO expenses (name, amount) VALUES ($1, $2) RETURNING id",
		newExpense.Name,
		newExpense.Amount,
	).Scan(&newExpense.ID)

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Expense successfully created.
	w.WriteHeader(http.StatusCreated)

	// Send created expense back as JSON.
	json.NewEncoder(w).Encode(newExpense)
}

// updateExpense updates an existing expense.
func updateExpense(w http.ResponseWriter, r *http.Request) {
	// Read ID from URL.
	idString := r.URL.Query().Get("id")

	// Convert ID from string to integer.
	id, err := strconv.Atoi(idString)
	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Create an empty Expense.
	var newExpense Expense

	// Decode incoming JSON.
	err = json.NewDecoder(r.Body).Decode(&newExpense)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate expense name.
	if newExpense.Name == "" {
		writeJSONError(w, "Expense name is required", http.StatusBadRequest)
		return
	}

	// Validate expense amount.
	if newExpense.Amount <= 0 {
		writeJSONError(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}

	// Update expense in PostgreSQL.
	result, err := db.Exec(
		"UPDATE expenses SET name = $1, amount = $2 WHERE id = $3",
		newExpense.Name,
		newExpense.Amount,
		id,
	)

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check how many rows were updated.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// No matching expense existed.
	if rowsAffected == 0 {
		writeJSONError(w, "Expense not found", http.StatusNotFound)
		return
	}

	// Keep the ID from the URL.
	newExpense.ID = id

	// Send updated expense back as JSON.
	json.NewEncoder(w).Encode(newExpense)
}

// deleteExpense deletes an existing expense.
func deleteExpense(w http.ResponseWriter, r *http.Request) {
	// Read ID from URL.
	idString := r.URL.Query().Get("id")

	// Convert ID from string to integer.
	id, err := strconv.Atoi(idString)
	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Delete expense from PostgreSQL.
	result, err := db.Exec(
		"DELETE FROM expenses WHERE id = $1",
		id,
	)

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check whether anything was actually deleted.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// No matching expense existed.
	if rowsAffected == 0 {
		writeJSONError(w, "Expense not found", http.StatusNotFound)
		return
	}

	// Successful DELETE returns no body.
	w.WriteHeader(http.StatusNoContent)
}
