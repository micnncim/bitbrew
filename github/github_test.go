package github_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	gogithub "github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"

	"github.com/micnncim/bitbrew/github"
	"github.com/micnncim/bitbrew/plugin"
)

type fakeGithubSearchService struct {
	code func(ctx context.Context, query string, opt *gogithub.SearchOptions) (*gogithub.CodeSearchResult, *gogithub.Response, error)
}

func (s *fakeGithubSearchService) Code(ctx context.Context, query string, opt *gogithub.SearchOptions) (*gogithub.CodeSearchResult, *gogithub.Response, error) {
	return s.code(ctx, query, opt)
}

func Test_extractDesc(t *testing.T) {
	closedTagDescPattern := regexp.MustCompile(`<bitbar.desc>.*</bitbar.desc>`)
	leftTagDescPattern := regexp.MustCompile(`<bitbar.desc>.*$`)

	cases := []struct {
		name       string
		candidates []string
		want       string
	}{
		{
			name: "found desc in closed tag",
			candidates: []string{
				"#!/usr/bin/env ruby\n\n# <bitbar.title>Brew Services</bitbar.title>",
				"<bitbar.desc>Shows and manages Homebrew services.</bitbar.desc>",
			},
			want: "Shows and manages Homebrew services.",
		},
		{
			name: "found desc in left tag",
			candidates: []string{
				"#!/usr/bin/env ruby\n\n# <bitbar.title>Brew Services</bitbar.title>",
				"<bitbar.desc>Shows and manages Homebrew services.",
			},
			want: "Shows and manages Homebrew services.",
		},
		{
			name: "not found desc",
			candidates: []string{
				"#!/usr/bin/env ruby\n\n# <bitbar.title>Brew Services</bitbar.title>",
			},
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := github.ExportExtractByTag(tc.candidates, closedTagDescPattern, leftTagDescPattern)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_service_Search(t *testing.T) {
	cases := []struct {
		name             string
		codeFunc         func(ctx context.Context, query string, opt *gogithub.SearchOptions) (*gogithub.CodeSearchResult, *gogithub.Response, error)
		extractByTagFunc func(candidates []string, closedPattern, leftPattern *regexp.Regexp) string
		want             plugin.Plugins
		wantErr          bool
	}{
		{
			name: "found plugins",
			codeFunc: func(ctx context.Context, query string, opt *gogithub.SearchOptions) (*gogithub.CodeSearchResult, *gogithub.Response, error) {
				name, path, url := "filename", "path", "url"
				ownerName, repoName := "owner", "repo"
				return &gogithub.CodeSearchResult{
					CodeResults: []gogithub.CodeResult{
						{
							Name:    &name,
							Path:    &path,
							HTMLURL: &url,
							Repository: &gogithub.Repository{
								Name: &repoName,
								Owner: &gogithub.User{
									Login: &ownerName,
								},
							},
						},
						{
							Name:    &name,
							Path:    &path,
							HTMLURL: &url,
							Repository: &gogithub.Repository{
								Name: &repoName,
								Owner: &gogithub.User{
									Login: &ownerName,
								},
							},
						},
					},
				}, &gogithub.Response{}, nil
			},
			extractByTagFunc: func(candidates []string, closedPattern, leftPattern *regexp.Regexp) string {
				return "text"
			},
			want: plugin.Plugins{
				{
					Name:         "text",
					Filename:     "filename",
					Description:  "text",
					Path:         "path",
					BitBarURL:    fmt.Sprintf(github.ExportBaseBitBarURL, "path"),
					GitHubURL:    "url",
					GitHubRawURL: fmt.Sprintf(github.ExportBaseGitHubRawURL, "owner", "repo", "path"),
				},
				{
					Name:         "text",
					Filename:     "filename",
					Description:  "text",
					Path:         "path",
					BitBarURL:    fmt.Sprintf(github.ExportBaseBitBarURL, "path"),
					GitHubURL:    "url",
					GitHubRawURL: fmt.Sprintf(github.ExportBaseGitHubRawURL, "owner", "repo", "path"),
				},
			},
			wantErr: false,
		},
		{
			name: "github client error",
			codeFunc: func(ctx context.Context, query string, opt *gogithub.SearchOptions) (*gogithub.CodeSearchResult, *gogithub.Response, error) {
				return nil, nil, errors.New("error")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		s := new(github.ExportService)
		s.ExportSetGithubSearchService(&fakeGithubSearchService{
			code: tc.codeFunc,
		})

		t.Run(tc.name, func(t *testing.T) {
			reset := github.ExportSetExtractByTag(tc.extractByTagFunc)

			got, err := s.Search(context.Background(), "q")
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err != nil)

			reset()
		})
	}
}

func Test_service_SearchByFilename(t *testing.T) {
	cases := []struct {
		name     string
		codeFunc func(ctx context.Context, query string, opt *gogithub.SearchOptions) (*gogithub.CodeSearchResult, *gogithub.Response, error)
		want     plugin.Plugins
		wantErr  bool
	}{
		{
			name: "found plugins",
			codeFunc: func(ctx context.Context, query string, opt *gogithub.SearchOptions) (*gogithub.CodeSearchResult, *gogithub.Response, error) {
				name, path, url := "filename", "path", "url"
				ownerName, repoName := "owner", "repo"
				return &gogithub.CodeSearchResult{
					CodeResults: []gogithub.CodeResult{
						{
							Name:    &name,
							Path:    &path,
							HTMLURL: &url,
							Repository: &gogithub.Repository{
								Name: &repoName,
								Owner: &gogithub.User{
									Login: &ownerName,
								},
							},
						},
						{
							Name:    &name,
							Path:    &path,
							HTMLURL: &url,
							Repository: &gogithub.Repository{
								Name: &repoName,
								Owner: &gogithub.User{
									Login: &ownerName,
								},
							},
						},
					},
				}, &gogithub.Response{}, nil
			},
			want: plugin.Plugins{
				{
					Filename:     "filename",
					Path:         "path",
					BitBarURL:    fmt.Sprintf(github.ExportBaseBitBarURL, "path"),
					GitHubURL:    "url",
					GitHubRawURL: fmt.Sprintf(github.ExportBaseGitHubRawURL, "owner", "repo", "path"),
				},
				{
					Filename:     "filename",
					Path:         "path",
					BitBarURL:    fmt.Sprintf(github.ExportBaseBitBarURL, "path"),
					GitHubURL:    "url",
					GitHubRawURL: fmt.Sprintf(github.ExportBaseGitHubRawURL, "owner", "repo", "path"),
				},
			},
			wantErr: false,
		},
		{
			name: "github client error",
			codeFunc: func(ctx context.Context, query string, opt *gogithub.SearchOptions) (*gogithub.CodeSearchResult, *gogithub.Response, error) {
				return nil, nil, errors.New("error")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		s := new(github.ExportService)
		s.ExportSetGithubSearchService(&fakeGithubSearchService{
			code: tc.codeFunc,
		})

		t.Run(tc.name, func(t *testing.T) {
			got, err := s.SearchByFilename(context.Background(), "q")
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
