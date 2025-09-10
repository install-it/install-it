package main

import (
	"context"
	"embed"
	"install-it/pkg/execute"
	"install-it/pkg/porter"
	"install-it/pkg/status"
	"install-it/pkg/storage"
	"install-it/pkg/sysinfo"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

var (
	dirRoot string
	// Path to the configuration directory
	dirConf string
	// Path to the driver directory
	dirDir string
	// Path to the WebView2 executable
	pathWV2 string
	// Type of the build version number
	buildVersion string
	// Version struct, parsed from [buildVersion]
	version *semver.Version
)

func init() {
	if buildVersion == "" {
		version, _ = semver.NewVersion("0.0.0")
	} else if v, err := semver.NewVersion(buildVersion); err != nil {
		panic(err)
	} else {
		version = v
	}

	if pathExe, err := os.Executable(); err != nil {
		panic(err)
	} else {
		dirRoot = filepath.Dir(pathExe)

		dirConf = filepath.Join(dirRoot, "conf")
		if _, err := os.Stat(dirConf); err != nil {
			if err := os.MkdirAll(dirConf, os.ModePerm); err != nil {
				panic(err)
			}
		}

		dirDir = filepath.Join(dirRoot, "drivers")
		if _, err := os.Stat(dirDir); err != nil {
			if err := os.MkdirAll(dirDir, os.ModePerm); err != nil {
				panic(err)
			}
		}

		for _, name := range [3]string{"network", "display", "miscellaneous"} {
			os.MkdirAll(filepath.Join(dirDir, name), os.ModePerm)
		}

		// WebView2 binary lookup
		pathWV2 = filepath.Join(dirRoot, "bin", "WebView2")
		if _, err := os.Stat(pathWV2); err != nil {
			pathWV2 = ""
		}
	}
}

func main() {
	app := &App{}
	mgt := &execute.CommandExecutor{}

	err := wails.Run(&options.App{
		Title:     "install-it",
		Width:     768,
		Height:    576,
		MinWidth:  640,
		MinHeight: 480,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			// working directory correction
			if cwd, err := os.Getwd(); err == nil {
				if pathExe, err := os.Executable(); err == nil && cwd != filepath.Dir(pathExe) {
					os.Chdir(filepath.Dir(pathExe))
				}
			}

			app.SetContext(ctx)
			mgt.SetContext(ctx)
		},
		Bind: []interface{}{
			app,
			mgt,
			&storage.AppSettingStorage{Store: &storage.FileStore{Path: filepath.Join(dirConf, "setting.json")}},
			&storage.DriverGroupStorage{Store: &storage.FileStore{Path: filepath.Join(dirConf, "groups.json")}},
			&porter.Porter{DirRoot: dirRoot, Message: make(chan string, 512), Targets: []string{dirConf, dirDir}},
			&sysinfo.SysInfo{},
		},
		EnumBind: []interface{}{
			[]struct {
				Value  storage.DriverType
				TSName string
			}{
				{storage.Network, "NETWORK"},
				{storage.Display, "DISPLAY"},
				{storage.Miscellaneous, "MISCELLANEOUS"},
			},
			[]struct {
				Value  storage.SuccessAction
				TSName string
			}{
				{storage.Nothing, "NOTHING"},
				{storage.Reboot, "REBOOT"},
				{storage.Shutdown, "SHUTDOWN"},
				{storage.Firmware, "FIRMWARE"},
			},
			[]struct {
				Value  status.Status
				TSName string
			}{
				{status.Pending, "PENDING"},
				{status.Running, "RUNNING"},
				{status.Completed, "COMPLETED"},
				{status.Failed, "FAILED"},
				{status.Aborting, "ABORTING"},
				{status.Aborted, "ABORTED"},
				{status.Skiped, "SKIPED"},
				{status.Speeded, "SPEEDED"},
				{status.Errored, "ERRORED"},
			},
		},
		Windows: &windows.Options{
			WebviewBrowserPath: pathWV2,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
