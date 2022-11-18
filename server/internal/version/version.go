package version

import (
	"runtime/debug"
)

type Version struct {
	BuildTime string
	Version   string
}

func Get() Version {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		revison := ""
		time := ""
		for _, setting := range buildInfo.Settings {
			if setting.Key == "vcs.revision" {
				revison = setting.Value
				revison = revison[:8]
			} else if setting.Key == "vcs.time" {
				time = setting.Value
			}
			if revison != "" && time != "" {
				break
			}
		}

		return Version{
			BuildTime: time,
			Version:   revison,
		}
	} else {
		return Version{
			BuildTime: "",
			Version:   "",
		}
	}
}
