package process_settings

import "reflect"

type TargetEvaluator struct {
	targetingContext map[string]interface{}
}

func sliceContains(slice []interface{}, item interface{}) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func deepMatch(a, b interface{}) bool {
	if a == nil || b == nil {
		return false
	}

	typeOfA := reflect.TypeOf(a).Kind()
	typeOfB := reflect.TypeOf(b).Kind()

	if typeOfA != typeOfB && typeOfA != reflect.Slice {
		return false
	}

	switch typeOfA {
	case reflect.Slice:
		return sliceContains(a.([]interface{}), b)
	case reflect.Map:
		return mapContains(a.(map[string]interface{}), b.(map[string]interface{}))
	default:
		return a == b
	}
}

func mapContains(map1, map2 map[string]interface{}) bool {
	for key, value := range map1 {
		if !deepMatch(value, map2[key]) {
			return false
		}
	}

	return true
}

func (t *TargetEvaluator) IsTargetMatch(settingsFile SettingsFile) bool {
	if settingsFile.Metadata.End {
		return false
	}

	if settingsFile.Target == nil {
		return true
	}

	return mapContains(settingsFile.Target, t.targetingContext)
}
