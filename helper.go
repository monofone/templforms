package templforms

import "strconv"

func isInvalidValue(err error, isFirst bool) string {
	if isFirst {
		return ""
	}
	if err != nil {
		return "true"
	}

	return "false"
}

func idAttribute(ID, name string) string {
	if ID == "" {
		return name
	}

	return ID
}

func numberAttribute(attr int) string {
	if attr > 0 {
		return strconv.Itoa(attr)
	}

	return ""
}
