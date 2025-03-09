package checks

// Check that runs verifications
type Check interface {
	// Run the check
	Run(i *CheckInput) ([]CheckFinding, error)
}
