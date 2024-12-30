package checks

func RegisterRemoteOPABundleChecks(bundleURLs []string) {
	for _, url := range bundleURLs {
		register(RemoteBundleCheck{
			BundleURL: url,
		})
	}
}
