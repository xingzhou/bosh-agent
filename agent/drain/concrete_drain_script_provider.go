package drain

import (
	"path/filepath"

	boshdirs "github.com/cloudfoundry/bosh-agent/settings/directories"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
)

type ConcreteDrainScriptProvider struct {
	cmdRunner   boshsys.CmdRunner
	fs          boshsys.FileSystem
	dirProvider boshdirs.DirectoriesProvider
}

func NewConcreteDrainScriptProvider(
	cmdRunner boshsys.CmdRunner,
	fs boshsys.FileSystem,
	dirProvider boshdirs.DirectoriesProvider,
) (provider ConcreteDrainScriptProvider) {
	provider.cmdRunner = cmdRunner
	provider.fs = fs
	provider.dirProvider = dirProvider
	return
}

func (p ConcreteDrainScriptProvider) NewDrainScript(templateName string) DrainScript {
	scriptPath := filepath.Join(p.dirProvider.JobsDir(), templateName, "bin", "drain")
	return NewConcreteDrainScript(p.fs, p.cmdRunner, scriptPath)
}
