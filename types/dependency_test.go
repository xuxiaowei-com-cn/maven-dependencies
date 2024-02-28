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

	parent, dependencies, err := Dependencies(str)
	if err != nil {
		t.Fatalf("字符串处理成 Dependency 异常: %v", err)
	}

	t.Log("parent", parent.GroupId, parent.GroupIdLine, parent.ArtifactId, parent.ArtifactIdLine, parent.Version, parent.VersionLine)
	for index, dependency := range dependencies {
		t.Log(index, dependency.GroupId, dependency.GroupIdLine, dependency.ArtifactId, dependency.ArtifactIdLine, dependency.Version, dependency.VersionLine)
	}

}
