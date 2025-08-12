package utils

import (
	"net/url"
	"testing"

	"github.com/junwei890/se-cli/parsers"
)

func TestCheckAbility(t *testing.T) {
	testCases := []struct {
		name     string
		visited  map[string]struct{}
		rules    parsers.Rules
		normURL  string
		expected bool
	}{
		{
			name: "F4: test case 1",
			visited: map[string]struct{}{
				"www.google.com/places": {},
			},
			rules:    parsers.Rules{},
			normURL:  "www.google.com/places",
			expected: false,
		},
		{
			name:     "F4: test case 2",
			visited:  map[string]struct{}{},
			rules:    parsers.Rules{},
			normURL:  "www.google.com/places",
			expected: true,
		},
		{
			name:    "F4: test case 3",
			visited: map[string]struct{}{},
			rules: parsers.Rules{
				Disallowed: []string{
					"www.google.com/maps",
				},
			},
			normURL:  "www.google.com/maps",
			expected: false,
		},
		{
			name:    "F4: test case 4",
			visited: map[string]struct{}{},
			rules: parsers.Rules{
				Disallowed: []string{
					"www.google.com/maps/",
				},
			},
			normURL:  "www.google.com/maps/place",
			expected: false,
		},
		{
			name:     "F4: test case 5",
			visited:  map[string]struct{}{},
			rules:    parsers.Rules{},
			normURL:  "www.google.com/maps",
			expected: true,
		},
		{
			name:    "F4: test case 6",
			visited: map[string]struct{}{},
			rules: parsers.Rules{
				Disallowed: []string{
					"www.google.com/*world",
				},
			},
			normURL:  "www.google.com/helloworld",
			expected: false,
		},
		{
			name:    "F4: test case 7",
			visited: map[string]struct{}{},
			rules: parsers.Rules{
				Disallowed: []string{
					"www.google.com/hello*",
				},
			},
			normURL:  "www.google.com/helloworld",
			expected: false,
		},
		{
			name:    "F4: test case 8",
			visited: map[string]struct{}{},
			rules: parsers.Rules{
				Disallowed: []string{
					"www.google.com/maps/",
				},
				Allowed: []string{
					"www.google.com/maps/places",
				},
			},
			normURL:  "www.google.com/maps/places",
			expected: true,
		},
		{
			name:    "F4: test case 9",
			visited: map[string]struct{}{},
			rules: parsers.Rules{
				Disallowed: []string{
					"www.google.com/maps/places",
				},
				Allowed: []string{
					"www.google.com/maps/",
				},
			},
			normURL:  "www.google.com/maps/places",
			expected: false,
		},
		{
			name:    "F4: test case 10",
			visited: map[string]struct{}{},
			rules: parsers.Rules{
				Disallowed: []string{
					"www.google.com/maps/",
				},
				Allowed: []string{
					"www.google.com/maps/",
				},
			},
			normURL:  "www.google.com/maps/places",
			expected: true,
		},
		{
			name:    "F4: test case 11",
			visited: map[string]struct{}{},
			rules: parsers.Rules{
				Disallowed: []string{
					"www.google.com/maps/places/",
				},
				Allowed: []string{
					"www.google.com/maps",
				},
			},
			normURL:  "www.google.com/maps/places/oregon",
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if comp := CheckAbility(testCase.visited, testCase.rules, testCase.normURL); comp != testCase.expected {
				t.Errorf("%s failed, %t != %t", testCase.name, comp, testCase.expected)
			}
		})
	}
}

func TestCheckDomain(t *testing.T) {
	dom, err := url.Parse("https://www.google.com")
	if err != nil {
		t.Errorf("error setting up test, unexpected error: %v", err)
	}

	testCases := []struct {
		name         string
		domain       *url.URL
		rawURL       string
		expected     bool
		errorPresent bool
	}{
		{
			name:         "F5: test case 1",
			domain:       dom,
			rawURL:       "https://gasdfas ",
			expected:     false,
			errorPresent: true,
		},
		{
			name:         "F5: test case 2",
			domain:       dom,
			rawURL:       "https://www.google.com/maps",
			expected:     true,
			errorPresent: false,
		},
		{
			name:         "F5: test case 3",
			domain:       dom,
			rawURL:       "https://www.github.com/repos",
			expected:     false,
			errorPresent: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := CheckDomain(testCase.domain, testCase.rawURL)
			if (err != nil) != testCase.errorPresent {
				t.Errorf("%s failed, expected error: %v", testCase.name, err)
			}
			if result != testCase.expected {
				t.Errorf("%s failed, %v != %v", testCase.name, result, testCase.expected)
			}
		})
	}
}
