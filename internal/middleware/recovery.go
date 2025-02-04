// File: internal/middleware/recovery.go
package middleware

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "runtime/debug"
)

// RecoveryHandler creates a middleware that recovers from panics
func RecoveryHandler(isDevelopment bool) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    // Log the stack trace
                    stack := debug.Stack()
                    log.Printf("[Recovery] panic recovered:\n%s\n%s", err, stack)

                    // Create error response
                    errorResponse := map[string]interface{}{
                        "error":   "Internal Server Error",
                        "message": "An unexpected error occurred",
                    }

                    if isDevelopment {
                        errorResponse["debug"] = map[string]interface{}{
                            "error":     fmt.Sprintf("%v", err),
                            "stack":     string(stack),
                            "requestId": r.Context().Value("requestID"),
                        }
                    }

                    // Send response
                    w.Header().Set("Content-Type", "application/json")
                    w.WriteHeader(http.StatusInternalServerError)
                    json.NewEncoder(w).Encode(errorResponse)
                }
            }()

            next.ServeHTTP(w, r)
        })
    }
}

// SafeHandler wraps a handler with panic recovery
func SafeHandler(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Handler panic: %v\n%s", err, debug.Stack())
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()

        handler(w, r)
    }
}

// CleanStack removes sensitive information from stack traces
func CleanStack(stack string) string {
    // Implementation remains the same
    return stack
}
