package command

import (
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/maven-dependencies/constant"
	"github.com/xuxiaowei-com-cn/maven-dependencies/file"
	"github.com/xuxiaowei-com-cn/maven-dependencies/flag"
	"github.com/xuxiaowei-com-cn/maven-dependencies/types"
	"log"
	"os"
)

func DependencyCommand() *cli.Command {
	return &cli.Command{
		Name:  "dependency",
		Usage: "Maven 坐标",
		Flags: []cli.Flag{
			flag.FilePathFlag(false),
			flag.GroupIdFlag(false), flag.ArtifactIdFlag(false), flag.VersionFlag(false),
		},
		Subcommands: []*cli.Command{
			EditDependencyCommand(),
		},
	}
}

func EditDependencyCommand() *cli.Command {
	return &cli.Command{
		Name:  "edit",
		Usage: "修改 Maven 坐标",
		Flags: []cli.Flag{
			flag.FilePathFlag(true),
			flag.GroupIdFlag(true), flag.ArtifactIdFlag(true), flag.VersionFlag(true),
			flag.AfterGroupIdFlag(), flag.AfterArtifactIdFlag(), flag.AfterVersionFlag(),
		},
		Action: func(context *cli.Context) error {
			var filePath = context.Path(constant.FilePath)
			var groupId = context.String(constant.GroupId)
			var artifactId = context.String(constant.ArtifactId)
			var version = context.String(constant.Version)
			var afterGroupId = context.String(constant.AfterGroupId)
			var afterArtifactId = context.String(constant.AfterArtifactId)
			var afterVersion = context.String(constant.AfterVersion)

			fileContext, err := file.ReadFileTrimSpace(filePath)
			if err != nil {
				return err
			}

			fileContextLines, err := file.ReadFileLines(filePath)
			if err != nil {
				return err
			}

			parsedDependencies, err := types.ParseDependencies(fileContext)
			if err != nil {
				return err
			}

			var dependency types.Dependency
			for _, parsedDependency := range parsedDependencies {
				if parsedDependency.GroupId == groupId && parsedDependency.ArtifactId == artifactId && parsedDependency.Version == version {
					dependency.GroupId = parsedDependency.GroupId
					dependency.GroupIdLine = parsedDependency.GroupIdLine
					dependency.ArtifactId = parsedDependency.ArtifactId
					dependency.ArtifactIdLine = parsedDependency.ArtifactIdLine
					dependency.Version = parsedDependency.Version
					dependency.VersionLine = parsedDependency.VersionLine
					break
				}
			}

			if dependency.GroupId != "" {
				var groupIdLine = dependency.GroupIdLine
				var artifactIdLine = dependency.ArtifactIdLine
				var versionLine = dependency.VersionLine

				var result string
				for index, line := range fileContextLines {
					if index != 0 {
						result += "\n"
					}
					if groupIdLine == index+1 && afterGroupId != "" {
						result += "<groupId>" + afterGroupId + "</groupId>"
					} else if artifactIdLine == index+1 && afterArtifactId != "" {
						result += "<artifactId>" + afterArtifactId + "</artifactId>"
					} else if versionLine == index+1 && afterVersion != "" {
						result += "<version>" + afterVersion + "</version>"
					} else {
						result += line
					}
				}

				err = os.WriteFile(filePath, []byte(result), 0644)
				if err != nil {
					return err
				}

			} else {
				log.Printf("Dependency not found")
			}

			return nil
		},
	}
}
