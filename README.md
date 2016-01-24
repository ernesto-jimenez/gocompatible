# gocompatible

[Gopher Gala 2016][gala] entry from [@ernesto_jimenez][twitter]. It's
aimed at testing backwards compatiblity with known dependent packages.
It can also check your packages will still work once upgrading a
dependency.

You can try it out quickly using docker:

```txt
docker pull quay.io/ernesto_jimenez/gocompatible
docker run --rm quay.io/ernesto_jimenez/gocompatible gocompatible
```

Check the changes from testify v1.0 to v1.1.3 do not break any of the
aws-go-sdk packages:

```txt
docker run --rm quay.io/ernesto_jimenez/gocompatible \
        gocompatible diff github.com/stretchr/testify/... \
        --filter github.com/aws \
        --from v1.0 --to v1.1.1 \
        --godoc --insecure
FAIL    github.com/aws/aws-sdk-go/awstesting - go get: exit status 2
0 skipped due to tests failing in v1.0
1 packages 1 failed:    1 failed get    0 build    0 test
```

The `aws-sdk-go` broke in v1.1.1 due to a backwards incompatible change
in `testify`. You can run the command with `-v` to get the output.

Check the upgrade to v1.1.3:

```txt
docker run --rm quay.io/ernesto_jimenez/gocompatible \
        gocompatible diff github.com/stretchr/testify/... \
        --filter github.com/aws \
        --from v1.0 --to v1.1.3 \
        --godoc --insecure
ok      github.com/aws/aws-sdk-go/awstesting
0 skipped due to tests failing in v1.0
1 packages 0 failed
```

Since executing code from random people in the internet is not a good
idea `gocompatible` comes with safeguards:

```txt
$ gocompatible test --godoc github.com/stretchr/testify/...
DANGEROUS ACTION!

Never run _test_ or _diff_ commands with --godoc
in your machine. This will download untrusted code
and run it, which is very dangerous.

You must always run those commands within an
isolated container. We've published a docker
image you can use to do it quickly.

Example:

gocompatible test github.com/stretchr/testify/assert \
  --filter github.com/uber \
  --godoc --docker

Which is equivalent to:

docker run --rm quay.io/ernesto_jimenez/gocompatible:latest \
  gocompatible test github.com/stretchr/testify/assert \
  --filter github.com/uber \
  --godoc --insecure
```

`gocompatible` always prepares a new `GOPATH` in a temporary directory
and works there to avoid polluting yours.

# Instructions

```txt
Backwards compatibility is really important in Go, specially with the
limitations around package versioning.

_gocompatible_ helps you keep your packages backwards compatible:

 - Find packages that depend on yours from your GOPATH or from godoc.org
 - Run tests from all the packages depending on yours
 - Check whether any dependeng packages break on different revisions

# Security

Never run _test_ or _diff_ commands with --godoc
in your machine. This will download untrusted code
and run it, which is very dangerous.

You must always run those commands within an
isolated container. We've published a docker
image you can use to do it quickly.

Example:

gocompatible test github.com/stretchr/testify/assert \
  --filter github.com/uber \
  --godoc --docker

Which is equivalent to:

docker run --rm quay.io/ernesto_jimenez/gocompatible:latest \
  gocompatible test github.com/stretchr/testify/assert \
  --filter github.com/uber \
  --godoc --insecure

Usage:
  gocompatible [command]

Examples:
# Find all pkgs on GOPATH depending on .
gocompatible dependents

# Find all pkgs on GOPATH depending on ./...
gocompatible dependents ./...

# Find all pkgs on GOPATH depending on testify
gocompatible dependents github.com/stretchr/testify/...

# Find all pkgs tracked by godoc.org as importing testify
gocompatible dependents --godoc github.com/stretchr/testify/...

# Run test for all packages in GOPATH dependent on testify/assert
gocompatible test github.com/stretchr/testify/assert

# Run tests from all uber packages tracked by godoc.org as dependent
on testify/assert
gocompatible test github.com/stretchr/testify/assert \
  --filter github.com/uber \
  --godoc --insecure

# Check whether changes between testify v1.0 and
# v1.1.1 broke any tests in aws-sdk-go
gocompatible diff github.com/stretchr/testify/... \
  --filter github.com/aws/aws-sdk-go/awstesting \
  --from v1.0 --to v1.1.1

Available Commands:
  dependents  List all packages we can find depending on the given
package
  diff        Compare what depending packages break between two diffent
revisions
  test        Run tests for all depending packages

Flags:
  -d, --docker          run the command in a docker container
  -f, --filter string   select dependents within the given path
      --godoc           find dependents in godoc.org
  -h, --help            help for gocompatible
  -l, --local           find dependents in your $GOPATH (default true)
  -r, --recurisve       check subpackages too
  -v, --verbose         verbose output

Use "gocompatible [command] --help" for more information about a
command.
```

## Installation

```
go get github.com/gophergala2016/gocompatible
```

## Example runs

### Runnign a gocompatible diff

```
$ gocompatible diff github.com/stretchr/testify/... \
  --filter github.com/aws/aws-sdk-go/awstesting/integration \
  --from v1.0 --to v1.1.1
ok      github.com/aws/aws-sdk-go/awstesting/integration/customizations/s3
FAIL    github.com/aws/aws-sdk-go/awstesting/integration/smoke - go get:
exit status 2
0 skipped due to tests failing in v1.0
2 packages 1 failed:    1 failed get    0 build    0 test
```

### Using gocompatible test/diff within docker

You must not simply run `gocompile test --godoc` or `gocompile diff --godoc`
from your computer since you will be running code from random people on
the internet. Instead, you can use our docker image to run
`gocompatible` within a container.

```
$ docker run --rm quay.io/ernesto_jimenez/gocompatible \
  gocompatible test github.com/raphael/goa \
  --godoc --insecure
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


# gocompatible diff --help

```txt
Usage:
  gocompatible diff [options] --from <commit> [<package>] [flags]

Examples:
# Check whether changes between testify v1.0 and
# v1.1.1 broke any tests in aws-sdk-go
gocompatible diff github.com/stretchr/testify/... \
  --filter github.com/aws/aws-sdk-go/awstesting \
  --from v1.0 --to v1.1.1

Flags:
  -c, --from string   commit/tag/branch of <package> to checkout and select only the packages with green builds
      --insecure      allows running testing packages from godoc
  -t, --to string     commit/tag/branch of the tested package we want to compare with from to see how many packages broke (default "master")

Global Flags:
  -d, --docker          run the command in a docker container
  -f, --filter string   select dependents within the given path
      --godoc           find dependents in godoc.org
  -l, --local           find dependents in your $GOPATH (default true)
  -r, --recurisve       check subpackages too
  -v, --verbose         verbose output
```

# gocompile test --help

```txt
This can be useful to consider whether to upgrade one of your
dependencies
across your organisation.

e.g.: Imagine you work at uber and whant to check your packages running
testify
will work with the latests version on master. You could run:

gocompatible githb.com/stretchr/testify/... \
  --filter github.com/uber

This will find all the packages on your GOPATH depending on testify and
run
their tests with the latest version of testify.

You could also checkout a specific branch or tag from testify:

gocompatible githb.com/stretchr/testify/... \
  --checkout v1.1
  --filter github.com/uber

Usage:
  gocompatible test [options] [<package>] [flags]

Examples:
# Run test for all packages in GOPATH dependent on testify/assert
gocompatible test github.com/stretchr/testify/assert

# Run tests from all uber packages tracked by godoc.org as dependent on
testify/assert
gocompatible test github.com/stretchr/testify/assert \
  --filter github.com/uber \
  --godoc --insecure

Flags:
  -c, --checkout string   The commit/tag/branch of the tested package we
want to checkout (default "master")
      --insecure          Allows running testing packages from godoc

Global Flags:
  -d, --docker          run the command in a docker container
  -f, --filter string   select dependents within the given path
      --godoc           find dependents in godoc.org
  -l, --local           find dependents in your $GOPATH (default true)
  -r, --recurisve       check subpackages too
  -v, --verbose         verbose output
```

# gocompatible dependents --help

```txt
This is very useful for:

 - Finding local packages depending on certain package.
 e.g: to consider upgrading a package.
 - Finding packages tracked by godoc.org as importing certain package.
 e.g: to help you maintain backwards compatibility of your open source
package.

Usage:
  gocompatible dependents [<package>] [flags]

Examples:
# Find all pkgs on GOPATH depending on .
gocompatible dependents

# Find all pkgs on GOPATH depending on ./...
gocompatible dependents ./...

# Find all pkgs on GOPATH depending on testify
gocompatible dependents github.com/stretchr/testify/...

# Find all pkgs tracked by godoc.org as importing testify
gocompatible dependents --godoc github.com/stretchr/testify/...

Global Flags:
  -d, --docker          run the command in a docker container
  -f, --filter string   select dependents within the given path
      --godoc           find dependents in godoc.org
  -l, --local           find dependents in your $GOPATH (default true)
  -r, --recurisve       check subpackages too
  -v, --verbose         verbose output
```

# Where does the idea come from?

While developing [testify][testify] we sometimes need to evaluate whether a
internal change could break the projects dependent on them.

I decided it would be a nice project to prototype for this year's [Gopher
Gala][gala].

[testify]: https://github.com/stretchr/testify
[godoc]: https://godoc.org
[gala]: http://gophergala.com
[twitter]: https://twitter.com/ernesto_jimenez
