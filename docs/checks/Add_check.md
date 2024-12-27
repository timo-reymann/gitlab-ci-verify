Add a new check
===

There are two ways to add a new check:

- [**Rego check**](#rego-check): Implement the check in [Rego]. This is preferred for simple checks
- [**Go check**](#go-check): Implement the check in Go. This is preferred for more complex checks

## [Rego] check

Checks are implemented in [Rego]. The checks are stored in the `pkg/checks/rego` directory.

You need to write a new [Rego] file for the check in `pkg/checks/rego/<your check>.rego` and register the check.

1. Create a new file for the check in `pkg/checks/*_check.rego`, e.g. `my_check.go`
2. Implement the check in [Rego] e.g.
   ```rego
   package my_check

   findings contains finding if {
        artifact_paths := input.mergedYaml.pages.artifacts.paths

        count([artifact_paths |
            some artifact_path in artifact_paths
            startswith(artifact_path, "public")
        ]) == 0

        finding := gitlab_ci_verify.warning(
            "CHECK-123",
            "message",
            yamlPathToLineNumber(".pages.artifacts.paths"),
        )
   }
   ```
3. Wrap the rego check in a Go check e.g.
   ```go
   package checks
   
   type MyCheck struct {
     InMemoryCheck
   }
   
   //go:embed my_check.rego
   var myRegoCheck string
   
   func NewMyCheck() MyCheck {
      return MyCheck{
         InMemoryCheck{
             RegoContent: myRegoCheck,
         },
      }
   }
   ```
4. Add test files in `pkg/checks/testdata/<your check>`
5. Register the check in `pkg/checks/checks.go`
    ```go
    package checks

    func init() {
        // other checks
        RegisterCheck(NewMyCheck())
    }
    ```

## Go check

Checks are implemented via the [`Check` interface](../../pkg/checks/check.go).

## Step by step

1. Create a new file for the check in `pkg/checks/*_check.go`, e.g. `my_check.go`
2. Implement Check e.g.
   ```go
   package checks   

   type MyCheck struct {}
   
   func (m MyCheck) Run(i *CheckInput) ([]CheckFinding, error) {
        return []CheckFinding{}, nil
   }
   ```
3. Add test files in `pkg/checks/testdata/<your check>`
4. Add tests for the check in `pkg/checks/*_check_test.go` e.g. `my_check_test.go`
    ```go
     package main

     import "testing"

     func TestMyCheck(t *testing.T) {
        c := MyCheck{}
        testCases := []struct {
            name             string
            file             string
            expectedFindings []CheckFinding{}
            // add what you need for test cases
         }{
            {
                name: "My test",
                file: "myFile.yml",
                expectedFindings: []CheckFinding{},
            },
        }
            
        for _, tc := range testCases {
            t.Run(tc.name, func (t *testing.T) {
                verifyFindings(t, tc.expectedFindings, checkMustSucceed(c.Run(&CheckInput{
                    CiYaml:        newCiYamlFromFile(t, path.Join("test_data", "my_check", tc.file)),
                    Configuration: &cli.Configuration{},
                })))
            })
        }
    }  
    ```
5. Register the check in `pkg/checks/checks.go`
    ```go
    package checks

    func init() {
        // other checks
        RegisterCheck(MyCheck{})
    }
    ```
   
[Rego]: https://www.openpolicyagent.org/docs/latest/policy-language/#what-is-rego