# tyk-apis

# Installation

```
go get github.com/TykTechnologies/tyk-apis
```

# Usage

```
tyk-apis tyk:targets="api"   paths="./test" output:dir=schemas
```

Where,

- `api` is the case insensitive name of struct you want to
generate schema for.

-  `./test` is the absolute path to the go package that
has the struct defined.

- `output:dir=schemas` says we will write generated open api documents inside `schemas` directory