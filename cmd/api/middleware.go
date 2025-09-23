package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (a *application) recoverPanic(next http.Handler) http.Handler {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
   // defer will be called when the stack unwinds
       defer func() {
           // recover() checks for panics
           err := recover();
           if err != nil {
               w.Header().Set("Connection", "close")
               a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
           }
       }()
       next.ServeHTTP(w,r)
   })  
}

func (a *application) enableCORS (next http.Handler) http.Handler {                             
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		// Let's check the request origin to see if it's in the trusted list
		origin := r.Header.Get("Origin")
		
		// Once we have a origin from the request header we need need to check
		if origin != "" {
			for i := range a.config.cors.trustedOrigins {
				if origin == a.config.cors.trustedOrigins[i] {
				w.Header().Set("Access-Control-Allow-Origin", origin)                                          
				break
				}
		}
		}


        next.ServeHTTP(w, r)
    })
	
}

func (a *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter *rate.Limiter
		lastSeen time.Time
	}

	var mu sync.Mutex
	var clients = make(map[string] *client)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3 * time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.config.limiter.enabled {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			a.serverErrorResponse(w, r, err)
			return
		}
		
		mu.Lock()

		_, found := clients[ip]
		if !found {
		clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(a.config.limiter.rps), a.config.limiter.burst)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			a.rateLimitExceededResponse(w, r)
			return
		}
	
		mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}