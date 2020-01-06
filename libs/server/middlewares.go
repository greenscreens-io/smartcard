/**
 * Copyright (C) 2015, 2016  Green Screens Ltd.
 */

package server

import (
	"log"
	"time"
	"net/http"
	"net/http/httputil"
	"golang.org/x/time/rate"
)

// Middleware - Define middleware method
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Logging is Generic http request logging
// Logging logs all requests with its path and the time it took to process
func Logging(isDebug bool) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			if (isDebug) {

				requestDump, err := httputil.DumpRequest(r, true)

				if err != nil {
					log.Println(err)
				}

				log.Println(string(requestDump))
			}

			// Do middleware things
			start := time.Now()
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Cors Implement CORS browser ehaders
func Cors(allowedMethod string) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			// allow cross domain AJAX requests
			w.Header().Set("Content-Type", "application/json")
			if origin := r.Header.Get("Origin"); origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods",
					"POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers",
					"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			}

			// Stop here if its Preflighted OPTIONS request
			if r.Method == "OPTIONS" {
				return
			}

			if r.Method != allowedMethod {
				response := "Only " + allowedMethod + " requests are allowed"
				http.Error(w, response, http.StatusMethodNotAllowed)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
func Method(m string) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Limit - Requst limiter for protection
func Limit(limiter *rate.Limiter) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			if limiter.Allow() == false {
				http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
				return
			}

			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {

	for _, m := range middlewares {
		f = m(f)
	}

	return f
}
