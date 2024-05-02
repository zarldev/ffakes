package info

import (
	_ "embed"
	"strings"
)

//go:generate sh -c "git describe --tags --abbrev=0 > version.tag"
//go:generate sh -c "date -u +'%Y-%m-%d %H:%M:%S' > build.date"
//go:generate sh -c "git rev-parse HEAD > git.sha"
var (

	//go:embed version.tag
	version string

	//go:embed build.date
	buildDate string

	//go:embed git.sha
	gitSha string

	AppInfo = NewAppInfo()
)

// Info contains the version, build date and git sha of the application
type info struct {
	version   string
	buildDate string
	gitSHA    string
}

func NewAppInfo() info {
	v := strings.TrimSpace(version)
	b := strings.TrimSpace(buildDate)
	g := strings.TrimSpace(gitSha)

	return info{
		version:   v,
		buildDate: b,
		gitSHA:    g,
	}
}
func (i info) String() string {
	return i.version + " " + i.buildDate + " " + i.gitSHA
}

func (i info) Short() string {
	return i.version
}

func (i info) Long() string {
	return i.String()
}

func (i info) Version() string {
	return i.version
}

func (i info) BuildDate() string {
	return i.buildDate
}

func (i info) GitSHA() string {
	return i.gitSHA
}
