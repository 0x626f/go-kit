package env

import "strings"

const (
	defaultKey = "DEFAULT"
)

// Router provides a fluent API for conditional value resolution based on
// environment variable values. It allows you to define different values
// for different environment configurations (e.g., development, staging, production).
//
// The router performs case-insensitive matching on the environment variable value
// and returns the corresponding configured value. If no match is found and a default
// is configured, the default value is returned.
//
// Example:
//
//	// Configure database URL based on environment
//	dbURL := RouterOn("APP_ENV").
//	    InCase("development", "localhost:5432").
//	    InCase("staging", "staging-db.example.com:5432").
//	    InCase("production", "prod-db.example.com:5432").
//	    WithDefault("localhost:5432").
//	    Resolve()
type Router struct {
	variable string
	cases    map[string]string
}

// RouterOn creates a new Router that will resolve values based on the specified
// environment variable. The router uses case-insensitive matching for the
// environment variable value.
//
// Parameters:
//   - variable: The name of the environment variable to use for routing
//
// Returns:
//   - *Router: A new Router instance ready to be configured with cases
//
// Example:
//
//	router := RouterOn("ENVIRONMENT")
func RouterOn(variable string) *Router {
	return &Router{variable: variable, cases: make(map[string]string)}
}

// WithDefault sets the default value to return when the environment variable
// doesn't match any configured cases or when the environment variable is not set.
//
// This method should typically be called after all InCase calls to provide a
// fallback value. If no default is set and no case matches, Resolve will return
// an empty string.
//
// Parameters:
//   - value: The default value to return when no case matches
//
// Returns:
//   - *Router: The Router instance for method chaining
//
// Example:
//
//	router := RouterOn("ENV").
//	    InCase("prod", "production-value").
//	    WithDefault("default-value")
func (router *Router) WithDefault(value string) *Router {
	router.cases[defaultKey] = value
	return router
}

// InCase adds a conditional case to the router. When the environment variable
// matches the specified option (case-insensitive), the corresponding value will
// be returned by Resolve().
//
// Multiple cases can be chained together to handle different environment values.
// The option matching is case-insensitive, so "Production", "PRODUCTION", and
// "production" are treated as the same case.
//
// Parameters:
//   - option: The environment variable value to match (case-insensitive)
//   - value: The value to return when this case matches
//
// Returns:
//   - *Router: The Router instance for method chaining
//
// Example:
//
//	router := RouterOn("ENVIRONMENT").
//	    InCase("development", "dev-config").
//	    InCase("production", "prod-config")
func (router *Router) InCase(option, value string) *Router {
	router.cases[strings.ToLower(option)] = value
	return router
}

// Resolve evaluates the router and returns the value corresponding to the current
// environment variable value. It performs case-insensitive matching against all
// configured cases.
//
// Resolution order:
//  1. If the environment variable matches a configured case (case-insensitive),
//     return that case's value
//  2. If no match is found and a default is configured via WithDefault,
//     return the default value
//  3. If the environment variable is not set, use "DEFAULT" as the lookup key
//
// Returns:
//   - string: The resolved value based on the current environment variable value
//
// Example:
//
//	// Assuming ENV=production
//	apiURL := RouterOn("ENV").
//	    InCase("development", "http://localhost:8080").
//	    InCase("production", "https://api.example.com").
//	    WithDefault("http://localhost:8080").
//	    Resolve()
//	// Returns: "https://api.example.com"
func (router *Router) Resolve() string {
	currentCase := GetEnv(router.variable, defaultKey)
	if currentCase == defaultKey {
		return router.cases[defaultKey]
	}
	return router.cases[strings.ToLower(currentCase)]
}
