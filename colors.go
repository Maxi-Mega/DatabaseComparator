package main

const enableColors = true // TODO: make this optional

const (
	colorReset         = "\033[0m"
	colorMissing       = "\033[31m"
	colorSame          = "\033[32m"
	colorDifferent     = "\033[34m"
	colorPartiallySame = "\033[36m"
	colorCommon        = "\033[37m"
)

// colorForMissing returns the given string colored as missing
func colorForMissing(str string) string {
	if enableColors {
		return colorMissing + str + colorReset
	} else {
		return str
	}
}

// colorForSame returns the given string colored as same
func colorForSame(str string) string {
	if enableColors {
		return colorSame + str + colorReset
	} else {
		return str
	}
}

// colorForDifferent returns the given string colored as different
func colorForDifferent(str string) string {
	if enableColors {
		return colorDifferent + str + colorReset
	} else {
		return str
	}
}

// colorForPartiallySame returns the given string colored as partially the same
func colorForPartiallySame(str string) string {
	if enableColors {
		return colorPartiallySame + str + colorReset
	} else {
		return str
	}
}

// colorForCommon returns the given string colored as common
func colorForCommon(str string) string {
	if enableColors {
		return colorCommon + str + colorReset
	} else {
		return str
	}
}
