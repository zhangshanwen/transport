package utils

import (
	"fmt"
	"strings"
)

// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//

var (
	foregroundMap = map[string]string{
		"black":        "30",
		"red":          "31",
		"green":        "32",
		"yellow":       "33",
		"blue":         "34",
		"purplish_red": "35",
		"ultramarine":  "36",
	}
)

func Print(color string, msg ...string) {
	color = strings.ToLower(color)
	color = foregroundMap[color]
	if color == "" {
		color = "30"
	}
	msg = append([]string{fmt.Sprintf("\u001B[1;%s;40m", color)}, msg...)
	fmt.Printf(strings.Join(msg, ""))
}

func Println(color string, msg ...string) {
	msg = append(msg, "\n")
	Print(color, msg...)
}
