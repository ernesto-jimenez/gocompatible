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

TODO

# Where does the idea come from?

While developing [testify][testify] we sometimes need to evaluate whether a
internal change could break the projects dependent on them.

I decided it would be a nice project to prototype during gophergala.

[testify]: https://github.com/stretchr/testify
