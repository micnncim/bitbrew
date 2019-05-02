package di

import (
	"github.com/micnncim/bitbrew/bitbrew"
	"github.com/micnncim/bitbrew/github"
)

func InitBitBrew(token, formulaPath, pluginFolder string) (bitbrew.Service, error) {
	gh, err := github.NewService(token)
	if err != nil {
		return nil, err
	}
	return bitbrew.NewService(gh, formulaPath, pluginFolder), nil
}
