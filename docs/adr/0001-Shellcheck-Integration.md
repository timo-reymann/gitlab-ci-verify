---
id: 001
status: accepted
date: 2024-08-29
deciders:
  - timo-reymann
---
# ShellCheck Integration

## Context and Problem Statement

To make the CLI tool useful and support shell script linting it is necessary to integrate shellcheck which is quite
mature and covers a lot of ground.

## Decision Drivers <!-- optional -->

* integration is straightforward
* embeddable into the application (no dependencies to install for the user)

## Considered Options

* Link the statically compiled library from shellcheck
* Use memexec and embed the binary directly

## Decision Outcome

Chosen option: "Use memexec and embed the binary directly," because it does not require to write any haskell or c code.

## Pros and Cons of the Options <!-- optional -->

### Link the statically compiled library from shellcheck

Shellcheck is written in Haskell and built via cabal-install.

Fortunately the build can produce a static library `libshellcheck.a` for all supported platforms.

* Good, because there is no hacky memexec workaround and the memory is shared with the application
* Good, because no need to parse the output as the data structures can be directly used
* Bad, because the libraries are not part of an release so they have to be built manually which takes some time and
  resources
* Bad, because that makes it harder to embed when done dynamically and when done statically requires a dedicated project
  to provide them

### Use memexec and embed the binary directly

On Linux & macOS it is pretty convenient to make a binary from memory executable and treat it like a regular file. For
Windows the workaround is to copy it to a temporary folder on usage and deleting it afterward.

Fortunately for Go, there exists a library that does exactly this.

* Good, because it is super easy to integrate as the binary can be included as is
* Good, because binaries can be directly copied to the repository, making it easy to build for everyone
* Good, because there is no need to use the library and interact with Haskell over CGO which allows static compilation
* Good, because the CLI interface is considerable more stable than the internal library
* Bad, because version updates need to be done manually by copying files
* Bad, because the output needs to be parsed again to be usable
* Bad, because it increases memory usage

## Links <!-- optional -->

* [go-memexec](https://github.com/amenzhinsky/go-memexec)
* [shellcheck](https://github.com/koalaman/shellcheck)
* [cabal-install](https://hackage.haskell.org/package/cabal-install)
