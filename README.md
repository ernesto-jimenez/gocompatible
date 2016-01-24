# gocompatible (WORK IN PROGRESS)

Backwards compatibility is a big deal in Go, specially with the
limitations around package versioning.

`gocompatible` helps you keep your packages backwards compatible by
testing your changes across all your dependents.

Our aim is to:

 - Help you find what packages depend on yours.
 - Automate testing of all the packages dependent on yours.
 - Take a diff or a Pull Request and show whether any dependent packages
   will get broken by those changes.

## Installation

```
go get github.com/gophergala2016/gocompatible
```

## Usage

```txt
Usage:
  gocompatible [command]

Available Commands:
  dependents  List all packages we can find depending on the given package
  test        Run tests for all depending packages

Flags:
  -f, --filter string   Filter dependents within the given path
      --godoc           Find dependents in godoc.org
  -h, --help            help for gocompatible
  -l, --local           Find dependents in your $GOPATH (default true)
  -r, --recurisve       Check subpackages too
  -v, --verbose         Verbose output

Use "gocompatible [command] --help" for more information about a
command.
```

## Examples

Find all packages from my GOPATH dependent on the package in the
current working directory and test all of them:

```
gocompatible dependents
```

Find all packages from my GOPATH dependent on
`github.com/stretchr/testify/assert` and test them:

```
gocompatible dependents github.com/stretchr/testify/assert
```

List all dependent packages depedent on `assert` tracked by [godoc.org][godoc].

```
gocompatible dependents --godoc github.com/stretchr/testify/assert
```

List all dependent packages depedent on `testify`'s packages tracked by [godoc.org][godoc].

```
gocompatible dependents -r --godoc github.com/stretchr/testify
```

Find and test all packages depedent on `assert` tracked by
[godoc.org][godoc]. **this should be run in a container to avoid running
random internet code pulled from the internet in your machine**.

```
gocompatible test --godoc --insecure github.com/stretchr/testify/assert
```

# Todo

 - [ ] get a pull request URL and test all dependencies before and after the PR.

# Where does the idea come from?

While developing [testify][testify] we sometimes need to evaluate whether a
internal change could break the projects dependent on them.

I decided it would be a nice project to prototype for this year's [Gopher
Gala][gala].

[testify]: https://github.com/stretchr/testify
[godoc]: https://godoc.org
[gala]: http://gophergala.com
