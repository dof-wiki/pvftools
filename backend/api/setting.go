package api

import "pvftools/backend/common/setting"

func (a *App) GetUserSettings() *setting.SettingStruct {
	return setting.UserSettings
}

func (a *App) SetUserSettings(s *setting.SettingStruct) error {
	return setting.SetUserSettings(s)
}
