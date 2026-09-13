# expense_tracker
# Expense Tracker

A simple full-stack expense tracker built with Go and PostgreSQL. Very basic frontend to begin with.

This project was built to learn the fundamentals of backend development, including REST API, HTTP, database integration, validation, testing, configuration, and deployment.

## Live Demo

https://expense-tracker-w2mu.onrender.com/

> The application is hosted on Render's free tier, so the first request may take a short time if the server has been inactive.

## Features

- Create expenses
- View all expenses
- Update expenses
- Delete expenses
- Input validation
- JSON error responses
- Persistent PostgreSQL storage
- Basic automated HTTP handler tests

## Tech Stack

- Go
- PostgreSQL
- HTML 
- Neon
- Render

## API

### Get all expenses

`GET /expenses`

### Create an expense

`POST /expenses`

Example request:

```json
{
  "name": "Coffee",
  "amount": 5
}
