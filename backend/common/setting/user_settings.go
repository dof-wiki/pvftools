package setting

import (
	"encoding/json"
	"errors"
	"os"
	"pvftools/backend/common/log"
)

const (
	ModePvfUtility = "pvfutility"
	ModeFile       = "file"
)

type SettingStruct struct {
	Mode            string `json:"mode"`
	Target          string `json:"target"`
	ConcurrentCount int    `json:"concurrent_count"` // 并发数
}

var UserSettings *SettingStruct

func defaultUserSettings() *SettingStruct {
	return &SettingStruct{
		Mode:            ModePvfUtility,
		Target:          "http://localhost:27000",
		ConcurrentCount: 10,
	}
}

func loadUserSettings() {
	UserSettings = new(SettingStruct)
	buf, err := os.ReadFile(UserSettingsFile())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			UserSettings = defaultUserSettings()
			log.LogInfo("load user settings, use default: %v", UserSettings)
			return
		} else {
			log.LogPanic("load user settings err %v", err)
		}
	}
	err = json.Unmarshal(buf, UserSettings)
	if err != nil {
		log.LogPanic("load user settings decode err %v", err)
	}
	log.LogInfo("load user settings %v", UserSettings)
}

func SetUserSettings(s *SettingStruct) error {
	buf, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = os.WriteFile(UserSettingsFile(), buf, 0644)
	if err != nil {
		return err
	}
	UserSettings = s
	return nil
}
