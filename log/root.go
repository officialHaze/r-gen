package log

import (
	"fmt"
)

type GenericSlugEnum struct {
	ERROR   string
	SUCCESS string
}

var GenericSlug GenericSlugEnum = GenericSlugEnum{
	ERROR:   "ERROR",
	SUCCESS: "SUCCESS",
}

func Successf(format string, v ...any) {
	m := fmt.Sprintf("[%s] %s\n", GenericSlug.SUCCESS, format)
	fmt.Printf(m, v...)
}

func Successln(msg string) {
	m := fmt.Sprintf("[%s] %s\n", GenericSlug.SUCCESS, msg)
	fmt.Println(m)
}

func Errorf(format string, v ...any) {
	m := fmt.Sprintf("[%s] %s\n", GenericSlug.ERROR, format)
	fmt.Printf(m, v...)
}

func Errorln(msg string) {
	m := fmt.Sprintf("[%s] %s\n", GenericSlug.ERROR, msg)
	fmt.Println(m)
}

func Printf(format string, v ...any) {
	fmt.Printf(format+"\n", v...)
}

func Println(format any) {
	fmt.Println(format)
}
