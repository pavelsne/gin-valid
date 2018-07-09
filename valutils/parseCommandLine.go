package valutils

import "strconv"

func isNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
