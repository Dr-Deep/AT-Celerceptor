package main

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/fatih/color"
)

const (
	bannerTop = ` ################################# 
##########..###.###.##.##.#########
##########.-###.###.##-##.#########
##########.#.##.###.##=##.#########
#########....##.###.##.##.#########
#########.##.##...#....##.#########
###################################`

	bannerBottom = `-----------------------------------
------####---##----##---##-##------
-------##---####---##---##-#-------
-------##---####---##---####-------
-------##--##--##--##---##-#-------
-------##--##--##--####-##-##------
 --------------------------------- `
)

func banner() {
	var (
		atBlue   = color.BgRGB(27, 50, 130) // Original Hex colors from playstore image
		atOrange = color.BgRGB(239, 124, 28)
	)

	atBlue.Println(bannerTop)
	atOrange.Println(bannerBottom)
}

func flagVersion() {
	var (
		gitRevision = "dbg-dist"
		buildTime   = "unknown"
	)

	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Printf(
			"by <github.com/dr-Deep>\n%s\n",
			runtime.Version(),
		)

		return
	}

	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			gitRevision = setting.Value[:7]
		}

		if setting.Key == "vcs.time" {
			buildTime = setting.Value
		}
	}

	//
	banner()
	fmt.Printf("Git: '%s'\n", gitRevision)
	fmt.Printf("Build-Time: '%s'\n", buildTime)
	fmt.Printf("==> by <github.com/dr-Deep> <==\n")
}
