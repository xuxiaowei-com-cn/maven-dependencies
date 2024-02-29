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

func Dependencies(str string) (Dependency, Dependency, []Dependency, error) {
	var current Dependency
	var parent Dependency
	var dependencies []Dependency

	lines := strings.Split(str, "\n")

	parentStart := false
	distributionManagementStart := false
	dependencyManagementStart := false
	dependenciesStart := false
	dependenciesExclusionStart := false
	buildStart := false
	reportingStart := false
	profileStart := false

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

		if strings.HasPrefix(line, "<distributionManagement") {
			distributionManagementStart = true
			continue
		} else if strings.HasPrefix(line, "</distributionManagement") {
			distributionManagementStart = false
			continue
		}

		if strings.HasPrefix(line, "<dependencyManagement") {
			dependencyManagementStart = true
			continue
		} else if strings.HasPrefix(line, "</dependencyManagement") {
			dependencyManagementStart = false
			continue
		}

		if strings.HasPrefix(line, "<build") {
			buildStart = true
			continue
		} else if strings.HasPrefix(line, "</build") {
			buildStart = false
			continue
		}

		if strings.HasPrefix(line, "<reporting") {
			reportingStart = true
			continue
		} else if strings.HasPrefix(line, "</reporting") {
			reportingStart = false
			continue
		}

		if strings.HasPrefix(line, "<profile") {
			profileStart = true
			continue
		} else if strings.HasPrefix(line, "</profile") {
			profileStart = false
			continue
		}

		if !parentStart && !distributionManagementStart && !dependencyManagementStart && !buildStart && !reportingStart && !profileStart && !dependenciesStart {
			// 不是 <parent> 标签内的坐标
			// 不是 <distributionManagement> 标签内的坐标
			// 不是 <dependencyManagement> 标签内的坐标
			// 不是 <build> 标签内的坐标
			// 不是 <reporting> 标签内的坐标
			// 不是 <profile> 标签内的坐标
			if strings.HasPrefix(line, "<groupId") {

				c, err := Context(line)
				if err != nil {
					return current, parent, nil, err
				}

				current.GroupId = c
				current.GroupIdLine = index + 1

			} else if strings.HasPrefix(line, "<artifactId") {
				c, err := Context(line)
				if err != nil {
					return current, parent, nil, err
				}

				current.ArtifactId = c
				current.ArtifactIdLine = index + 1

			} else if strings.HasPrefix(line, "<version") {
				c, err := Context(line)
				if err != nil {
					return current, parent, nil, err
				}

				current.Version = c
				current.VersionLine = index + 1
			}
		}

		if parentStart {
			if strings.HasPrefix(line, "<groupId") {

				c, err := Context(line)
				if err != nil {
					return current, parent, nil, err
				}

				parent.GroupId = c
				parent.GroupIdLine = index + 1

			} else if strings.HasPrefix(line, "<artifactId") {
				c, err := Context(line)
				if err != nil {
					return current, parent, nil, err
				}

				parent.ArtifactId = c
				parent.ArtifactIdLine = index + 1

			} else if strings.HasPrefix(line, "<version") {
				c, err := Context(line)
				if err != nil {
					return current, parent, nil, err
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

			if strings.HasPrefix(line, "<exclusion") {
				dependenciesExclusionStart = true
				continue
			} else if strings.HasPrefix(line, "</exclusion") {
				dependenciesExclusionStart = false
				add = false
				continue
			}

			if dependenciesExclusionStart {
				continue
			}

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
						return current, parent, nil, err
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
						return current, parent, nil, err
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
						return current, parent, nil, err
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

	return current, parent, dependencies, nil
}
