package types

type Project struct {
	Parent                 Parent
	GroupId                string
	GroupIdLine            int
	GroupIdLeft            int
	ArtifactId             string
	ArtifactIdLine         int
	ArtifactIdLeft         int
	Version                string
	VersionLine            int
	VersionLeft            int
	Description            string
	DescriptionLine        int
	DescriptionLeft        int
	DistributionManagement DistributionManagement
	DependencyManagement   DependencyManagement
	Dependencies           []Dependency
	Build                  Build
	Reporting              Reporting
}

type Parent struct {
	GroupId          string
	GroupIdLine      int
	GroupIdLeft      int
	ArtifactId       string
	ArtifactIdLine   int
	ArtifactIdLeft   int
	Version          string
	VersionLine      int
	VersionLeft      int
	RelativePath     string
	RelativePathLine int
	RelativePathLeft int
}

type DistributionManagement struct {
	Relocation Relocation
}

type Relocation struct {
	GroupId        string
	GroupIdLine    int
	GroupIdLeft    int
	ArtifactId     string
	ArtifactIdLine int
	ArtifactIdLeft int
	Version        string
	VersionLine    int
	VersionLeft    int
}

type DependencyManagement struct {
	Dependencies []Dependency
}

type Dependency struct {
	GroupId        string
	GroupIdLine    int
	GroupIdLeft    int
	ArtifactId     string
	ArtifactIdLine int
	ArtifactIdLeft int
	Version        string
	VersionLine    int
	VersionLeft    int
	Exclusions     []Exclusion
}

type Exclusion struct {
	GroupId        string
	GroupIdLine    int
	GroupIdLeft    int
	ArtifactId     string
	ArtifactIdLine int
	ArtifactIdLeft int
}

type Build struct {
	Extensions       []Extension
	PluginManagement PluginManagement
	Plugins          []Plugin
	Profiles         []Profile
}

type Extension struct {
	GroupId        string
	GroupIdLine    int
	GroupIdLeft    int
	ArtifactId     string
	ArtifactIdLine int
	ArtifactIdLeft int
	Version        string
	VersionLine    int
	VersionLeft    int
}

type PluginManagement struct {
	Plugins []Plugin
}

type Plugin struct {
	GroupId        string
	GroupIdLine    int
	GroupIdLeft    int
	ArtifactId     string
	ArtifactIdLine int
	ArtifactIdLeft int
	Version        string
	VersionLine    int
	VersionLeft    int
	Dependencies   []Dependency
}

type Reporting struct {
	Plugins []Plugin
}

type Profile struct {
	Build                  Build
	DistributionManagement DistributionManagement
	DependencyManagement   DependencyManagement
	Dependencies           []Dependency
	Reporting              Reporting
}
