package file

import (
	"testing"
)

func Test_ReadFile(t *testing.T) {

	filePath := "../LICENSE"

	str, err := ReadFile(filePath)
	if err != nil {
		t.Fatalf("无法读取文件: %v", err)
	}

	t.Log(str)
}

func Test_ReadFileLines(t *testing.T) {

	filePath := "../LICENSE"

	lines, err := ReadFileLines(filePath)
	if err != nil {
		t.Fatalf("无法读取文件: %v", err)
	}

	for index, line := range lines {
		t.Log(index, line)
	}
}

func Test_ReadFileLinesTrimSpace(t *testing.T) {

	filePath := "../LICENSE"

	lines, err := ReadFileLinesTrimSpace(filePath)
	if err != nil {
		t.Fatalf("无法读取文件: %v", err)
	}

	for index, line := range lines {
		t.Log(index, line)
	}
}
