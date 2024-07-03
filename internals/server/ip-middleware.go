package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type authHandler func(w http.ResponseWriter, r *http.Request) error

func middlewareWhitelistedIP(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-FORWARDED-FOR")
		log.Printf("\nIP Address: %v\n", r.RemoteAddr)

		whitelistedIPs := os.Getenv("WHITELISTED_IPS")
		isGlobalIP := strings.Contains(whitelistedIPs, "0.0.0.0")

		if !isGlobalIP && !strings.Contains(whitelistedIPs, ip) {
			responseWithError(w, 401, fmt.Sprintf("Invalid Caller: %v", ip))
			return
		}

		handler(w, r)
	}
}
