# gocompatible

`gocompatible` helps you check wheter some new developments in your
package will the packages depending on it.

When you have some changes to your open source package, `gocompatible`
will look for other public packges depending on it and check whether
those packages will be broken with your changes.

This is COSI! (Continuous Open Source Integration)

## Installation

```
go get github.com/gophergala2016/gocompatible
```

## Usage

Print a list of all packages from my GOPATH dependent on the package in
my current directory:

```
gocompatible -print
```

Find all packages from my GOPATH dependent on the package in the
current working directory and test all of them:

```
gocompatible
```

Find all packages from my GOPATH dependent on
`github.com/stretchr/testify/assert` and test them:

```
gocompatible -pkg github.com/stretchr/testify/assert
```

List all dependent packages depedent on `assert` tracked by [godoc.org][godoc].

```
gocompatible -godoc -print -pkg github.com/stretchr/testify/assert
```

List all dependent packages depedent on `testify`'s packages tracked by [godoc.org][godoc].

```
gocompatible -godoc -print -pkg github.com/stretchr/testify -subpkgs
```

Find and test all packages depedent on `assert` tracked by
[godoc.org][godoc]. **this should be run in a container to avoid running
random internet code pulled from the internet in your machine**.

```
gocompatible -godoc -pkg github.com/stretchr/testify/assert
```

# Todo

 - [ ] get a pull request URL and test all dependencies before and after the PR.

# Where does the idea come from?

While developing [testify][testify] we sometimes need to evaluate whether a
internal change could break the projects dependent on them.

I decided it would be a nice project to prototype during gophergala.

[testify]: https://github.com/stretchr/testify
[godoc]: https://godoc.org
