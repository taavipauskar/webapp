package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
)

type contextKey string

const contextUserKey contextKey = "user_ip"

func (app *application) ipFromContext(ctx context.Context) string {
	return ctx.Value(contextUserKey).(string)
}

func (app *application) addIPToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = context.Background()
		// get ip
		ip, err := getIP(r)
		if err != nil {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
			if len(ip) == 0 {
				ip = "unknown"
			}
			ctx = context.WithValue(r.Context(), contextUserKey, ip)
		} else {
			ctx = context.WithValue(r.Context(), contextUserKey, ip)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "unknown", err
	}
	log.Print(ip)
	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", fmt.Errorf("%q is not a valid IP:port", ip)
	}

	fwd := r.Header.Get("X-Forwarded-For")
	if len(fwd) == 0 {
		ip = fwd
	}

	if len(ip) == 0 {
		ip = "forward"
	}
	return ip, nil
}
