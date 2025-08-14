package utils

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/junwei890/se-cli/parsers"
)

func CheckAbility(visited map[string]struct{}, rules parsers.Rules, normURL string) bool {
	if _, ok := visited[normURL]; ok {
		return false
	} else {
		visited[normURL] = struct{}{}
	}

	green := true
	disallowedOn := ""
	allowedOn := ""

	for _, url := range rules.Disallowed {
		match, err := path.Match(url, normURL)
		if err != nil {
			log.Println(fmt.Errorf("can't match %s", url))
			continue
		}

		if !match {
			match = strings.HasPrefix(normURL, url)
		}

		if match {
			disallowedOn = url
			green = !match
			break
		}
	}

	for _, url := range rules.Allowed {
		match, err := path.Match(url, normURL)
		if err != nil {
			log.Println(fmt.Errorf("can't match %s", url))
			continue
		}

		if !match {
			match = strings.HasPrefix(normURL, url)
		}

		if match {
			allowedOn = url
			green = match
			break
		}
	}

	if disallowedOn != "" && allowedOn != "" {
		if len(disallowedOn) > len(allowedOn) {
			green = false
		} else {
			green = true
		}
	}

	return green
}

func CheckDomain(domain *url.URL, rawURL string) (bool, error) {
	structure, err := url.Parse(rawURL)
	if err != nil {
		return false, fmt.Errorf("can't parse %s", rawURL)
	}

	if structure.Hostname() != domain.Hostname() {
		return false, nil
	}

	return true, nil
}
