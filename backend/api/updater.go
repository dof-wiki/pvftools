package api

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"pvftools/backend/common/consts"
	"pvftools/backend/common/log"
)

func (a *App) downloadNewVersion() {
	a.updater.AutoUpdate()
	runtime.EventsEmit(a.ctx, consts.EventSelfUpdateReady)
	result, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "自动更新",
		Message:       "新版本已下载完成, 是否立即执行更新?",
		DefaultButton: "Ok",
		Buttons:       []string{"Ok", "No"},
	})
	log.LogInfo("Result: %v", result)
	if result == "Ok" {
	}
}

func (a *App) CheckUpdate() bool {
	return a.checkUpdate()
}

func (a *App) DownloadNewVersion() {
	go a.downloadNewVersion()
}

func (a *App) DoUpdate() {
	a.updater.DoUpdate()
}
