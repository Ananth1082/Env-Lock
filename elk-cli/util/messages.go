package util

import "fmt"

var red = NewSAColor(ColorRedF)
var green = NewSAColor(ColorGreenF)
var yellow = NewSAColor(ColorYellowF)
var blue = NewSAColor(ColorBlueF)

func PrintError(errMsg string) {
	fmt.Println(red.GetColoredString(errMsg))
}

func PrintSuccess(successMsg string) {
	fmt.Println(green.GetColoredString(successMsg))
}

func PrintWarning(warningMsg string) {
	fmt.Println(yellow.GetColoredString(warningMsg))
}

func PrintPrompt(promptMsg string) {
	fmt.Print(blue.GetColoredString(promptMsg))
}
