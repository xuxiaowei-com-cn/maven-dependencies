package types

import (
	"regexp"
	"strings"
)

type Dependency struct {
	GroupId        string
	GroupIdLine    int
	ArtifactId     string
	ArtifactIdLine int
	Version        string
	VersionLine    int
}

func Context(str string) (string, error) {
	re, err := regexp.Compile(">[^<]*<")
	if err != nil {
		return "", err
	}

	matches := re.FindAllString(str, -1)
	if matches == nil {
		return "", nil
	}

	ml := len(matches)

	if ml == 0 {
		return "", nil
	}

	return strings.TrimRight(strings.TrimLeft(matches[0], ">"), "<"), nil
}

func ParseDependencies(str string) ([]Dependency, error) {
	var parsedDependencies []Dependency

	lines := strings.Split(str, "\n")

	dependenciesStart := false
	dependencyStart := false
	add := false
	for index, line := range lines {
		if strings.HasPrefix(line, "<dependencies") {
			dependenciesStart = true
			continue
		} else if strings.HasPrefix(line, "</dependencies") {
			dependenciesStart = false
			continue
		}

		if dependenciesStart {

			if strings.HasPrefix(line, "<dependency") {
				dependencyStart = true
				continue
			} else if strings.HasPrefix(line, "</dependency") {
				dependencyStart = false
				add = false
				continue
			}

			if dependencyStart {
				if strings.HasPrefix(line, "<groupId") {

					c, err := Context(line)
					if err != nil {
						return nil, err
					}

					if add {
						dependency := parsedDependencies[len(parsedDependencies)-1]
						dependency.GroupId = c
						dependency.GroupIdLine = index + 1
					} else {
						parsedDependencies = append(parsedDependencies, Dependency{
							GroupId:     c,
							GroupIdLine: index + 1,
						})
						add = true
					}

				} else if strings.HasPrefix(line, "<artifactId") {
					c, err := Context(line)
					if err != nil {
						return nil, err
					}

					if add {
						dependency := parsedDependencies[len(parsedDependencies)-1]
						dependency.ArtifactId = c
						dependency.ArtifactIdLine = index + 1
						parsedDependencies[len(parsedDependencies)-1] = dependency
					} else {
						parsedDependencies = append(parsedDependencies, Dependency{
							ArtifactId:     c,
							ArtifactIdLine: index + 1,
						})
						add = true
					}
				} else if strings.HasPrefix(line, "<version") {
					c, err := Context(line)
					if err != nil {
						return nil, err
					}

					if add {
						dependency := parsedDependencies[len(parsedDependencies)-1]
						dependency.Version = c
						dependency.VersionLine = index + 1
						parsedDependencies[len(parsedDependencies)-1] = dependency
					} else {
						parsedDependencies = append(parsedDependencies, Dependency{
							Version:     c,
							VersionLine: index + 1,
						})
						add = true
					}
				}
			}
		}
	}

	return parsedDependencies, nil
}
