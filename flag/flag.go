package flag

import (
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/maven-dependencies/constant"
)

func FilePathFlag(required bool) cli.Flag {
	return &cli.PathFlag{
		Name:     constant.FilePath,
		Usage:    "File Path",
		Required: required,
	}
}

func GroupIdFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.GroupId,
		Usage:    "Maven groupId",
		Required: required,
	}
}

func ArtifactIdFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.ArtifactId,
		Usage:    "Maven artifactId",
		Required: required,
	}
}

//func VersionFlag() cli.Flag {
//	return &cli.StringFlag{
//		Name:  constant.Version,
//		Usage: "Maven version",
//	}
//}

func AfterGroupIdFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  constant.AfterGroupId,
		Usage: "After Maven groupId",
	}
}

func AfterArtifactIdFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  constant.AfterArtifactId,
		Usage: "After Maven artifactId",
	}
}

func AfterVersionFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  constant.AfterVersion,
		Usage: "After Maven version",
	}
}
