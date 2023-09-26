package services

// unique removes duplicate records
func unique(stringSlice []string) []string {
	// create a map with all the values as key
	keys := make(map[string]bool)
	returnSlice := []string{}
	for _, item := range stringSlice {
		if _, value := keys[item]; !value {
			keys[item] = true
			returnSlice = append(returnSlice, item)
		}
	}
	return returnSlice
}
