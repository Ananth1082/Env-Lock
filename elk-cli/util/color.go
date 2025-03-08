package util

import "fmt"

// these are the basic 8 colors
const (
	ColorGreyF = 30 + iota
	ColorRedF
	ColorGreenF
	ColorYellowF
	ColorBlueF
	ColorMagentaF
	ColorCyanF
	ColorWhiteF

	ColorGreyB = 40 + iota
	ColorRedB
	ColorGreenB
	ColorYellowB
	ColorBlueB
	ColorMagentaB
	ColorCyanB
	ColorWhiteB

	MAX_COLORS = 16
)

const (
	FOREGROUND = iota
	BACKGROUND
)

/*
formatStr is the ANSI color code string with %s as the placeholder for the string to be colored
Name is the user defined name of the color
ColorType is either FOREGROUND or BACKGROUND
*/
type Color struct {
	formatStr string
	ColorType int
}

type RGBColor struct {
	Red   int
	Green int
	Blue  int
	Color
}

type SAColor struct {
	ColorCode int
	Color
}

type Color256 struct {
	ColorCode int
	Color
}

func NewRGBColor(r, g, b int, cType int) *RGBColor {
	fstr := ""
	if cType == FOREGROUND {
		fstr = "\033[38;2;%d;%d;%dm%s\033[0m"
	} else {
		fstr = "\033[48;2;%d;%d;%dm%s\033[0m"
	}
	return &RGBColor{
		Red:   r,
		Green: g,
		Blue:  b,
		Color: Color{
			formatStr: fmt.Sprintf(fstr, r, g, b, "%s"),
			ColorType: cType,
		},
	}
}

func NewSAColor(colorCode int) *SAColor {
	fstr := "\033[%dm%s\033[0m"
	ctype := FOREGROUND

	if colorCode >= 40 {
		ctype = BACKGROUND
	}
	return &SAColor{
		ColorCode: colorCode,
		Color: Color{
			formatStr: fmt.Sprintf(fstr, colorCode, "%s"),
			ColorType: ctype,
		},
	}
}

func NewColor256(colorCode int, cType int) *Color256 {
	fstr := "\033[38;5;%dm%s\033[0m"
	if cType == BACKGROUND {
		fstr = "\033[48;5;%dm%s\033[0m"
	}

	return &Color256{
		ColorCode: colorCode,
		Color: Color{
			formatStr: fmt.Sprintf(fstr, colorCode, "%s"),
			ColorType: cType,
		},
	}
}

func (c *Color) GetColoredString(str string) string {
	return fmt.Sprintf(c.formatStr, str)
}
