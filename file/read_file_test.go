package file

import (
	"testing"
)

func Test_ReadFile(t *testing.T) {

	path := "../LICENSE"
	//path := "https://repo1.maven.org/maven2/org/springframework/boot/spring-boot-dependencies/2.7.18/spring-boot-dependencies-2.7.18.pom"

	str, err := ReadFile(path)
	if err != nil {
		t.Fatalf("无法读取文件: %v", err)
	}

	t.Log(str)
}

func Test_ReadFileTrimSpace(t *testing.T) {
	path := "../LICENSE"
	//path := "https://repo1.maven.org/maven2/org/springframework/boot/spring-boot-dependencies/2.7.18/spring-boot-dependencies-2.7.18.pom"

	str, err := ReadFileTrimSpace(path)
	if err != nil {
		t.Fatalf("无法读取文件: %v", err)
	}

	t.Log(str)
}

func Test_ReadFileLines(t *testing.T) {

	path := "../LICENSE"
	//path := "https://repo1.maven.org/maven2/org/springframework/boot/spring-boot-dependencies/2.7.18/spring-boot-dependencies-2.7.18.pom"

	lines, err := ReadFileLines(path)
	if err != nil {
		t.Fatalf("无法读取文件: %v", err)
	}

	for index, line := range lines {
		t.Log(index, line)
	}
}

func Test_ReadFileLinesTrimSpace(t *testing.T) {

	path := "../LICENSE"
	// path := "https://repo1.maven.org/maven2/org/springframework/boot/spring-boot-dependencies/2.7.18/spring-boot-dependencies-2.7.18.pom"

	lines, err := ReadFileLinesTrimSpace(path)
	if err != nil {
		t.Fatalf("无法读取文件: %v", err)
	}

	for index, line := range lines {
		t.Log(index, line)
	}
}
