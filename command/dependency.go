package command

import (
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/maven-dependencies/constant"
	"github.com/xuxiaowei-com-cn/maven-dependencies/file"
	"github.com/xuxiaowei-com-cn/maven-dependencies/flag"
	"github.com/xuxiaowei-com-cn/maven-dependencies/types"
	"log"
	"os"
	"strings"
)

func DependencyCommand() *cli.Command {
	return &cli.Command{
		Name:  "dependency",
		Usage: "Maven 坐标",
		Flags: []cli.Flag{
			flag.FilePathFlag(false),
			flag.GroupIdFlag(false), flag.ArtifactIdFlag(false), flag.VersionFlag(),
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
			flag.GroupIdFlag(true), flag.ArtifactIdFlag(true), flag.VersionFlag(),
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

			if afterGroupId == "" && afterArtifactId == "" && afterVersion == "" {
				log.Printf("未接收到修改结果，取消本次任务")
				return nil
			}

			fileContext, err := file.ReadFileTrimSpace(filePath)
			if err != nil {
				return err
			}

			fileContextLines, err := file.ReadFileLines(filePath)
			if err != nil {
				return err
			}

			_, _, dependencies, err := types.Dependencies(fileContext)
			if err != nil {
				return err
			}

			var dependencyResult types.Dependency
			for _, dependency := range dependencies {
				if dependency.GroupId == groupId && dependency.ArtifactId == artifactId && dependency.Version == version {
					dependencyResult.GroupId = dependency.GroupId
					dependencyResult.GroupIdLine = dependency.GroupIdLine
					dependencyResult.ArtifactId = dependency.ArtifactId
					dependencyResult.ArtifactIdLine = dependency.ArtifactIdLine
					dependencyResult.Version = dependency.Version
					dependencyResult.VersionLine = dependency.VersionLine
					break
				}
			}

			if dependencyResult.GroupId != "" {
				var groupIdLine = dependencyResult.GroupIdLine
				var artifactIdLine = dependencyResult.ArtifactIdLine
				var versionLine = dependencyResult.VersionLine

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

				if dependencyResult.Version == "" && afterVersion != "" {
					lines := strings.Split(result, "\n")
					result = ""

					for index, line := range lines {
						if index != 0 {
							result += "\n"
						}
						result += line
						if groupIdLine == index+1 && afterGroupId != "" {
							result += "\n<version>" + afterVersion + "</version>"
						}
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
