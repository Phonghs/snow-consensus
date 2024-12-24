package util

func GetKeys[K comparable, V any](inputMap map[K]V) []K {
	if inputMap == nil || len(inputMap) == 0 {
		return []K{}
	}

	keys := make([]K, 0, len(inputMap))
	for key := range inputMap {
		keys = append(keys, key)
	}

	return keys
}

func GetMapExcludingKey[K comparable, V any](originalMap map[K]V, excludeKey K) map[K]V {
	resultMap := make(map[K]V)

	for key, value := range originalMap {
		if key != excludeKey {
			resultMap[key] = value
		}
	}

	return resultMap
}
