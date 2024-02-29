package command

import (
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/maven-dependencies/constant"
	"github.com/xuxiaowei-com-cn/maven-dependencies/file"
	"github.com/xuxiaowei-com-cn/maven-dependencies/flag"
	"github.com/xuxiaowei-com-cn/maven-dependencies/types"
	"os"
	"strings"
)

func DescriptionCommand() *cli.Command {
	return &cli.Command{
		Name:  "description",
		Usage: "Maven 描述",
		Flags: []cli.Flag{
			flag.FilePathFlag(false),
			flag.AfterDescriptionFlag(false),
		},
		Subcommands: []*cli.Command{
			EditDescriptionCommand(),
		},
	}
}

func EditDescriptionCommand() *cli.Command {
	return &cli.Command{
		Name:  "edit",
		Usage: "修改 Maven 描述",
		Flags: []cli.Flag{
			flag.FilePathFlag(true),
			flag.AfterDescriptionFlag(true),
		},
		Action: func(context *cli.Context) error {
			var filePath = context.Path(constant.FilePath)
			var afterDescription = context.Path(constant.AfterDescription)

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
				if project.DescriptionLine == index+1 {
					result += "<description>" + afterDescription + "</description>"
				} else {
					result += line
				}
			}

			if project.Description == "" {
				lines := strings.Split(result, "\n")
				result = ""

				for index, line := range lines {
					if index != 0 {
						result += "\n"
					}
					result += line
					if project.ArtifactIdLine == index+1 {
						result += "\n<description>" + afterDescription + "</description>"
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
