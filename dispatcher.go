package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/fatih/color"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Dispatch setups the forwarding rules in Echo engine.
func Dispatch(engine *echo.Echo, rules map[string]*Rule) error {
	for path, rule := range rules {
		if rule.Type != "proxy" {
			return fmt.Errorf("Unsupported type: %s", rule.Type)
		}
		rule.Path = path

		// Traces
		color.Red("%s:", rule.Name)
		color.White("  Matching on %s", rule.Path)
		if rule.Rewrite != nil {
			color.White("  Rewriting %s", rule.Rewrite)
		}
		color.White("  Proxifying on %s", rule.URL)

		u, err := url.Parse(rule.URL)
		if err != nil {
			return errors.Wrap(err, "dispatch:")
		}
		handler := httputil.NewSingleHostReverseProxy(u) // Proxy

		rw := rule.Rewrite
		director := handler.Director
		handler.Director = func(req *http.Request) {
			color.White("Recieved: %s", req.URL.String())

			director(req)
			if rw != nil {
				req.URL.Path = rw.Perform(req.URL.Path)
			}

			color.White("  %s  %s", color.GreenString("⇨"), req.URL.String())
		}

		// Match rule
		engine.Any(path, echo.WrapHandler(handler))
	}

	return nil
}
