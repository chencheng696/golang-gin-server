package main

//供template使用
import (
	"html/template"
	"strconv"
	"strings"
	"time"
)

func tmplMapValue(m map[string]string, k string) string {
	if v, ok := m[k]; ok {
		return v
	} else {
		return ``
	}
}

func tmplFormatAsDate(t time.Time, format string) string {
	if t.IsZero() {
		return ``
	}
	if format == `yyyy-MM-dd hh:mm:ss` {
		return t.Format("2006-01-02 15:04:05")
	} else if format == `yyyy-MM-dd hh:mm` {
		return t.Format("2006-01-02 15:04")
	} else if format == `yyyy-MM-dd` {
		return t.Format("2006-01-02")
	} else if format != `` {
		return t.Format(format)
	} else {
		return t.Format("2006-01-02 15:04:05.0000")
	}
}

func tmplInt64ToString(val int64) string {
	return strconv.FormatInt(val, 10)
}

//HTML转义
func tmplUnescaped(x string) interface{} {
	return template.HTML(x)
}

//js转义
func tmplUnescapedjs(x string) interface{} {
	return template.JS(x)
}

//css转义
func tmplUnescapedcss(x string) interface{} {
	return template.CSS(x)
}

//超出部分省略号
func tmplFormatOverflow(x string, l int) string {
	xRune := []rune(x)
	if len(xRune) <= l {
		return x
	} else {
		return string(xRune[:l]) + `...`
	}
}

//格式化pict的url
func tmplFormatPictUrl(val string) string {
	if len(val) > 4 {
		if strings.ToLower(val[:4]) == `http` {
			return val
		} else {
			return `/pict/` + val
		}
	} else {
		return val
	}
}

func tmplMultAdd(args ...int) int {
	sum := 0
	for _, v := range args {
		sum += v
	}
	return sum
}

func tmplMultMinus(args ...int) int {
	sum := 0
	for i, v := range args {
		if i == 0 {
			sum = v
		} else {
			sum -= v
		}
	}
	return sum
}

func tmplMultTimes(args ...int) int {
	sum := 0
	for i, v := range args {
		if i == 0 {
			sum = v
		} else {
			sum *= v
		}
	}
	return sum
}

func tmplMultDivided(args ...int) int {
	sum := 0
	for i, v := range args {
		if i == 0 {
			sum = v
		} else {
			sum /= v
		}
	}
	return sum
}
