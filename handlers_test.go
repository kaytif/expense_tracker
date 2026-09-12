package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test unsupported HTTP methods.
func TestExpenseHandler_MethodNotAllowed(t *testing.T) {
	tests := []struct {
		name   string
		method string
	}{
		{"PATCH", http.MethodPatch},
		{"OPTIONS", http.MethodOptions},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, "/expenses", nil)
			rec := httptest.NewRecorder()

			expenseHandler(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405, got %d", rec.Code)
			}
		})
	}
}

// Test invalid POST requests.
func TestCreateExpense_InvalidInput(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"empty name", `{"name":"","amount":5}`},
		{"zero amount", `{"name":"Coffee","amount":0}`},
		{"negative amount", `{"name":"Coffee","amount":-5}`},
		{"broken JSON", `{"name":"Coffee","amount":}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := strings.NewReader(test.body)

			req := httptest.NewRequest(http.MethodPost, "/expenses", body)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			expenseHandler(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", rec.Code)
			}
		})
	}
}

// Test PUT requests containing invalid IDs.
func TestUpdateExpense_InvalidID(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"missing ID", "/expenses"},
		{"negative ID", "/expenses?id=-1"},
		{"zero ID", "/expenses?id=0"},
		{"non-number ID", "/expenses?id=abc"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := strings.NewReader(`{"name":"Coffee","amount":5}`)

			req := httptest.NewRequest(http.MethodPut, test.url, body)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			expenseHandler(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", rec.Code)
			}
		})
	}
}

// Test DELETE requests containing invalid IDs.
func TestDeleteExpense_InvalidID(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"missing ID", "/expenses"},
		{"negative ID", "/expenses?id=-1"},
		{"zero ID", "/expenses?id=0"},
		{"non-number ID", "/expenses?id=abc"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, test.url, nil)
			rec := httptest.NewRecorder()

			expenseHandler(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", rec.Code)
			}
		})
	}
}