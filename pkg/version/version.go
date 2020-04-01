package version

import "fmt"

const (
	VersionMajor  = 0
	VersionMinor  = 1
	VersionBugfix = 0
)

func String() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionBugfix)
}
