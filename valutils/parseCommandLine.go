package valutils

import "strconv"

// IsValidPort checks whether a handed in interface
// can be parsed into a valid server port number.
func IsValidPort(port interface{}) bool {
	if port == nil {
		return false
	}
	checkPort, ok := port.(string)
	if !ok || !isNumeric(checkPort) {
		return false
	}
	if len(checkPort) != 4 {
		return false
	}
	return true
}

func isNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
