package router

import (
	"testing"

	"github.com/AnatoleLucet/traefik-test/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestRouter_Match(t *testing.T) {
	t.Run("matche by host", func(t *testing.T) {
		rules := []config.Rule{
			{
				If:   config.If{Host: "example.com", Path: "/users/*/friends", Method: "GET"},
				Then: config.Then{Forward: "http://localhost:8080"},
			},
			{
				If:   config.If{Host: "www.example.com", Path: "/users/*/friends", Method: "GET"},
				Then: config.Then{Forward: "http://localhost:8090"},
			},
		}

		rtr := New(rules)
		rule, ok := rtr.Match(Request{Host: "www.example.com", Path: "/users/123/friends", Method: "GET"})

		assert.True(t, ok)
		assert.Equal(t, rules[1], rule)
	})

	t.Run("match by path", func(t *testing.T) {
		rules := []config.Rule{
			{
				If:   config.If{Host: "example.com", Path: "/users/*/friends", Method: "GET"},
				Then: config.Then{Forward: "http://localhost:8080"},
			},
			{
				If:   config.If{Host: "example.com", Path: "/users/*/profile", Method: "GET"},
				Then: config.Then{Forward: "http://localhost:8090"},
			},
		}

		rtr := New(rules)
		rule, ok := rtr.Match(Request{Host: "example.com", Path: "/users/123/profile", Method: "GET"})

		assert.True(t, ok)
		assert.Equal(t, rules[1], rule)
	})

	t.Run("match by method", func(t *testing.T) {
		rules := []config.Rule{
			{
				If:   config.If{Host: "example.com", Path: "/users/*/friends", Method: "GET"},
				Then: config.Then{Forward: "http://localhost:8080"},
			},
			{
				If:   config.If{Host: "example.com", Path: "/users/*/friends", Method: "POST"},
				Then: config.Then{Forward: "http://localhost:8090"},
			},
		}

		rtr := New(rules)
		rule, ok := rtr.Match(Request{Host: "example.com", Path: "/users/123/friends", Method: "POST"})

		assert.True(t, ok)
		assert.Equal(t, rules[1], rule)
	})

	t.Run("match by host path and method", func(t *testing.T) {
		rules := []config.Rule{
			{
				If:   config.If{Host: "example.com", Path: "/users/*/profile", Method: "POST"},
				Then: config.Then{Forward: "http://localhost:8080"},
			},
			{
				If:   config.If{Host: "www.example.com", Path: "/users/*/friends", Method: "POST"},
				Then: config.Then{Forward: "http://localhost:8080"},
			},
			{
				If:   config.If{Host: "example.com", Path: "/users/*/profile", Method: "GET"},
				Then: config.Then{Forward: "http://localhost:8080"},
			},
			{
				If:   config.If{Host: "www.example.com", Path: "/users/*/profile", Method: "POST"},
				Then: config.Then{Forward: "http://localhost:8090"},
			},
		}

		rtr := New(rules)
		rule, ok := rtr.Match(Request{Host: "www.example.com", Path: "/users/123/profile", Method: "POST"})

		assert.True(t, ok)
		assert.Equal(t, rules[3], rule)
	})

	t.Run("match by specificity", func(t *testing.T) {
		rules := []config.Rule{
			{
				If:   config.If{Host: "www.example.com"},
				Then: config.Then{Forward: "http://localhost:8080"},
			},
			{
				If:   config.If{Host: "www.example.com", Path: "/users/*/profile"},
				Then: config.Then{Forward: "http://localhost:8080"},
			},
			{
				If:   config.If{Host: "www.example.com", Path: "/users/*/profile", Method: "POST"},
				Then: config.Then{Forward: "http://localhost:8090"},
			},
			// also a prefect match like above. used to make sure natural order is preserved.
			// only the rule above should match. not this one.
			{
				If:   config.If{Host: "www.example.com", Path: "/users/*/profile", Method: "POST"},
				Then: config.Then{Forward: "http://localhost:8070"},
			},
		}

		rtr := New(rules)
		rule, ok := rtr.Match(Request{Host: "www.example.com", Path: "/users/123/profile", Method: "POST"})

		assert.True(t, ok)
		assert.Equal(t, rules[2], rule)
	})

	t.Run("no match", func(t *testing.T) {
		rules := []config.Rule{
			{
				If:   config.If{Host: "example.com", Path: "/users/*/profile", Method: "POST"},
				Then: config.Then{Forward: "http://localhost:8080"},
			},
			{
				If:   config.If{Host: "www.example.com", Path: "/users/*/friends", Method: "POST"},
				Then: config.Then{Forward: "http://localhost:8090"},
			},
		}

		rtr := New(rules)
		rule, ok := rtr.Match(Request{Host: "www.example.com", Path: "/users/123/profile", Method: "POST"})

		assert.False(t, ok)
		assert.Equal(t, config.Rule{}, rule)
	})
}

func TestMatchSegments(t *testing.T) {
	t.Run("host", func(t *testing.T) {
		t.Run("exact match", func(t *testing.T) {
			tests := [][]string{
				// no tld
				{"localhost", "localhost"},

				// tld
				{"example.com", "example.com"},
				{"www.example.com", "www.example.com"},
				{"api.example.com", "api.example.com"},

				// ports
				{"localhost:8080", "localhost:8080"},
				{"example.com:8080", "example.com:8080"},
				{"www.example.com:8080", "www.example.com:8080"},
			}

			for _, c := range tests {
				assert.Truef(t, matchSegments(splitHost(c[0]), splitHost(c[1])), "expected '%s' to match '%s'", c[0], c[1])
			}
		})

		t.Run("wildcard match", func(t *testing.T) {
			tests := [][]string{
				// leading wildcard
				{"example.com", "*.com"},
				{"www.example.com", "*.com"},
				{"www.example.com", "*.example.com"},

				// trailing wildcard
				{"example.com", "example.*"},
				{"www.example.com", "www.*"},
				{"www.example.com", "www.example.*"},

				// middle wildcard
				{"api.example.com", "api.*.com"},
				{"api.dev.example.com", "api.*.com"},
			}

			for _, c := range tests {
				assert.Truef(t, matchSegments(splitHost(c[0]), splitHost(c[1])), "expected '%s' to match '%s'", c[0], c[1])
			}
		})

		t.Run("no match", func(t *testing.T) {
			tests := [][]string{
				{"example.com", "www.example.com"},
				{"www.example.com", "api.example.com"},

				// leading wildcard
				{"example.com", "*.example.org"},
				{"example.com", "*.org"},
				{"www.example.com", "*.org"},

				// trailing wildcard
				{"localhost", "localhost.*"},
				{"notexample.com", "example.*"},
				{"www.example.com", "example.*"},

				// middle wildcard
				{"api.example.com", "www.*.org"},
			}

			for _, c := range tests {
				assert.Falsef(t, matchSegments(splitHost(c[0]), splitHost(c[1])), "expected '%s' NOT to match '%s'", c[0], c[1])
			}
		})
	})

	t.Run("path", func(t *testing.T) {
		t.Run("exact match", func(t *testing.T) {
			tests := [][]string{
				{"/users", "/users"},
				{"/users/123", "/users/123"},
				{"/users/123/profile", "/users/123/profile"},
			}

			for _, c := range tests {
				assert.Truef(t, matchSegments(splitPath(c[0]), splitPath(c[1])), "expected '%s' to match '%s'", c[0], c[1])
			}
		})

		t.Run("wildcard match", func(t *testing.T) {
			tests := [][]string{
				// leading wildcard
				{"/users/profile", "*/profile"},
				{"/a/b/profile", "*/profile"},

				// trailing wildcard
				{"/users/123", "/users/*"},
				{"/users/123/profile", "/users/*"},

				// middle wildcard
				{"/users/123/profile", "/users/*/profile"},
				{"/users/123/abc/profile", "/users/*/profile"},
			}

			for _, c := range tests {
				assert.Truef(t, matchSegments(splitPath(c[0]), splitPath(c[1])), "expected '%s' to match '%s'", c[0], c[1])
			}
		})

		t.Run("no match", func(t *testing.T) {
			tests := [][]string{
				{"/users", "/api"},
				{"/users/123", "/users/456"},

				// leading wildcard
				{"/users", "*/users"},
				{"/users", "*/profile"},
				{"/users", "*/user/profile"},

				// trailing wildcard
				{"/users", "/users/*"},
				{"/api/users", "/users/*"},

				// middle wildcard
				{"/users/123/profile", "/users/*/settings"},
			}

			for _, c := range tests {
				assert.Falsef(t, matchSegments(splitPath(c[0]), splitPath(c[1])), "expected '%s' NOT to match '%s'", c[0], c[1])
			}
		})
	})
}

func TestMatchMethod(t *testing.T) {
	t.Run("exact match", func(t *testing.T) {
		tests := [][]string{
			{"GET", "GET"},
			{"POST", "POST"},
			{"PUT", "PUT,PATCH"},
		}

		for _, c := range tests {
			assert.Truef(t, matchMethod(c[0], splitMethod(c[1])), "expected '%s' to match '%s'", c[0], c[1])
		}
	})

	t.Run("no match", func(t *testing.T) {
		tests := [][]string{
			{"GET", "POST"},
			{"POST", "GET,PUT"},
			{"PATCH", "GET,POST,PUT"},
		}

		for _, c := range tests {
			assert.Falsef(t, matchMethod(c[0], splitMethod(c[1])), "expected '%s' NOT to match '%s'", c[0], c[1])
		}
	})
}
