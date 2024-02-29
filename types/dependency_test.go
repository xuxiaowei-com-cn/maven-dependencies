package types

import (
	"github.com/xuxiaowei-com-cn/maven-dependencies/file"
	"testing"
)

func Test_Dependencies(t *testing.T) {

	path := "https://repo1.maven.org/maven2/org/springframework/boot/spring-boot-dependencies/2.7.18/spring-boot-dependencies-2.7.18.pom"

	str, err := file.ReadFileTrimSpace(path)
	if err != nil {
		t.Fatalf("无法读取文件: %v", err)
	}

	project, err := Dependencies(str)
	if err != nil {
		t.Fatalf("字符串处理成 Dependency 异常: %v", err)
	}

	t.Log("parent", project.Parent.GroupId, project.Parent.GroupIdLine, project.Parent.ArtifactId, project.Parent.ArtifactIdLine, project.Parent.Version, project.Parent.VersionLine)
	t.Log("current", project.GroupId, project.GroupIdLine, project.ArtifactId, project.ArtifactIdLine, project.Version, project.VersionLine)
	for index, dependency := range project.Dependencies {
		t.Log(index, dependency.GroupId, dependency.GroupIdLine, dependency.ArtifactId, dependency.ArtifactIdLine, dependency.Version, dependency.VersionLine)
	}

}
