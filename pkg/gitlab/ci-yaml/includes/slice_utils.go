package includes

import "slices"

// FilterByTypes filters the slice of includes by the given types
func FilterByTypes(includes []Include, includeTypes ...string) []Include {
	var filteredIncludes []Include

	for _, include := range includes {
		if slices.Contains(includeTypes, include.Type()) {
			filteredIncludes = append(filteredIncludes, include)
		}
	}

	return filteredIncludes
}

// Equals checks if two include slices are equal
func Equals(a, b []Include) bool {
	if len(a) != len(b) {
		return false
	}

	for idx, include := range a {
		if !include.Equals(b[idx]) {
			return false
		}
	}
	return true
}

// Unique returns a slice with unique elements in order
func Unique(includes []Include) []Include {
	var uniqueIncludes []Include

	for _, include := range includes {
		alreadyIncluded := false
		for _, existingInclude := range uniqueIncludes {
			if include.Equals(existingInclude) {
				alreadyIncluded = true
				break
			}
		}
		if !alreadyIncluded {
			uniqueIncludes = append(uniqueIncludes, include)
		}
	}

	return uniqueIncludes
}
