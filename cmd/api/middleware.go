package main

import (
	"fmt"
	"net/http"
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
