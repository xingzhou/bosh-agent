package infrastructure

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshplatform "github.com/cloudfoundry/bosh-agent/platform"
	boshsettings "github.com/cloudfoundry/bosh-agent/settings"
)

type genericInfrastructure struct {
	platform       boshplatform.Platform
	settingsSource SettingsSource

	networkingType          string
	staticEphemeralDiskPath string

	logTag string
	logger boshlog.Logger
}

const (
	NetworkingTypeDHCP   = "dhcp"
	NetworkingTypeManual = "manual"
)

func NewGenericInfrastructure(
	platform boshplatform.Platform,
	settingsSource SettingsSource,
	networkingType string,
	staticEphemeralDiskPath string,
	logger boshlog.Logger,
) genericInfrastructure {
	return genericInfrastructure{
		platform:       platform,
		settingsSource: settingsSource,

		networkingType:          networkingType,
		staticEphemeralDiskPath: staticEphemeralDiskPath,

		logTag: "genericInfrastructure",
		logger: logger,
	}
}

func (inf genericInfrastructure) SetupSSH(username string) error {
	publicKey, err := inf.settingsSource.PublicSSHKeyForUsername(username)
	if err != nil {
		return bosherr.WrapError(err, "Getting public key")
	}

	if len(publicKey) > 0 {
		return inf.platform.SetupSSH(publicKey, username)
	}

	return nil
}

func (inf genericInfrastructure) GetSettings() (boshsettings.Settings, error) {
	return inf.settingsSource.Settings()
}

// Existing examples:
// - vSphere: manual
// - AWS, Openstack: dhcp
// - Warden, Dummy: empty
func (inf genericInfrastructure) SetupNetworking(networks boshsettings.Networks) error {
	switch {
	case inf.networkingType == NetworkingTypeDHCP:
		return inf.platform.SetupDhcp(networks)

	case inf.networkingType == NetworkingTypeManual:
		return inf.platform.SetupManualNetworking(networks)

	default:
		return nil
	}
}

// Existing examples:
// - vSphere: static configuration "/dev/sdb"
// - AWS, Openstack: allows empty device path
// - AWS, OpenStack, Warden, Dummy: allows normalization
func (inf genericInfrastructure) GetEphemeralDiskPath(diskSettings boshsettings.DiskSettings) string {
	if len(diskSettings.Path) == 0 {
		return ""
	}

	if len(inf.staticEphemeralDiskPath) > 0 {
		return inf.staticEphemeralDiskPath
	}

	return inf.platform.NormalizeDiskPath(diskSettings)
}
