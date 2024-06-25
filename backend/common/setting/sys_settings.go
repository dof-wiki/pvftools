package setting

import (
	"github.com/labstack/gommon/log"
	"os"
	"path"

	"github.com/adrg/xdg"
)

func UserSettingsFile() string {
	p, _ := xdg.ConfigFile(path.Join("pvftools", "settings.json"))
	return p
}

func DBFile() string {
	p, _ := xdg.DataFile(path.Join("pvftools", "data.db"))
	log.Debug("dbfile:", p)
	return p
}

func UpdaterDir() string {
	p, _ := xdg.DataFile(path.Join("pvftools", "updater"))
	log.Debug("updater dir:", p)
	os.MkdirAll(p, 0755)
	return p
}
