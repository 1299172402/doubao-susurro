package assets

import (
	_ "embed"
)

//go:embed static/index.html
var IndexPage []byte

//go:embed static/logo.png
var LogoPNG []byte

//go:embed static/startup/windows.vbs
var StartupWindowsVBS string

//go:embed static/startup/macos.plist
var StartupMacOSPlist string

//go:embed static/startup/linux.desktop
var StartupLinuxDesktop string
