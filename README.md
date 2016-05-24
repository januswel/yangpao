揚炮 yangpao v0.0.1
===================

version number raiser for [semver](http://semver.org/)

inspired by [bmp](https://github.com/kt3k/bmp)


synopsis
--------

- raises version numbers in your project files in keeping with a setting file
    - major version
    - minor version
    - patch version


requirements
------------

- yangpao dose not need any requirements.


install
-------

[release page](https://github.com/januswel/yangpao/releases)

1. get an executable
    - binary download
        - Get a newest binary for your os/arch from release page.
        - Rename it to `yangpao`.
        - Put your `yangpao` into a directory within your $PATH.
    - go get
        - `go get github.com/januswel/yangpao`
2. generate a setting file
    - Change the directory to your project root
    - Run `yangpao -g`


usage
-----

### .yangpao.toml

Edit your .yangpao.toml

```toml:.yangpao.toml
current = "2.1.3"

[[files]]
# matches like "2.1.3" in README.md
path = "README.md"

[[files]]
# matches like "v2.1.3" in release_tag
path = "release_tag"
prefix = "v"

[[files]]
# matches like "ver 2.1.3" in version.txt
path = "assets/version.txt"
prefix = "ver "

[[files]]
# matches like "yangpao 2.1.3 version" in src/public/index.html
path = "src/public/index.html"
prefix = "yangpao "
postfix = " version"
```

### yangpao

```sh
# shows current version and checks consistency
yangpao

# generates setting file on current directory
yangpao --generate
yangpao -g

# raises patch version number
yangpao --patch
yangpao -p

# raises minor version number
yangpao --minor
yangpao -m

# raises major version number
# short option is not provided to prevent operational erros
# because this operation means the version upgrade includes incompatible changes
yangpao --major
```
