package mock

import "net/http"

type Check struct {
	Passed bool
}

func GenerateCheckedMiddlewares(checks ...*Check) []func(http.Handler) http.Handler {
	var handlers = make([]func(http.Handler) http.Handler, len(checks))
	for k, c := range checks {
		func(c *Check) { // Wrapping the function call like this creates a closure, making sure `c` does not change before evaluation.
			handlers[k] = func(handler http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.Passed = true
					handler.ServeHTTP(w, r)
				})
			}
		}(c)
	}
	return handlers
}
