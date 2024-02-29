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

func Dependencies(str string) (Project, error) {
	var project Project

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
					return project, err
				}

				project.GroupId = c
				project.GroupIdLine = index + 1

			} else if strings.HasPrefix(line, "<artifactId") {
				c, err := Context(line)
				if err != nil {
					return project, err
				}

				project.ArtifactId = c
				project.ArtifactIdLine = index + 1

			} else if strings.HasPrefix(line, "<version") {
				c, err := Context(line)
				if err != nil {
					return project, err
				}

				project.Version = c
				project.VersionLine = index + 1
			} else if strings.HasPrefix(line, "<description") {
				c, err := Context(line)
				if err != nil {
					return project, err
				}

				project.Description = c
				project.DescriptionLine = index + 1
			}
		}

		if parentStart {
			if strings.HasPrefix(line, "<groupId") {

				c, err := Context(line)
				if err != nil {
					return project, err
				}

				project.Parent.GroupId = c
				project.Parent.GroupIdLine = index + 1

			} else if strings.HasPrefix(line, "<artifactId") {
				c, err := Context(line)
				if err != nil {
					return project, err
				}

				project.Parent.ArtifactId = c
				project.Parent.ArtifactIdLine = index + 1

			} else if strings.HasPrefix(line, "<version") {
				c, err := Context(line)
				if err != nil {
					return project, err
				}

				project.Parent.Version = c
				project.Parent.VersionLine = index + 1
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
						return project, err
					}

					if add {
						dependency := project.Dependencies[len(project.Dependencies)-1]
						dependency.GroupId = c
						dependency.GroupIdLine = index + 1
					} else {
						project.Dependencies = append(project.Dependencies, Dependency{
							GroupId:     c,
							GroupIdLine: index + 1,
						})
						add = true
					}

				} else if strings.HasPrefix(line, "<artifactId") {
					c, err := Context(line)
					if err != nil {
						return project, err
					}

					if add {
						dependency := project.Dependencies[len(project.Dependencies)-1]
						dependency.ArtifactId = c
						dependency.ArtifactIdLine = index + 1
						project.Dependencies[len(project.Dependencies)-1] = dependency
					} else {
						project.Dependencies = append(project.Dependencies, Dependency{
							ArtifactId:     c,
							ArtifactIdLine: index + 1,
						})
						add = true
					}
				} else if strings.HasPrefix(line, "<version") {
					c, err := Context(line)
					if err != nil {
						return project, err
					}

					if add {
						dependency := project.Dependencies[len(project.Dependencies)-1]
						dependency.Version = c
						dependency.VersionLine = index + 1
						project.Dependencies[len(project.Dependencies)-1] = dependency
					} else {
						project.Dependencies = append(project.Dependencies, Dependency{
							Version:     c,
							VersionLine: index + 1,
						})
						add = true
					}
				}
			}
		}
	}

	return project, nil
}
