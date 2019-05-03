package bitbrew

import (
	"github.com/micnncim/bitbrew/plugin"
)

type (
	ExportBitbrew = bitbrew
)

var (
	ExportBitbrewAddFormula    = (*bitbrew).addFormula
	ExportBitbrewRemoveFormula = (*bitbrew).removeFormula
)

func (b *ExportBitbrew) ExportPlugins() plugin.Plugins {
	return b.plugins
}

func (b *ExportBitbrew) ExportSetPlugins(plugins plugin.Plugins) {
	b.plugins = plugins
}

func (b *ExportBitbrew) ExportSetFormulaPath(formulaPath string) {
	b.formulaPath = formulaPath
}

func (b *ExportBitbrew) ExportSetPluginFolder(pluginFolder string) {
	b.pluginFolder = pluginFolder
}
