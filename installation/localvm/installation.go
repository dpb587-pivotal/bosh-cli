package localvm

import (
	biinstall "github.com/cloudfoundry/bosh-cli/installation"
	biui "github.com/cloudfoundry/bosh-cli/ui"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type Installation struct{}

var _ biinstall.Installation = &Installation{}

func (i Installation) Target() biinstall.Target {
	return biinstall.Target{}
}

func (i Installation) Job() biinstall.InstalledJob {
	return biinstall.InstalledJob{}
}

func (i Installation) WithRunningRegistry(boshlog.Logger, biui.Stage, func() error) error {
	return nil
}

func (i Installation) StartRegistry() error {
	return nil
}

func (i Installation) StopRegistry() error {
	return nil
}
