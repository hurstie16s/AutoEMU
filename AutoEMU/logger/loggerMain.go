package logger

import (
	"fmt"
)

const reset = "\033[0m"
const red = "\033[31m"
const green = "\033[32m"
const yellow = "\033[33m"
const blue = "\033[34m"
const purple = "\033[35m"
const cyan = "\033[36m"

//const grey   = "\033[37m"
//const white  = "\033[97m"

var SuccessMessage = fmt.Sprintf("%sSuccess: %s", green, reset)
var InfoMessage = fmt.Sprintf("%sInfo: %s", cyan, reset)
var ConfigMessage = fmt.Sprintf("%sConfig: %s", blue, reset)
var warnMessage = fmt.Sprintf("%sWarning: %s", yellow, reset)
var ErrorMessage = fmt.Sprintf("%sError: %s", purple, reset)
var FatalMessage = fmt.Sprintf("%sFatal: %s", red, reset)
