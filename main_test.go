package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
)

type Project struct {
	Parent string

	GroupId         string
	GroupIdLine     int
	GroupIdLeft     int
	ArtifactId      string
	ArtifactIdLine  int
	ArtifactIdLeft  int
	Version         string
	VersionLine     int
	VersionLeft     int
	Packaging       string
	PackagingLine   int
	PackagingLeft   int
	Description     string
	DescriptionLine int
	DescriptionLeft int

	Dependencies []Dependency
}

type Dependency struct {
	GroupId        string
	GroupIdLine    int
	GroupIdLeft    int
	ArtifactId     string
	ArtifactIdLine int
	ArtifactIdLeft int
	Version        string
	VersionLine    int
	VersionLeft    int
}

func Test(t *testing.T) {

	fileContextLines, err := ReadFileLines("https://repo1.maven.org/maven2/io/xuxiaowei/nacos/nacos-core/2.3.1/nacos-core-2.3.1.pom")
	if err != nil {
		assert.NoError(t, err)
	}

	var project Project

	for index, line := range fileContextLines {

		var lineTrimSpace = strings.TrimSpace(line)

		if project.GroupId == "" {
			valueStr, valueLine, valueLeft := value(fileContextLines, index, lineTrimSpace, "groupId")
			project.GroupId = valueStr
			project.GroupIdLine = valueLine
			project.GroupIdLeft = valueLeft
		}

		if project.ArtifactId == "" {
			valueStr, valueLine, valueLeft := value(fileContextLines, index, lineTrimSpace, "artifactId")
			project.ArtifactId = valueStr
			project.ArtifactIdLine = valueLine
			project.ArtifactIdLeft = valueLeft
		}

		if project.Version == "" {
			valueStr, valueLine, valueLeft := value(fileContextLines, index, lineTrimSpace, "version")
			project.Version = valueStr
			project.VersionLine = valueLine
			project.VersionLeft = valueLeft
		}

		if project.Packaging == "" {
			valueStr, valueLine, valueLeft := value(fileContextLines, index, lineTrimSpace, "packaging")
			project.Packaging = valueStr
			project.PackagingLine = valueLine
			project.PackagingLeft = valueLeft
		}

		if project.Description == "" {
			valueStr, valueLine, valueLeft := value(fileContextLines, index, lineTrimSpace, "description")
			project.Description = valueStr
			project.DescriptionLine = valueLine
			project.DescriptionLeft = valueLeft
		}

		if project.Dependencies == nil {

			if strings.HasPrefix(lineTrimSpace, "<dependencies") {

				project.Dependencies = []Dependency{}

				var start = index + 1
				var end int

				for i := 0; i < len(fileContextLines)-index; i++ {

					var possibleValue = fileContextLines[index+i]

					var possibleValueTrimSpace = strings.TrimSpace(possibleValue)

					if possibleValueTrimSpace == "" {
						continue
					}

					if strings.HasSuffix(possibleValueTrimSpace, "</dependencies>") || strings.HasSuffix(possibleValueTrimSpace, "</dependencies") {
						end = index + i + 1
						break
					}
				}

				t.Log(start, end)

				for i := 0; i < len(fileContextLines)-index; i++ {

					var possibleValue = fileContextLines[index+i]

					t.Log(possibleValue)
				}

			}
		}

	}

	assert.Equal(t, "io.xuxiaowei.nacos", project.GroupId)
	assert.Equal(t, 4, project.GroupIdLine)
	assert.Equal(t, 11, project.GroupIdLeft)
	assert.Equal(t, "nacos-core", project.ArtifactId)
	assert.Equal(t, 5, project.ArtifactIdLine)
	assert.Equal(t, 14, project.ArtifactIdLeft)
	assert.Equal(t, "2.3.1", project.Version)
	assert.Equal(t, 6, project.VersionLine)
	assert.Equal(t, 11, project.VersionLeft)
	assert.Equal(t, "", project.Packaging)
	assert.Equal(t, 0, project.PackagingLine)
	assert.Equal(t, 0, project.PackagingLeft)
	assert.Equal(t, "nacos-core", project.Description)
	assert.Equal(t, 8, project.DescriptionLine)
	assert.Equal(t, 15, project.DescriptionLeft)
}

func value(lines []string, index int, lineTrimSpace string, name string) (string, int, int) {
	if strings.HasPrefix(lineTrimSpace, "<"+name) {

		var valueTmp string
		var valueTmpSum string

		for i := 0; i < len(lines)-index; i++ {
			var possibleValue = lines[index+i]
			valueTmpSum += possibleValue

			var possibleValueTrimSpace = strings.TrimSpace(possibleValue)

			if possibleValueTrimSpace == "" {
				continue
			}

			valueTmp += possibleValueTrimSpace

			if strings.HasSuffix(possibleValueTrimSpace, "</"+name+">") || strings.HasSuffix(possibleValueTrimSpace, "</"+name) {

				context, err := Context(valueTmp)
				if err != nil {
					panic(err)
				}

				var contextIndex = strings.Index(valueTmpSum, context)

				var contextIndexLeft = valueTmpSum[:contextIndex]
				var countLine = strings.Count(contextIndexLeft, "\r")
				var left = strings.Index(lines[countLine+index], context)

				return context, countLine + index + 1, left + 1
			}
		}
	}
	return "", 0, 0
}

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

func isURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

func get(path string) ([]byte, error) {

	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func ReadFile(path string) (string, error) {
	var bytes []byte
	var err error
	if isURL(path) {
		bytes, err = get(path)
	} else {
		bytes, err = os.ReadFile(path)
	}
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ReadFileLines(path string) ([]string, error) {
	str, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(str, "\n")
	return lines, nil
}
