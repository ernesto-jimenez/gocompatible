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
  diff        Compare what depending packages break between two diffent revisions
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
$ gocompatible dependents
```

Find all packages from my GOPATH dependent on
`github.com/stretchr/testify/assert` and test them:

```
$ gocompatible dependents github.com/stretchr/testify/assert
```

List all dependent packages depedent on `assert` tracked by [godoc.org][godoc].

```
$ gocompatible dependents --godoc github.com/stretchr/testify/assert
```

List all dependent packages depedent on `testify`'s packages tracked by [godoc.org][godoc].

```
$ gocompatible dependents --godoc github.com/stretchr/testify/...
```

Find and test all packages depedent on `assert` tracked by
[godoc.org][godoc]. **this should be run in a container to avoid running
random internet code pulled from the internet in your machine**.

```
$ gocompatible test --godoc --insecure github.com/stretchr/testify/assert
```


```
$ gocompatible diff github.com/stretchr/testify/... --filter github.com/aws/aws-sdk-go/awstesting/integration --from v1.0 --to v1.1.1
ok
github.com/aws/aws-sdk-go/awstesting/integration/customizations/s3
FAIL    github.com/aws/aws-sdk-go/awstesting/integration/smoke - go get:
exit status 2
0 skipped due to tests failing in v1.0
2 packages 1 failed:    1 failed get    0 build    0 test
```

# Using gocompatible test/diff within docker

You must not simply run `gocompile test --godoc` or `gocompile diff --godoc`
from your computer since you will be running code from random people on
the internet. Instead, you can use our docker image to run
`gocompatible` within a container.

```
$ docker run --rm quay.io/ernesto_jimenez/gocompatible gocompatible test github.com/raphael/goa --godoc --insecure
FAIL    github.com/deferpanic/dpgoa/example - go get: exit status 2
ok      github.com/deferpanic/dpgoa/middleware
ok      github.com/goadesign/goa-cellar
FAIL    github.com/gopheracademy/congo - go get: exit status 2
ok      github.com/raphael/goa-cellar
FAIL    github.com/raphael/goa-middleware/cors - go test: exit status 1
ok      github.com/raphael/goa-middleware/gzip
FAIL    github.com/raphael/goa-middleware/jwt - go test: exit status 1
FAIL    github.com/raphael/goa-middleware/middleware - go get: exit status 1
FAIL    github.com/raphael/goa/examples/cellar - go get: exit status 1
FAIL    github.com/deferpanic/dpgoa/example/client/deferpanic-cli - go get: exit status 1
FAIL    github.com/raphael/goa/examples/cellar/swagger - go get: exit status 1
12 packages 8 failed:    6 failed get    0 build    2 test
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
