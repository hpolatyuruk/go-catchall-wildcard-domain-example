package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key string

const (
	// CustomerNameContextKey represents the key to get customer id from request context
	CustomerNameContextKey key = "CustomerName"
)

func parseCustomerName(host string) string {
	index := strings.Index(host, ".")
	if index < 0 {
		panic(fmt.Errorf("Unexpected host format %s", host))
	}
	return host[0:index]
}

func existsCustomer(name string) (bool, error) {
	// access database here to check whether customer exists or not by name
	return true, nil
}

func customerMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			customerName := parseCustomerName(r.Host)

			exists, err := existsCustomer(customerName)
			if err != nil {
				panic(err)
			}

			if exists == false {
				w.WriteHeader(404)
				w.Write([]byte("Customer not found!"))
				return
			}

			// set customer name to the request context because we will need customer name in our handlers
			ctx := context.WithValue(r.Context(), CustomerNameContextKey, customerName)

			// call next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var customerName string = r.Context().Value(CustomerNameContextKey).(string)

	// do customer related operations here. getting results from db by customer etc.

	w.Write([]byte(fmt.Sprintf("You are seeing %s's page", customerName)))
}

func main() {

	router := http.NewServeMux()

	router.HandleFunc("/", indexHandler)

	middleware := customerMiddleware()
	routerWithMiddleware := middleware(router)

	err := http.ListenAndServe(":80", routerWithMiddleware)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
	}
}
