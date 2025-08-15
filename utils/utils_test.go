package utils

import (
	"errors"
	"net/url"
	"os"
	"reflect"
	"slices"
	"testing"
)

func TestNormalize(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expected     string
		errorPresent bool
	}{
		{
			name:         "Normalize: test case 1",
			input:        "http://www.hello.com/world",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "Normalize: test case 2",
			input:        "http://www.hello.com/world/",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "Normalize: test case 3",
			input:        "https://www.hello.com/world",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "Normalize: test case 4",
			input:        "https://www.hello.com/world/",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "Normalize: test case 5",
			input:        "https://www.hello.com/world?unit=testing",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "Normalize: test case 6",
			input:        "https://www.hello.com/world?unit=testing#foo",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "Normalize: test case 7",
			input:        "com.invalid .https://",
			expected:     "",
			errorPresent: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := Normalize(testCase.input)
			if (err != nil) != testCase.errorPresent {
				t.Errorf("%s failed, unexpected error: %v", testCase.name, err)
			}
			if result != testCase.expected {
				t.Errorf("%s failed, %s != %s", testCase.name, result, testCase.expected)
			}
		})
	}
}

func TestParseHTML(t *testing.T) {
	page, err := os.ReadFile("./test_files/example.html")
	if err != nil {
		t.Errorf("error setting up test, unexpected error: %v", err)
	}

	domain, err := url.Parse("https://www.google.com")
	if err != nil {
		t.Errorf("error setting up test, unexpected error: %v", err)
	}

	testCase := struct {
		name     string
		domain   *url.URL
		page     []byte
		expected Response
	}{
		name:   "ParseHTML: test case 1",
		domain: domain,
		page:   page,
		expected: Response{
			Title: "Mixed Links Example",
			Content: []string{
				"home",
				"|",
				"services",
				"|",
				"github",
				"this site has a mix of internal and external links for demonstration purposes.",
				"learn more",
				"about us",
				"or check out our",
				"portfolio",
				".",
				"visit our",
				"documentation",
				"or read the latest",
				"tech news",
				".",
				"questions? reach out via our",
				"contact page",
				".",
				"© 2025 mixed link example",
			},
			Links: []string{
				"https://www.google.com/",
				"https://www.google.com/services",
				"https://www.github.com",
				"https://www.google.com/about",
				"https://www.example.com/portfolio",
				"https://www.google.com/docs",
				"https://news.ycombinator.com",
				"https://www.google.com/contact",
			},
		},
	}

	t.Run(testCase.name, func(t *testing.T) {
		result, err := ParseHTML(testCase.domain, testCase.page)
		if err != nil {
			t.Errorf("%s failed, unexpected error: %v", testCase.name, err)
		}
		if result.Title != testCase.expected.Title {
			t.Errorf("%s failed, %s != %s", testCase.name, result.Title, testCase.expected.Title)
		}
		if comp := slices.Equal(result.Content, testCase.expected.Content); !comp {
			t.Errorf("%s failed, %v != %v", testCase.name, result.Content, testCase.expected.Content)
		}
		if comp := slices.Equal(result.Links, testCase.expected.Links); !comp {
			t.Errorf("%s failed, %v != %v", testCase.name, result.Links, testCase.expected.Links)
		}
	})
}

func TestParseRobots(t *testing.T) {
	textFile, err := os.ReadFile("./test_files/example.txt")
	if err != nil {
		t.Errorf("error setting up test, unexpected error: %v", err)
	}

	testCase := struct {
		name     string
		url      string
		file     []byte
		expected Rules
	}{
		name: "ParseRobots: test case 1",
		url:  "www.google.com",
		file: textFile,
		expected: Rules{
			Allowed: []string{
				"www.google.com/archive",
				"www.google.com/year",
				"www.google.com/list",
				"www.google.com/abs",
				"www.google.com/pdf",
				"www.google.com/html",
				"www.google.com/catchup",
			},
			Disallowed: []string{
				"www.google.com/user",
				"www.google.com/e-print",
				"www.google.com/src",
				"www.google.com/ps",
				"www.google.com/dvi",
				"www.google.com/cookies",
				"www.google.com/form",
				"www.google.com/find",
				"www.google.com/view",
				"www.google.com/ftp",
				"www.google.com/refs",
				"www.google.com/cits",
				"www.google.com/format",
				"www.google.com/PS_cache",
				"www.google.com/Stats",
				"www.google.com/seek-and-destroy",
				"www.google.com/IgnoreMe",
				"www.google.com/oai2",
				"www.google.com/auth",
				"www.google.com/tb",
				"www.google.com/tb-recent",
				"www.google.com/trackback",
				"www.google.com/prevnext",
				"www.google.com/ct",
				"www.google.com/api",
				"www.google.com/search",
				"www.google.com/set_author_id",
				"www.google.com/show-email",
			},
			Delay: 15,
		},
	}

	t.Run(testCase.name, func(t *testing.T) {
		result, err := ParseRobots(testCase.url, testCase.file)
		if err != nil {
			t.Errorf("%s failed, unexpected error: %v", testCase.name, err)
		}
		if comp := slices.Equal(result.Allowed, testCase.expected.Allowed); !comp {
			t.Errorf("%s failed, %v != %v", testCase.name, result.Allowed, testCase.expected.Allowed)
		}
		if comp := slices.Equal(result.Disallowed, testCase.expected.Disallowed); !comp {
			t.Errorf("%s failed, %v != %v", testCase.name, result.Disallowed, testCase.expected.Disallowed)
		}
		if result.Delay != testCase.expected.Delay {
			t.Errorf("%s failed, %v != %v", testCase.name, result.Delay, testCase.expected.Delay)
		}
	})
}

func TestCheckAbility(t *testing.T) {
	testCases := []struct {
		name     string
		visited  map[string]struct{}
		rules    Rules
		normURL  string
		expected bool
	}{
		{
			name: "CheckAbility: test case 1",
			visited: map[string]struct{}{
				"www.google.com/places": {},
			},
			rules:    Rules{},
			normURL:  "www.google.com/places",
			expected: false,
		},
		{
			name:     "CheckAbility: test case 2",
			visited:  map[string]struct{}{},
			rules:    Rules{},
			normURL:  "www.google.com/places",
			expected: true,
		},
		{
			name:    "CheckAbility: test case 3",
			visited: map[string]struct{}{},
			rules: Rules{
				Disallowed: []string{
					"www.google.com/maps",
				},
			},
			normURL:  "www.google.com/maps",
			expected: false,
		},
		{
			name:    "CheckAbility: test case 4",
			visited: map[string]struct{}{},
			rules: Rules{
				Disallowed: []string{
					"www.google.com/maps/",
				},
			},
			normURL:  "www.google.com/maps/place",
			expected: false,
		},
		{
			name:     "CheckAbility: test case 5",
			visited:  map[string]struct{}{},
			rules:    Rules{},
			normURL:  "www.google.com/maps",
			expected: true,
		},
		{
			name:    "CheckAbility: test case 6",
			visited: map[string]struct{}{},
			rules: Rules{
				Disallowed: []string{
					"www.google.com/*world",
				},
			},
			normURL:  "www.google.com/helloworld",
			expected: false,
		},
		{
			name:    "CheckAbility: test case 7",
			visited: map[string]struct{}{},
			rules: Rules{
				Disallowed: []string{
					"www.google.com/hello*",
				},
			},
			normURL:  "www.google.com/helloworld",
			expected: false,
		},
		{
			name:    "CheckAbility: test case 8",
			visited: map[string]struct{}{},
			rules: Rules{
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
			name:    "CheckAbility: test case 9",
			visited: map[string]struct{}{},
			rules: Rules{
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
			name:    "CheckAbility: test case 10",
			visited: map[string]struct{}{},
			rules: Rules{
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
			name:    "CheckAbility: test case 11",
			visited: map[string]struct{}{},
			rules: Rules{
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
			name:         "CheckDomain: test case 1",
			domain:       dom,
			rawURL:       "https://gasdfas ",
			expected:     false,
			errorPresent: true,
		},
		{
			name:         "CheckDomain: test case 2",
			domain:       dom,
			rawURL:       "https://www.google.com/maps",
			expected:     true,
			errorPresent: false,
		},
		{
			name:         "CheckDomain: test case 3",
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

func TestQueue(t *testing.T) {
	queue := &Queue{"a", "b", "c", "d", "e"}

	size := queue.Size()
	if size != 5 {
		t.Errorf("Queue: test case 1 failed, %d != %d", size, 5)
	}

	queue.Enqueue("f")
	if comp := reflect.DeepEqual(*queue, Queue{"a", "b", "c", "d", "e", "f"}); !comp {
		t.Errorf("Queue: test case 2 failed: %v != %v", *queue, Queue{"a", "b", "c", "d", "e", "f"})
	}

	size = queue.Size()
	if size != 6 {
		t.Errorf("Queue: test case 3 failed, %d != %d", size, 6)
	}

	popped, err := queue.Dequeue()
	if err != nil {
		t.Errorf("Queue: test case 4 failed, unexpected error: %v", err)
	}
	if popped != "a" {
		t.Errorf("Queue: test case 5 failed, %s != %s", popped, "a")
	}
	if comp := reflect.DeepEqual(*queue, Queue{"b", "c", "d", "e", "f"}); !comp {
		t.Errorf("Queue: test case 6 failed, %v != %v", *queue, Queue{"b", "c", "d", "e", "f"})
	}

	first, err := queue.Peek()
	if err != nil {
		t.Errorf("Queue: test case 7 failed, unexpected error: %v", err)
	}
	if first != "b" {
		t.Errorf("Queue: test case 8 failed, %s != %s", first, "b")
	}

	for size := queue.Size(); size > 0; size-- {
		queue.Dequeue()
	}

	empty := queue.CheckEmpty()
	if !empty {
		t.Errorf("Queue: test case 9 failed, %v != %v", empty, true)
	}

	_, err = queue.Dequeue()
	if (err != nil) != true {
		t.Errorf("Queue: test case 10 failed, expected error: %s", errors.New("queue empty"))
	}

	_, err = queue.Peek()
	if (err != nil) != true {
		t.Errorf("Queue: test case 11 failed, expected error: %s", errors.New("queue empty"))
	}
}
