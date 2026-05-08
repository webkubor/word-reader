// package main — 每个 Go 程序的入口包，必须有且只有一个 main 包
package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// go:embed 是 Go 的编译指令，在编译时把 frontend/dist 目录的所有文件
// 嵌入到二进制文件中，这样桌面应用不需要额外的前端文件
// embed.FS 是一个文件系统类型，类似前端的 public 目录被打包进 exe/app
//
//go:embed all:frontend/dist
var assets embed.FS

// main — 程序入口，类似前端的 index.js
func main() {
	// NewApp() 创建应用实例，定义在 app.go 中
	app := NewApp()

	// wails.Run 启动桌面应用，传入配置项
	err := wails.Run(&options.App{
		// 窗口标题
		Title: "Word Reader",
		// 窗口初始宽高（像素）
		Width:  520,
		Height: 680,
		// 窗口最小宽高，用户不能再缩小
		MinWidth:  420,
		MinHeight: 500,
		// AssetServer 配置前端静态资源
		// Wails 内置了一个 HTTP 服务器来服务前端文件
		// 类似前端的 devServer，但这里是嵌入在桌面应用里的
		AssetServer: &assetserver.Options{
			Assets: assets, // 上面 embed 嵌入的前端文件
		},
		// 窗口背景色（深色主题）
		BackgroundColour: &options.RGBA{R: 22, G: 19, B: 17, A: 1},
		// OnStartup — 应用启动时的回调函数，类似前端的 onMounted
		// app.startup 定义在 app.go 中，负责初始化数据库
		OnStartup: app.startup,
		// Bind — 把 Go 对象绑定到前端 JS 可以调用
		// 绑定后，前端可以通过 window.go.main.App.xxx() 调用 App 的方法
		// 类似前端的 expose API 给渲染进程
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
