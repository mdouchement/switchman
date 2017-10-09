package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type (
	// A Config holds all the yaml configuration file.
	Config struct {
		// Listen is the binding interface of Switchman.
		Listen string `json:"listen"`
		// Rules are the forwarding rule used by switchman.
		Rules map[string]*Rule `json:"rules"`
	}

	// A Rewrite holds all stuff used for rewriting a Path from a Rule.
	Rewrite struct {
		once sync.Once
		re   *regexp.Regexp
		// From is a Regexp that is used to parse all matching routes.
		From string `json:"from"`
		// To is the pattern of the destination route.
		To string `json:"to"`
	}

	// A Rule holds the proxy pass rule.
	Rule struct {
		// Path is the raw path of the request.
		Path string `json:"-"`
		// Name is a descriptor of the rule.
		Name string `json:"name"`
		// Type is the rule's type.
		Type string `json:"type"`
		// URL is the endpoint of the destination server.
		URL string `json:"url"`
		// Rewrite is path rewriter stuff.
		Rewrite *Rewrite `json:"rewrite"`
	}
)

// String implements the Stringer interface.
func (rw *Rewrite) String() string {
	return fmt.Sprintf("From: %s - To: %s", rw.From, rw.To)
}

// Perform runs the rewrite on the given data.
func (rw *Rewrite) Perform(data string) string {
	rw.once.Do(func() {
		rw.re = regexp.MustCompile(rw.From)
	})

	ndst := string(rw.To) // Dup the string

	m := rw.re.FindStringSubmatch(data)
	for k, v := range MatcherLookup(m, rw.re) {
		ndst = strings.Replace(ndst, fmt.Sprintf("<%s>", k), v, 1)
	}

	return ndst
}

// MatcherLookup returns the map value of the named captures.
func MatcherLookup(match []string, re *regexp.Regexp) map[string]string {
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}
	return result
}
