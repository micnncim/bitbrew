package bitbrew

import (
	"github.com/micnncim/bitbrew/plugin"
)

type (
	ExportService = service
)

func (s *ExportService) ExportPlugins() plugin.Plugins {
	return s.plugins
}

func (s *ExportService) ExportSetPlugins(plugins plugin.Plugins) {
	s.plugins = plugins
}

func (s *ExportService) ExportSetFormulaPath(formulaPath string) {
	s.formulaPath = formulaPath
}

func (s *ExportService) ExportSetPluginFolder(pluginFolder string) {
	s.pluginFolder = pluginFolder
}
