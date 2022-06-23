package controller

import (
	"fmt"
	"net/http"
	"os"
	"xm/client"
	apiError "xm/error"
)

// protect makes sure that caller is authorized to make the call before invoking actual handler
func protect(ipLocationClient client.IPLocationClient, handlerFunc func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getUserIP(r)

		country, err := ipLocationClient.GetLocation(ip)
		if err != nil || country != originCountry() {
			fmt.Println(err)
			respondJSON(w, http.StatusUnauthorized, map[string]string{"error": apiError.ErrorCodeInvalidRequestOrigin})
			return
		}

		handlerFunc(w, r)
	}
}

// getUserIP gets ip address from request based on
// X-Real-Ip - fetches first true IP (if the requests sits behind multiple NAT sources/load balancer)
// X-Forwarded-For - if for some reason X-Real-Ip is blank and does not return response, get from X-Forwarded-For
// Remote Address - last resort (usually won't be reliable as this might be the last ip or if it is a naked http request to server ie no load balancer)
func getUserIP(r *http.Request) string {
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	return ipAddress
}

func originCountry() string {
	origin := os.Getenv("ORIGIN_COUNTRY")
	if len(origin) == 0 {
		origin = "CY"
		// todo change it to CY when pushed
	}
	return origin
}
