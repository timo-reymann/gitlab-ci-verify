Add a new check
===

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