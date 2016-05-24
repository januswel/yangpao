package core

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Settings struct {
	Current string
	Files   []File
}

type File struct {
	Path    string
	Prefix  string
	Postfix string
}

type Versions struct {
	Current string
	Files   []Version
}

type Version struct {
	Path  string
	Lines []string
}

const SETTINGS_FILE_NAME = ".yangpao.toml"
const VERSION_PATTERN = `\d+\.\d+\.\d+`
const (
	PATCH = iota
	MINOR
	MAJOR
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func SearchSettingFile() error {
	previous := ""
	current, err := os.Getwd()
	if err != nil {
		return err
	}

	for previous != current {
		if Exists(SETTINGS_FILE_NAME) {
			return nil
		}

		previous = current
		err = os.Chdir("..")
		if err != nil {
			return err
		}
		current, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	return fmt.Errorf("%s is not found. Run %s -g to generate", SETTINGS_FILE_NAME, os.Args[0])
}

func GenerateSettingFile() {
	var settings Settings

	settings.Current = "0.0.0"

	var file File

	file.Path = "README.md"
	file.Prefix = "AwesomeProject "
	file.Postfix = " version"
	settings.Files = append(settings.Files, file)

	file.Path = "release_tag"
	file.Prefix = "v"
	file.Postfix = ""
	settings.Files = append(settings.Files, file)

	WriteBackSettings(settings)
}

func CheckVersions(versions *Versions) error {
	var settings Settings
	if err := ParseSettings(SETTINGS_FILE_NAME, &settings); err != nil {
		return err
	}

	versions.Current = settings.Current

	for _, file := range settings.Files {
		if !Exists(file.Path) {
			return fmt.Errorf("file is not exist: %s", file.Path)
		}
		var version Version
		version.Path = file.Path

		pattern := fmt.Sprintf("%s%s%s", file.Prefix, VERSION_PATTERN, file.Postfix)
		matcher := regexp.MustCompile(pattern)

		raw, err := ioutil.ReadFile(file.Path)
		if err != nil {
			return err
		}

		founds := matcher.FindAllString(string(raw), -1)
		version.Lines = founds

		versions.Files = append(versions.Files, version)
	}

	return nil
}

func Upgrade(which int) (string, error) {
	var settings Settings
	if err := ParseSettings(SETTINGS_FILE_NAME, &settings); err != nil {
		return "", err
	}

	newVersion, err := UpgradeVersion(settings.Current, which)
	if err != nil {
		return "", err
	}
	settings.Current = newVersion

	Replace(settings)
	WriteBackSettings(settings)

	return settings.Current, nil
}

func ParseSettings(settingsFileName string, settings *Settings) error {
	if _, err := toml.DecodeFile(SETTINGS_FILE_NAME, &settings); err != nil {
		return err
	}

	if settings.Current == "" {
		return fmt.Errorf("current version is empty or undefined")
	}

	r := regexp.MustCompile(VERSION_PATTERN)
	if !r.MatchString(settings.Current) {
		return fmt.Errorf("version is not like semver")
	}

	for _, file := range settings.Files {
		if file.Path == "" {
			return fmt.Errorf("path for each file is requied")
		}
	}

	return nil
}

func UpgradeVersion(current string, which int) (string, error) {
	switch which {
	case PATCH:
		return IncrementSpecifiedVersion(current, 2)
	case MINOR:
		return IncrementSpecifiedVersion(current, 1)
	case MAJOR:
		return IncrementSpecifiedVersion(current, 0)
	}

	return "", fmt.Errorf("specify patch, minor, or major to upgrade")
}

func IncrementSpecifiedVersion(current string, index int) (string, error) {
	split := strings.Split(current, ".")
	version, err := strconv.Atoi(split[index])
	if err != nil {
		return "", err
	}
	split[index] = strconv.Itoa(version + 1)
	for i := range split {
		if index < i {
			split[i] = "0"
		}
	}
	return strings.Join(split, "."), nil
}

func WriteBackSettings(settings Settings) error {
	file, err := os.OpenFile(SETTINGS_FILE_NAME, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	if err := encoder.Encode(settings); err != nil {
		return err
	}

	if _, err := file.Write(buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

func Replace(settings Settings) error {
	for _, file := range settings.Files {
		if !Exists(file.Path) {
			return fmt.Errorf("file is not exist: %s", file.Path)
		}

		raw, err := ioutil.ReadFile(file.Path)
		if err != nil {
			return err
		}

		pattern := fmt.Sprintf("%s(%s)%s", file.Prefix, VERSION_PATTERN, file.Postfix)
		matcher := regexp.MustCompile(pattern)
		target := fmt.Sprintf("%s%s%s", file.Prefix, settings.Current, file.Postfix)

		replaced := matcher.ReplaceAllString(string(raw), target)

		err = ioutil.WriteFile(file.Path, []byte(replaced), 0)
		if err != nil {
			return err
		}
	}

	return nil
}
