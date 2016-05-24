package main

import (
	"flag"
	"fmt"
	"github.com/januswel/yangpao/core"
	"os"
)

func Xor(a, b bool) bool {
	return (a || b) && !(a && b)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, r)
		}
	}()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--major] [-m-|--minor] [-p|--patch]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Choose only one option from below to upgrade version\n")
		fmt.Fprintf(os.Stderr, "    -g|--generate   generates setting file .yangpao.toml\n")
		fmt.Fprintf(os.Stderr, "    --major         upgrades major version\n")
		fmt.Fprintf(os.Stderr, "    -m|--minor      upgrades major version\n")
		fmt.Fprintf(os.Stderr, "    -p|--patch      upgrades patch version\n")
	}
	currentVersionShort := flag.Bool("c", false, "")
	currentVersionLong := flag.Bool("current", false, "")
	generateShort := flag.Bool("g", false, "")
	generateLong := flag.Bool("generate", false, "")
	patchShort := flag.Bool("p", false, "")
	patchLong := flag.Bool("patch", false, "")
	minorShort := flag.Bool("m", false, "")
	minorLong := flag.Bool("minor", false, "")
	majorLong := flag.Bool("major", false, "")

	flag.Parse()

	needsToShowCurrentVersion := *currentVersionShort || *currentVersionLong
	needsToGenerateSettingFile := *generateShort || *generateLong
	patch := *patchShort || *patchLong
	minor := *minorShort || *minorLong
	major := *majorLong

	if err := core.SearchSettingFile(); err != nil {
		panic(err)
	}

	if needsToShowCurrentVersion {
		currentVersion, err := core.CurrentVersion()
		if err != nil {
			panic(err)
		}
		fmt.Println(currentVersion)
		os.Exit(0)
	}

	if needsToGenerateSettingFile {
		core.GenerateSettingFile()
		os.Exit(0)
	}

	if !(patch || minor || major) {
		ShowVersions()
		os.Exit(0)
	}

	if !Xor(Xor(patch, minor), major) || (patch && minor && major) {
		flag.Usage()
		os.Exit(2)
	}

	Upgrade(patch, minor, major)
}

func ShowVersions() {
	var versions core.Versions
	if err := core.CheckVersions(&versions); err != nil {
		panic(err)
	}

	fmt.Printf("current version: %s\n", versions.Current)
	for _, file := range versions.Files {
		fmt.Printf("%s\n", file.Path)
		for _, line := range file.Lines {
			fmt.Printf("    %s\n", line)
		}
	}
}

func Upgrade(patch, minor, major bool) {
	var which int

	if patch {
		which = core.PATCH
	}
	if minor {
		which = core.MINOR
	}
	if major {
		which = core.MAJOR
	}

	newVersion, err := core.Upgrade(which)
	if err != nil {
		panic(err)
	}

	fmt.Printf("upgraded to %s\n", newVersion)
}
