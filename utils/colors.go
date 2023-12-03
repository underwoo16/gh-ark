package utils

import (
	"fmt"

	"github.com/mgutz/ansi"
)

var Yellow = ansi.ColorFunc("222+h")
var Blue = ansi.ColorFunc("75+h")
var Green = ansi.ColorFunc("120+h")

func TestColors() string {
	// iterate range 0 .. 255
	for i := 0; i < 256; i++ {
		color := fmt.Sprint(i) + "+h"
		str := fmt.Sprintf("%d", i)
		fmt.Println(ansi.Color(str, color))
	}
}
