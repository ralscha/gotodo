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
		revision := ""
		buildTime := ""
		for _, setting := range buildInfo.Settings {
			if setting.Key == "vcs.revision" {
				revision = shortRevision(setting.Value)
			} else if setting.Key == "vcs.time" {
				buildTime = setting.Value
			}
			if revision != "" && buildTime != "" {
				break
			}
		}

		return Version{
			BuildTime: buildTime,
			Version:   revision,
		}
	} else {
		return Version{
			BuildTime: "",
			Version:   "",
		}
	}
}

func shortRevision(revision string) string {
	if len(revision) <= 8 {
		return revision
	}
	return revision[:8]
}
