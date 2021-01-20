# tyk-apis

# Installation

```
go get github.com/TykTechnologies/tyk-apis
```

# Usage

```
tyk-apis tyk:targets="apidefinition"   paths="./apidef"
```

Where, `apidefinition` is the case insensitive name of struct you want to
generate schema for, and `./apidef` is the absolute path to the go package that
has the struct defined
