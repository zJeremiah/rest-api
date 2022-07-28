package version

import (
	"fmt"
	"runtime"
)

var (
	// These flags are set at build time with `-ldflags "-X path.to.package.Version x.x.x"` etc.
	// in the Makefile.

	// Version specifies the git hash for this build.
	Version = "-"
	// BuildTimeUTC is specifies the build time.
	BuildTimeUTC = "-"
	// AppName is defined at build time.
	AppName = "-"
)

type Struct struct {
	Build     string `json:"build"`
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
}

type Response struct {
	App    string `json:"app_name"`
	Struct `json:"build_info"`
}

// Get returns version data as a string.
func Get() string {
	return fmt.Sprintf(
		"app: %s\nversion: %s (built w/%s)\nUTC Build Time: %s",
		AppName,
		Version,
		runtime.Version(),
		BuildTimeUTC,
	)
}

func JSON() Response {
	return Response{
		App: AppName,
		Struct: Struct{
			runtime.Version(),
			Version,
			BuildTimeUTC,
		},
	}
}
