/*
miningPoolCli – open-source tonuniverse mining pool client

Copyright (C) 2021 tonuniverse.com

This file is part of miningPoolCli.

miningPoolCli is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

miningPoolCli is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with miningPoolCli.  If not, see <https://www.gnu.org/licenses/>.
*/

package initp

import (
	"flag"
	"fmt"
	"miningPoolCli/config"
	"miningPoolCli/utils/api"
	"miningPoolCli/utils/getminer"
	"miningPoolCli/utils/gpuwrk"
	"miningPoolCli/utils/mlog"
	"os"
	"runtime"
)

func InitProgram() []gpuwrk.GPUstruct {
	config.Configure()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, config.Texts.GlobalHelpText)
	}

	flag.StringVar(&config.ServerSettings.AuthKey, "pool-id", "", "")
	flag.StringVar(&config.ServerSettings.MiningPoolServerURL, "url", "https://pool.tonuniverse.com", "")
	flag.BoolVar(&config.UpdateStatsFile, "stats", false, "") // for Hive OS

	flag.BoolVar(&config.NetSrv.RunThis, "serve-stat", false, "") // run http server with miner stat

	flag.Parse()

	switch "" {
	case config.ServerSettings.AuthKey:
		mlog.LogFatal("Flag -pool-id is required; for help run with -h flag")
	}

	mlog.LogText(config.Texts.Logo)
	mlog.LogText(config.Texts.WelcomeAdditionalMsg)

	os, architecture := runtime.GOOS, runtime.GOARCH

	if os == config.OSType.Win {
		mlog.LogOk("Supported OS detected: " + os + "/" + architecture)
	} else if os == config.OSType.Macos {
		mlog.LogFatal("Unsupported OS detected: " + "Mac OS")
	} else if os == config.OSType.Linux && architecture == "amd64" {
		mlog.LogOk("Supported OS detected: " + os + "/" + architecture)
	} else {
		mlog.LogFatal("Unsupported OS detected: " + os + "/" + architecture)
	}
	mlog.LogInfo("Using mining pool API url: " + config.ServerSettings.MiningPoolServerURL)
	config.OS.OperatingSystem, config.OS.Architecture = os, architecture

	api.Auth()

	getminer.UbubntuGetMiner()
	gpusArray := gpuwrk.SearchGpus()

	mlog.LogPass()
	gpuwrk.LogGpuList(gpusArray)
	mlog.LogInfo("Launching the mining processes...")

	return gpusArray
}
