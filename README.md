揚炮 yangpao
============

version number raiser for [semver](http://semver.org/)
inspired by [bmp](https://github.com/kt3k/bmp)

synopsis
--------

- raises version numbers in your project files in keeping with a setting file
    - major version
    - minor version
    - patch version
- `commit`s whole changes in target files, and creates a `tag` for the commit with new version number using git
    - see "requirements"


requirements
------------

- git
    - If you need to commit changes by yangpao with `--tag`, `-t` option.
    - Make sure that git is installed and the path to git binary is added into $PATH


install
-------

[release page](https://github.com/januswel/yangpao/releases)

1. getting executable
    - Get a newest binary for your os/arch from release page.
    - Rename it to `yangpao`.
    - Put your `yangpao` into a directory within your $PATH.
2. generating a setting file
    - Change the directory to your project root
    - Run `yangpao -g`


usage
-----

### .yangpao.yml

Edit your .yangpao.yml

```yml
---
current: 2.1.3
paths:
  # matches like "2.1.3" in README.md
  - README.md

  # matches like "v2.1.3" in README.md
  - release_tag
    prefix: v

  # matches like "ver 2.1.3" in README.md
  - release_tag
    prefix: 'ver '

  # matches like "2.1.3 version" in README.md
  - release_tag
    postfix: ' version'
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
# breaking change
yangpao --major

# raise a patch version and commit the changes in one shot
yangpao --patch --tag
yangpao -pt
```
