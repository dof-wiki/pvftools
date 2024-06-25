package main

import (
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"pvftools/backend/api"
	"pvftools/backend/common/ctx"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := api.NewApp()

	defer func() {
		if err := recover(); err != nil {
			runtime.MessageDialog(*ctx.Ctx, runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Title:   "发生了严重错误",
				Message: fmt.Sprintf("%v", err),
			})
		}
	}()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "PvfTools",
		Width:  1440,
		Height: 960,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
