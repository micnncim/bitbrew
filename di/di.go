package di

import (
	"errors"

	"github.com/micnncim/bitbrew/bitbrew"
	"github.com/micnncim/bitbrew/github"
)

func InitBitBrew(token, formulaPath, pluginFolder string) (bitbrew.Bitbrew, error) {
	if token == "" {
		return nil, errors.New("github token is missing")
	}
	gh, err := github.NewService(token)
	if err != nil {
		return nil, err
	}
	return bitbrew.New(gh, formulaPath, pluginFolder), nil
}
