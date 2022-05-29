package util

func MergeMaps[V interface{}](maps ...map[string]V) map[string]V {
	switch len(maps) {
	case 0:
		return nil
	case 1:
		return maps[0]
	}

	population := make(map[string]V)
	for _, m := range maps {
		if m == nil {
			continue
		}

		for key, value := range m {
			population[key] = value
		}
	}
	return population
}
