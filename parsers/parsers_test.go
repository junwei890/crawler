package parsers

import (
	"net/url"
	"os"
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
			name:         "F1: test case 1",
			input:        "http://www.hello.com/world",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "F1: test case 2",
			input:        "http://www.hello.com/world/",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "F1: test case 3",
			input:        "https://www.hello.com/world",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "F1: test case 4",
			input:        "https://www.hello.com/world/",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "F1: test case 5",
			input:        "https://www.hello.com/world?unit=testing",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "F1: test case 6",
			input:        "https://www.hello.com/world?unit=testing#foo",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "F1: test case 7",
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
		name:   "F2: test case 1",
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
		name: "F3: test case 1",
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
