package github

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/micnncim/bitbrew/plugin"
)

type Service interface {
	Search(ctx context.Context, q string) (plugin.Plugins, error)
	SearchByFilename(ctx context.Context, q string) (plugin.Plugins, error)
}

type githubSearchService interface {
	Code(ctx context.Context, query string, opt *github.SearchOptions) (*github.CodeSearchResult, *github.Response, error)
}

type service struct {
	githubSearchService
}

const (
	// baseGitHubRawURL: https://raw.githubusercontent.com/matryer/bitbar-plugins/master/Dev/Homebrew/brew-services.10m.rb
	baseGitHubRawURL = "https://raw.githubusercontent.com/%s/%s/master/%s"
	// baseBitBarURL: https://getbitbar.com/plugins/Dev/Homebrew/brew-services.10m.rb
	baseBitBarURL = "https://getbitbar.com/plugins/%s"
)

var (
	closedTagTitlePattern   = regexp.MustCompile(`<bitbar.title>.*</bitbar.title>`)
	leftTagTitlePattern     = regexp.MustCompile(`<bitbar.title>.*$`)
	closedTagDescPattern    = regexp.MustCompile(`<bitbar.desc>.*</bitbar.desc>`)
	leftTagDescPattern      = regexp.MustCompile(`<bitbar.desc>.*$`)
	closedTagPattern        = regexp.MustCompile(`<.+>(.+)<.+>`)
	leftHalfRightTagPattern = regexp.MustCompile(`<.*>(.+)<.*`)
	leftTagPattern          = regexp.MustCompile(`<.*>(.+)`)
)

var (
	extractByTagFunc = extractByTag
)

func NewService(token string) (Service, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	tc := oauth2.NewClient(context.Background(), ts)

	return &service{
		githubSearchService: github.NewClient(tc).Search,
	}, nil
}

func (s *service) Search(ctx context.Context, q string) (plugin.Plugins, error) {
	opt := github.ListOptions{PerPage: 20}
	var plugins plugin.Plugins

	for {
		result, resp, err := s.githubSearchService.Code(ctx, q, &github.SearchOptions{
			TextMatch:   true,
			ListOptions: opt,
		})
		if err != nil {
			return nil, err
		}

		for _, r := range result.CodeResults {
			textMatches := make([]string, 0, len(r.TextMatches))
			for _, match := range r.TextMatches {
				textMatches = append(textMatches, *match.Fragment)
			}

			name := extractByTagFunc(textMatches, closedTagTitlePattern, leftTagTitlePattern)
			if name == "" {
				continue
			}

			plugins = append(plugins, &plugin.Plugin{
				Name:         name,
				Filename:     *r.Name,
				Description:  extractByTagFunc(textMatches, closedTagDescPattern, leftTagDescPattern),
				Path:         *r.Path,
				BitBarURL:    fmt.Sprintf(baseBitBarURL, *r.Path),
				GitHubURL:    *r.HTMLURL,
				GitHubRawURL: fmt.Sprintf(baseGitHubRawURL, *r.Repository.Owner.Login, *r.Repository.Name, *r.Path),
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return plugins, nil
}

func (s *service) SearchByFilename(ctx context.Context, q string) (plugin.Plugins, error) {
	opt := github.ListOptions{PerPage: 20}
	var plugins plugin.Plugins

	for {
		result, resp, err := s.githubSearchService.Code(ctx, q, &github.SearchOptions{
			ListOptions: opt,
		})
		if err != nil {
			return nil, err
		}

		for _, r := range result.CodeResults {
			plugins = append(plugins, &plugin.Plugin{
				Filename:     *r.Name,
				Path:         *r.Path,
				BitBarURL:    fmt.Sprintf(baseBitBarURL, *r.Path),
				GitHubURL:    *r.HTMLURL,
				GitHubRawURL: fmt.Sprintf(baseGitHubRawURL, *r.Repository.Owner.Login, *r.Repository.Name, *r.Path),
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return plugins, nil
}

func extractByTag(candidates []string, closedPattern, leftPattern *regexp.Regexp) string {
	for _, c := range candidates {
		result := closedPattern.FindString(c)
		if result != "" {
			matches := closedTagPattern.FindStringSubmatch(result)
			if len(matches) > 1 {
				return matches[1]
			}
			return result
		}

		result = leftPattern.FindString(c)
		if result != "" {
			matches := leftHalfRightTagPattern.FindStringSubmatch(result)
			if len(matches) > 1 {
				return matches[1]
			}
			matches = leftTagPattern.FindStringSubmatch(result)
			if len(matches) > 1 {
				return matches[1]
			}
			return result
		}
	}

	return ""
}
