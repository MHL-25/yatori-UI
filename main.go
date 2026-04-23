package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	appMenu := createMenu(app)

	err := wails.Run(&options.App{
		Title:     "Yatori-UI - 智能网课助手",
		Width:     1280,
		Height:    800,
		MinWidth:  1024,
		MinHeight: 680,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour:   &options.RGBA{R: 15, G: 15, B: 26, A: 255},
		OnStartup:          app.startup,
		OnBeforeClose:      app.beforeClose,
		HideWindowOnClose:  true,
		Menu:               appMenu,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			WebviewBrowserPath:   "",
		},
		Frameless:   true,
		StartHidden: false,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func createMenu(app *App) *menu.Menu {
	appMenu := menu.NewMenu()
	fileMenu := appMenu.AddSubmenu("文件")
	fileMenu.AddText("刷新", keys.CmdOrCtrl("r"), func(_ *menu.CallbackData) {})
	fileMenu.AddSeparator()
	fileMenu.AddText("导入配置", keys.CmdOrCtrl("o"), func(cd *menu.CallbackData) {
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("退出程序", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		app.shouldQuit = true
		runtime.Quit(app.ctx)
	})

	helpMenu := appMenu.AddSubmenu("帮助")
	helpMenu.AddText("关于", nil, func(_ *menu.CallbackData) {})

	return appMenu
}
