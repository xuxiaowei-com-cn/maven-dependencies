package types

import (
	"regexp"
	"strings"
)

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

func Dependencies(str string) (Dependency, []Dependency, error) {
	var parent Dependency
	var dependencies []Dependency

	lines := strings.Split(str, "\n")

	parentStart := false
	dependenciesStart := false
	dependencyStart := false
	add := false
	for index, line := range lines {

		if strings.HasPrefix(line, "<parent") {
			parentStart = true
			continue
		} else if strings.HasPrefix(line, "</parent") {
			parentStart = false
			continue
		}

		if parentStart {
			if strings.HasPrefix(line, "<groupId") {

				c, err := Context(line)
				if err != nil {
					return parent, nil, err
				}

				parent.GroupId = c
				parent.GroupIdLine = index + 1

			} else if strings.HasPrefix(line, "<artifactId") {
				c, err := Context(line)
				if err != nil {
					return parent, nil, err
				}

				parent.ArtifactId = c
				parent.ArtifactIdLine = index + 1

			} else if strings.HasPrefix(line, "<version") {
				c, err := Context(line)
				if err != nil {
					return parent, nil, err
				}

				parent.Version = c
				parent.VersionLine = index + 1
			}
		}

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
						return parent, nil, err
					}

					if add {
						dependency := dependencies[len(dependencies)-1]
						dependency.GroupId = c
						dependency.GroupIdLine = index + 1
					} else {
						dependencies = append(dependencies, Dependency{
							GroupId:     c,
							GroupIdLine: index + 1,
						})
						add = true
					}

				} else if strings.HasPrefix(line, "<artifactId") {
					c, err := Context(line)
					if err != nil {
						return parent, nil, err
					}

					if add {
						dependency := dependencies[len(dependencies)-1]
						dependency.ArtifactId = c
						dependency.ArtifactIdLine = index + 1
						dependencies[len(dependencies)-1] = dependency
					} else {
						dependencies = append(dependencies, Dependency{
							ArtifactId:     c,
							ArtifactIdLine: index + 1,
						})
						add = true
					}
				} else if strings.HasPrefix(line, "<version") {
					c, err := Context(line)
					if err != nil {
						return parent, nil, err
					}

					if add {
						dependency := dependencies[len(dependencies)-1]
						dependency.Version = c
						dependency.VersionLine = index + 1
						dependencies[len(dependencies)-1] = dependency
					} else {
						dependencies = append(dependencies, Dependency{
							Version:     c,
							VersionLine: index + 1,
						})
						add = true
					}
				}
			}
		}
	}

	return parent, dependencies, nil
}
