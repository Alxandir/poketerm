package term

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
)

var reader = bufio.NewReader(os.Stdin)

func ShowNoResponseDialog(str string, args ...interface{}) string {
	finalString := fmt.Sprintf(str, args...)
	finalString = fmt.Sprintf("%-60v", finalString)
	return showDialog(false, finalString)
}

func ShowNoWaitDialog(str string, args ...interface{}) {
	showDialogWithoutConfirmation(str, args...)
	return
}

func lengthCheck(str string) string {
	if len(str) <= 0 || len(str) > 16 {
		return "Input is an invalid length"
	}
	return ""
}

func ShowInputDialog(str string, args ...interface{}) string {
	return ShowInputDialogValidated(lengthCheck, str, args...)
}

func ShowInputDialogValidated(check func(string) string, str string, args ...interface{}) string {
	result := showDialog(true, str, args...)
	validationError := check(result)
	if len(validationError) > 0 {
		fmt.Println(validationError)
		return ShowInputDialogValidated(check, str, args...)
	}
	return result
}

func showDialog(expectResponse bool, str string, args ...interface{}) (output string) {
	str += "\t"
	if !expectResponse {
		str += "\t\t(Press Enter to continue)"
	}
	fmt.Printf(str, args...)
	output, _ = reader.ReadString('\n')
	output = strings.Replace(output, "\n", "", -1)
	return
}

func showDialogWithoutConfirmation(str string, args ...interface{}) {
	fmt.Printf(str, args...)
	return
}

func PrettyPrint(str string) {
	myFigure := figure.NewFigure(str, "", true)
	myFigure.Print()
}
