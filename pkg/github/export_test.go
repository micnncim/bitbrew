package github

import (
	"regexp"
)

var (
	ExportExtractByTag = extractByTag
)

const (
	ExportBaseGitHubRawURL = baseGitHubRawURL
	ExportBaseBitBarURL    = baseBitBarURL
)

type (
	ExportService = service
)

func (s *ExportService) ExportSetGithubSearchService(searchService githubSearchService) {
	s.githubSearchService = searchService
}

func ExportSetExtractByTag(f func(candidates []string, closedPattern, leftPattern *regexp.Regexp) string) (resetFunc func()) {
	var org func(candidates []string, closedPattern, leftPattern *regexp.Regexp) string
	org, extractByTagFunc = extractByTagFunc, f
	return func() {
		extractByTagFunc = org
	}
}
