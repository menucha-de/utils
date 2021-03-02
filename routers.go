package utils

import (
	"net/http"
)

// Route a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes a list of routes
type Routes []Route
