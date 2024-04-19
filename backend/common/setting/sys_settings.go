package setting

import (
	"path"

	"github.com/adrg/xdg"
)

func UserSettingsFile() string {
	p, _ := xdg.ConfigFile(path.Join("pvftools", "settings.json"))
	return p
}

func DBFile() string {
	p, _ := xdg.DataFile(path.Join("pvftools", "data.db"))
	return p
}
