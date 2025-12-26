package command

import (
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/maven-dependencies/constant"
	"github.com/xuxiaowei-com-cn/maven-dependencies/file"
	"github.com/xuxiaowei-com-cn/maven-dependencies/flag"
	"github.com/xuxiaowei-com-cn/maven-dependencies/types"
)

func CurrentCommand() *cli.Command {
	return &cli.Command{
		Name:  "current",
		Usage: "Maven 当前坐标",
		Flags: []cli.Flag{
			flag.FilePathFlag(false),
			flag.AfterGroupIdFlag(), flag.AfterArtifactIdFlag(),
			flag.AfterVersionFlag(),
		},
		Subcommands: []*cli.Command{
			EditCurrentCommand(),
		},
	}
}

func EditCurrentCommand() *cli.Command {
	return &cli.Command{
		Name:  "edit",
		Usage: "修改 Maven 当前坐标",
		Flags: []cli.Flag{
			flag.FilePathFlag(true),
			flag.AfterGroupIdFlag(), flag.AfterArtifactIdFlag(),
			flag.AfterVersionFlag(),
		},
		Action: func(context *cli.Context) error {
			var filePath = context.Path(constant.FilePath)
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

			project, err := types.Dependencies(fileContext)
			if err != nil {
				return err
			}

			var result string
			for index, line := range fileContextLines {
				if index != 0 {
					result += "\n"
				}
				if project.GroupIdLine == index+1 && afterGroupId != "" {
					result += "<groupId>" + afterGroupId + "</groupId>"
				} else if project.ArtifactIdLine == index+1 && afterArtifactId != "" {
					result += "<artifactId>" + afterArtifactId + "</artifactId>"
				} else if project.VersionLine == index+1 && afterVersion != "" {
					result += "<version>" + afterVersion + "</version>"
				} else {
					result += line
				}
			}

			if project.Version == "" && afterVersion != "" {
				lines := strings.Split(result, "\n")
				result = ""

				for index, line := range lines {
					if index != 0 {
						result += "\n"
					}
					result += line
					if project.ArtifactIdLine == index+1 {
						result += "\n<version>" + afterVersion + "</version>"
					}
				}
			}

			if project.GroupId == "" && afterGroupId != "" {
				lines := strings.Split(result, "\n")
				result = ""

				for index, line := range lines {
					if index != 0 {
						result += "\n"
					}
					result += line
					if project.ArtifactIdLine == index+1 && afterGroupId != "" {
						result += "\n<groupId>" + afterGroupId + "</groupId>"
					}
				}
			}

			err = os.WriteFile(filePath, []byte(result), 0644)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
