package cmd

import (
	"github.com/micnncim/bitbrew/bitbrew"
)

var (
	ExportBrowseFunc    = browseFunc
	ExportInstallFunc   = installFunc
	ExportUninstallFunc = uninstallFunc
	ExportSearchFunc    = searchFunc
	ExportSyncFunc      = syncFunc
	ExportListFunc      = listFunc
)

func ExportSetInitBitbrewClient(f func(token, formulaPath, pluginFolder string) (bitbrew.Client, error)) (resetFunc func()) {
	var org func(token, formulaPath, pluginFolder string) (bitbrew.Client, error)
	org, initBitbrewClient = initBitbrewClient, f
	return func() {
		initBitbrewClient = org
	}
}
