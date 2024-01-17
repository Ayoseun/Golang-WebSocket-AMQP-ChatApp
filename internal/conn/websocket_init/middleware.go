package websocketinit

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// authMiddleware is a middleware function to handle authentication and authorization.
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get the expected Authorization header value from the environment
	uri := os.Getenv("AUTHHEADER")

	return func(w http.ResponseWriter, r *http.Request) {
		// Extract and decode the Authorization header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			// Return Forbidden status if Authorization header is missing
			http.Error(w, "Forbidden: Invalid Authorization header", http.StatusForbidden)
			return
		}

		// Perform authentication and authorization checks (replace this with your actual logic)
		if authHeader != uri {
			// Return Forbidden status if credentials are invalid
			http.Error(w, "Forbidden: Invalid credentials", http.StatusForbidden)
			return
		}

		// Call the next handler in the chain if authentication is successful
		next.ServeHTTP(w, r)
	}
}
