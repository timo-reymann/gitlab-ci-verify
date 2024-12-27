# Integration of custom policies

* Status: proposed
* Date: 2024-12-23

Technical Story: https://github.com/timo-reymann/gitlab-ci-verify/issues/34

## Context and Problem Statement

Currently, gitlab-ci-verify only supports check written in Golang. This means every check needs to be implemented in Go.
This is a high barrier for new contributors and also limits the flexibility of the tool.

Additionally, it makes it basically impossible to add project or organization level checks, as these are highly dynamic
and need to be configurable.

## Decision Drivers <!-- optional -->

* battle tested solution
* easy to write and maintain checks
* support for project and organization level checks

## Considered Options

* [Rego] checks, integrated using the Go SDK
* Interpreted go using [yaegi]
* Configuration file using [expr]

## Decision Outcome

Chosen option: "rego",
because it is already well established in the cloud native ecosystem and provides a good balance between flexibility and
ease of use.

## Pros and Cons of the Options <!-- optional -->

### Rego

Rego is a policy language that is part of the Open Policy Agent project. It is designed to be easy to write and maintain
policies. It can be integrated into Go applications using the Go SDK.

* Good, because it is a battle-tested solution
* Good, because it is easy to write and maintain checks
* Good, because it supports project and organization level checks through remote bundles
* Bad, because it requires learning a new language

### yaegi

yaegi is an embedded Go interpreter. It can be used to execute Go code at runtime. This would allow users to write
checks in Go using the same syntax as the rest of the tool.

It is used by traefik to allow users to write plugins in Go.

* Good, because it allows users to write checks in Go
* Bad, because it requires learning a new language

### expr

expr is a library that allows users to evaluate Go-like expressions at runtime. This would allow users to write checks
in a configuration file using a simplified Go-like syntax.

* Good, because it allows users to write checks in a configuration file
* Bad, because it requires learning a new language
* Bad, because it is not very flexible and can become messy for a bit more complex checks

## Links <!-- optional -->

* [Rego]
* [yaegi]
* [expr]

[Rego]: https://www.openpolicyagent.org/docs/latest/policy-language/#what-is-rego

[yaegi]: https://github.com/traefik/yaegi

[expr]: https://github.com/expr-lang/expr

<!-- markdownlint-disable-file MD013 -->