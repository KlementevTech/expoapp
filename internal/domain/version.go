package domain

type Version struct {
	version string
}

func NewVersion(version string) *Version {
	return &Version{
		version: version,
	}
}

func (v *Version) Version() string {
	return v.version
}
