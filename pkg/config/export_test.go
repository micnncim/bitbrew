package config

import (
	"github.com/spf13/viper"
)

var (
	ExportInitConfigFunc = initConfigFunc
)

func ExportSetDefaultConfigName(s string) (resetFunc func()) {
	var org string
	org, defaultConfigName = defaultConfigName, s
	return func() {
		defaultConfigName = org
	}
}

func ExportSetInitConfigFunc(f func(string) (*Config, error)) (resetFunc func()) {
	var org func(string) (*Config, error)
	org, initConfigFunc = initConfigFunc, f
	return func() {
		initConfigFunc = org
	}
}

func ExportSetNewDefaultViperFunc(f func(string) (*viper.Viper, error)) (resetFunc func()) {
	var org func(string) (*viper.Viper, error)
	org, newDefaultViperFunc = newDefaultViperFunc, f
	return func() {
		newDefaultViperFunc = org
	}
}
