package state

// MergeStateWithData will merge data from state with other data
// the data from State will be overwritten
func MergeStateWithData(stateData map[string]interface{}, data ...map[string]interface{}) map[string]interface{} {
	mergedData := make(map[string]interface{}, len(stateData))
	for k, v := range stateData {
		mergedData[k] = v
	}

	if len(data) > 0 {
		for k, v := range data[0] {
			if _, found := stateData[k]; !found {
				mergedData[k] = v
			}
		}

	}

	return mergedData
}
