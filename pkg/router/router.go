package router

import (
	"slices"
	"strings"

	"github.com/AnatoleLucet/traefik-test/pkg/config"
)

// NOTE: this router *could* definitly be optimized further with precompiled rules
// e.g. in a trie for O(1) matching. but this O(n) implementation is already fast enough and
// somewhat optimized. no need to overcomplicate the thing now :)

type Request struct {
	Host   string
	Path   string
	Method string
}

// Router is responsible for finding the correct rule to apply for a given request.
type Router struct {
	rules   []CompiledRule // [i]rule
	hosts   [][]string     // [i][sub, domain, com]
	paths   [][]string     // [i][user, *, friends]
	methods [][]string     // [i][GET, POST, PUT]
}

func New(rules []CompiledRule) *Router {
	// rules are sorted by specificity.
	// if multiple rules match the given request, the rule with the most specificity
	// will win (e.g. a `path+method` condition with always win over a single `path` condition).
	rules = slices.Clone(rules)
	slices.SortStableFunc(rules, func(a, b CompiledRule) int {
		return specificity(b.Rule) - specificity(a.Rule)
	})

	// store the rules as an soa to gain a few cycles by not pulling the
	// whole config.Rule struct in cache lines when iterating to find a match.
	hosts := make([][]string, len(rules))
	paths := make([][]string, len(rules))
	methods := make([][]string, len(rules))
	for i, rule := range rules {
		hosts[i] = splitHost(rule.Rule.If.Host)
		paths[i] = splitPath(rule.Rule.If.Path)
		methods[i] = splitMethod(rule.Rule.If.Method)
	}

	return &Router{
		rules:   rules,
		hosts:   hosts,
		paths:   paths,
		methods: methods,
	}
}

// Match finds the rule that matches best the given request.
func (r *Router) Match(req Request) (CompiledRule, bool) {
	host := splitHost(req.Host)
	path := splitPath(req.Path)

	for i := range r.rules {
		if matchSegments(host, r.hosts[i]) &&
			matchSegments(path, r.paths[i]) &&
			matchMethod(req.Method, r.methods[i]) {
			return r.rules[i], true
		}
	}

	return CompiledRule{}, false
}

// "www.example.com" -> ["com", "example", "www"]
func splitHost(host string) []string {
	trimmed := strings.Trim(host, ".")
	if trimmed == "" {
		return nil
	}

	segements := strings.Split(trimmed, ".")
	slices.Reverse(segements) // reverse to match tld first (e.g. [com, domain, sub])
	return segements
}

// "/users/*/friends" -> ["users", "*", "friends"]
func splitPath(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}

	return strings.Split(trimmed, "/")
}

// "GET,POST" -> ["GET", "POST"]
func splitMethod(method string) []string {
	trimmed := strings.Trim(method, ",")
	if trimmed == "" {
		return nil
	}

	return strings.Split(trimmed, ",")
}

// compares segments with wildcard support.
// e.g. source [com, example, www] matches target [com, *, www] but not [com, example, api]
func matchSegments(source, target []string) bool {
	if len(target) == 0 {
		return true // no constraints
	}

	// fast exit path for exact match
	if len(source) != len(target) && !slices.Contains(target, "*") {
		return false
	}

	si, ti := 0, 0
	for si < len(source) && ti < len(target) {
		if target[ti] == "*" {
			ti++
			si++

			// trailing wildcard, rest of the segments match
			if ti == len(target) {
				return true
			}

			// consume segments until we find the next segment in the target
			for si < len(source) && source[si] != target[ti] {
				si++
			}

			continue
		}

		// segment doesn't match
		if source[si] != target[ti] {
			return false
		}
		si++
		ti++
	}

	// make sure we traversed every segments
	return si == len(source) && ti == len(target)
}

// checks if source is in the list of allowed methods
// e.g. source "GET" matches target ["GET", "POST"] but not ["POST", "PUT"]
func matchMethod(source string, target []string) bool {
	if len(target) == 0 {
		return true // no constraints
	}

	return slices.Contains(target, source)
}

// calculates the specificity score of a rule. higher score means more specific.
func specificity(rule config.Rule) int {
	score := 0

	// 1000 for exact host segment, 100 for wildcard host segment
	for _, seg := range splitHost(rule.If.Host) {
		if seg == "*" {
			score += 100
		} else {
			score += 1000
		}
	}

	// 10 for exact path segment, 1 for wildcard path segment
	for _, seg := range splitPath(rule.If.Path) {
		if seg == "*" {
			score += 1
		} else {
			score += 10
		}
	}

	// 5 for method constraint
	if rule.If.Method != "" {
		score += 5
	}

	return score
}
