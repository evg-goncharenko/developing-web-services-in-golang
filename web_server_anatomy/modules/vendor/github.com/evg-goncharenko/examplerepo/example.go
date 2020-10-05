package examplerepo

// FirstName is the key in map
const FirstName string = "FirstName"

// GiveName return name and flat gtat in exists in code
func GiveName(m map[string]string) (string, bool) {
	val, exist := m[FirstName]
	return val, exist
}


