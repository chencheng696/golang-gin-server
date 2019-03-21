package main

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func ValidPattern() map[string][2]string {

	pattern := make(map[string][2]string)

	pattern[`t`] = [2]string{``, `文本`}
	pattern[`s`] = [2]string{`^.*$`, `文字列一行`}
	pattern[`w`] = [2]string{`^\w*$`, `英数字`}
	pattern[`d`] = [2]string{`^\d*$`, `数字`}
	pattern[`i`] = [2]string{`^[\+\-]{0,1}[1-9]{1,}\d*$`, `数字`}
	pattern[`f`] = [2]string{`^[\+\-]{0,1}[1-9]{1,}\d*\.?\d*$`, `数字`}
	pattern[`e`] = [2]string{`^([^@\s]+)@((?:[-a-z0-9]+\.)+[a-z]{2,})$`, `邮件店址`}
	pattern[`ss`] = [2]string{`^[^\'\"\\<>]*$`, `「\',",\\,<,>」以外的文字`}
	pattern[`date`] = [2]string{`^[\d\-\/]*$`, `日期`}
	pattern[`time`] = [2]string{`^[\d:]*$`, `时间`}
	pattern[`wd`] = [2]string{`^[a-zA-Z0-9]+$`, `英数字`}
	pattern[`pass`] = [2]string{`^(([0-9]+[a-zA-Z]+)|([a-zA-Z]+[0-9]+))[a-zA-Z0-9]*$`, `半角英数字组合（8位以上、15位一下`}

	return pattern
}

/*
*	checkMap定义
*	key=>变量名
*	value=>map[string]string
*		type:内容正则表达式
*		name:字段名字
*		isempty:是否允许空
*		maxlen:最大长度
*		minlen:最小长度
*
*	客户端提交过来非法字段 则报错
 */
func ValidForm(c *gin.Context, checkMap map[string]interface{}) (e map[string]string, r map[string]interface{}) {
	return ValidForm2(c, checkMap, true)
}

func ValidForm2(c *gin.Context, checkMap map[string]interface{}, isStrict bool) (e map[string]string, r map[string]interface{}) {

	checkMap[`cmd`] = map[string]string{
		`type`: `:ss`,
	}
	checkMap[`pageNo`] = map[string]string{
		`type`: `:d`,
	}
	checkMap[`pageCount`] = map[string]string{
		`type`: `:d`,
	}
	checkMap[`searchNo`] = map[string]string{
		`type`: `:d`,
	}
	checkMap[`selectNo`] = map[string]string{
		`type`: `:d`,
	}
	checkMap[`langCode`] = map[string]string{
		`type`: `:s`,
	}
	checkMap[`langcode`] = map[string]string{
		`type`: `:s`,
	}

	req := c.Request
	//req.ParseForm()
	//req.ParseMultipartForm(c.engine.MaxMultipartMemory)

	GetValues := req.URL.Query()
	PostValues := req.PostForm

	form, _ := c.MultipartForm()

	pattern := ValidPattern()

	e = make(map[string]string)
	r = make(map[string]interface{})

	for k, _ := range checkMap {
		r[k] = ``
	}

	//get
	if GetValues != nil {
		for k, v := range GetValues {
			if len(v) > 1 {
				r[k] = v
			} else {
				r[k] = v[0]
			}
		}
	}
	//post
	if PostValues != nil {
		for k, v := range PostValues {
			if len(v) > 1 {
				r[k] = v
			} else {
				r[k] = v[0]
			}
		}
	}
	//MultipartForm
	if form != nil {
		if form.Value != nil {
			for k, v := range form.Value {
				if len(v) > 1 {
					r[k] = v
				} else {
					r[k] = v[0]
				}
			}
		}
		if form.File != nil {
			for k, v := range form.File {
				if len(v) > 1 {
					r[k] = v
				} else {
					r[k] = v[0]
				}
			}
		}
	}

	if isStrict {
		for k, _ := range r {
			if _, ok := checkMap[k]; !ok {
				e[k] = `[` + k + `]非法字段，请确认！`
			}
		}
		if len(e) > 0 {
			return e, r
		}
	}

	for k, v := range checkMap {

		e[k] = ``

		attr := v.(map[string]string)

		sType := ``
		if value, ok := attr[`type`]; ok {
			sType = value
		}

		sName := ``
		if value, ok := attr[`name`]; ok {
			sName = value
		}
		if len(sName) == 0 {
			sName = k
		}

		isEmpty := `1`
		if len(sType) > 0 {
			if sType[:1] != `:` {
				isEmpty = `0`
			} else {
				sType = sType[1:]
			}
		}

		maxLen := ``
		if value, ok := attr[`maxlen`]; ok {
			maxLen = value
		}

		minLen := ``
		if value, ok := attr[`minlen`]; ok {
			minLen = value
		}

		fixLen := ``
		if value, ok := attr[`fixlen`]; ok {
			fixLen = value
		}

		value, ok := r[k]
		if !ok {
			if isEmpty == `0` {
				e[k] = sName + `不能为空！`
			}
			continue
		}

		subVal := make([]string, 0)
		t := reflect.TypeOf(value)
		if t.Kind() == reflect.String {
			subVal = append(subVal, value.(string))
		} else if t.Kind() == reflect.Slice {
			if reflect.SliceOf(t).Kind() == reflect.String {
				subVal = value.([]string)
			} else {
				continue //文件组
			}
		} else {
			continue //文件
		}

		if len(subVal) == 0 {
			if isEmpty == `0` {
				e[k] = sName + `不能为空！`
			}
			continue
		}

		for _, item := range subVal {
			if isEmpty == `0` && len(item) == 0 {
				e[k] = sName + `不能为空！`
				break
			}

			if len(sType) > 0 {
				p, ok := pattern[sType]
				if !ok {
					e[k] = sName + `校验类型(` + sType + `)不存在！`
					break
				}

				match, _ := regexp.MatchString(p[0], item)
				if !match {
					e[k] = sName + `只能是` + p[1]
					break
				}

				min := -1
				if len(minLen) > 0 {
					m, err := strconv.Atoi(minLen)
					if err != nil {
						e[k] = sName + `的minLen设置不正确`
						break
					} else {
						min = m
					}
				}
				max := -1
				if len(maxLen) > 0 {
					m, err := strconv.Atoi(maxLen)
					if err != nil {
						e[k] = sName + `的maxLen设置不正确`
						break
					} else {
						max = m
					}
				}
				if max != -1 && max < min {
					e[k] = sName + `的maxLen设置不正确`
					break
				}

				if min >= 0 && max >= 0 {
					if len(item) < min || len(item) > max {
						e[k] = sName + `的长度范围是` + minLen + `到` + maxLen + `位`
						break
					}
				} else if min >= 0 {
					if len(item) < min {
						e[k] = sName + `的长度最小是` + minLen + `位`
						break
					}
				} else if max >= 0 {
					if len(item) > max {
						e[k] = sName + `的长度最大是` + maxLen + `位`
						break
					}
				}

				if len(fixLen) > 0 {
					m, err := strconv.Atoi(fixLen)
					if err != nil {
						e[k] = sName + `的fixLen设置不正确`
						break
					} else {
						if len(item) != m {
							e[k] = sName + `的长度必须固定` + fixLen + `位`
							break
						}
					}
				}
			}
		}
	}

	r[`commomMsg`] = ``
	return e, r
}

func ValidUpload(config []map[string]string, list [][]string) error {

	errorMsg := make([]string, 0)

	pattern := ValidPattern()

	for i, v := range list {
		for j, w := range v {
			check := config[j]

			sKey := ``
			if value, ok := check[`key`]; ok {
				sKey = value
			}

			sType := ``
			if value, ok := check[`type`]; ok {
				sType = value
			}

			sName := ``
			if value, ok := check[`name`]; ok {
				sName = value
			}
			if len(sName) == 0 {
				sName = sKey
			}

			isEmpty := `1`
			if len(sType) > 0 {
				if sType[:1] != `:` {
					isEmpty = `0`
				} else {
					sType = sType[1:]
				}
			}

			maxLen := ``
			if value, ok := check[`maxlen`]; ok {
				maxLen = value
			}

			minLen := ``
			if value, ok := check[`minlen`]; ok {
				minLen = value
			}

			fixLen := ``
			if value, ok := check[`fixlen`]; ok {
				fixLen = value
			}

			if len(w) == 0 {
				if isEmpty == `0` {
					errorMsg = append(errorMsg, `第`+strconv.Itoa(i+1)+`行【`+sName+`】不能为空！`)
				}
				continue
			}

			if len(sType) > 0 {
				p, ok := pattern[sType]
				if !ok {
					errorMsg = append(errorMsg, `第`+strconv.Itoa(i+1)+`行【`+sName+`】校验类型(`+sType+`)不存在！`)
					continue
				}

				match, _ := regexp.MatchString(p[0], w)
				if !match {
					errorMsg = append(errorMsg, `第`+strconv.Itoa(i+1)+`行【`+sName+`】只能是`+p[1])
					continue
				}

				min := -1
				if len(minLen) > 0 {
					m, err := strconv.Atoi(minLen)
					if err != nil {
						errorMsg = append(errorMsg, `第`+strconv.Itoa(i+1)+`行【`+sName+`】的minLen设置不正确`)
						continue
					} else {
						min = m
					}
				}
				max := -1
				if len(maxLen) > 0 {
					m, err := strconv.Atoi(maxLen)
					if err != nil {
						errorMsg = append(errorMsg, `第`+strconv.Itoa(i+1)+`行【`+sName+`】的maxLen设置不正确`)
						continue
					} else {
						max = m
					}
				}
				if max != -1 && max < min {
					errorMsg = append(errorMsg, `第`+strconv.Itoa(i+1)+`行【`+sName+`】的maxLen设置不正确`)
					continue
				}

				if min >= 0 && max >= 0 {
					if len(w) < min || len(w) > max {
						errorMsg = append(errorMsg, `第`+strconv.Itoa(i+1)+`行【`+sName+`】的长度范围是`+minLen+`到`+maxLen+`位`)
						continue
					}
				} else if min >= 0 {
					if len(w) < min {
						errorMsg = append(errorMsg, `第`+strconv.Itoa(i+1)+`行【`+sName+`】的长度最小是`+minLen+`位`)
						continue
					}
				} else if max >= 0 {
					if len(w) > max {
						errorMsg = append(errorMsg, `第`+strconv.Itoa(i+1)+`行【`+sName+`】的长度最大是`+maxLen+`位`)
						continue
					}
				}

				if len(fixLen) > 0 {
					m, err := strconv.Atoi(fixLen)
					if err != nil {
						errorMsg = append(errorMsg, `第`+strconv.Itoa(i+1)+`行【`+sName+`】的fixLen设置不正确`)
						continue
					} else {
						if len(w) != m {
							errorMsg = append(errorMsg, `第`+strconv.Itoa(i+1)+`行【`+sName+`】的长度必须固定`+fixLen+`位`)
							continue
						}
					}
				}
			}
		}
	}

	if len(errorMsg) > 0 {
		return errors.New(strings.Join(errorMsg, "\n"))
	} else {
		return nil
	}
}

func CheckValidResult(data map[string]string) bool {
	for _, v := range data {
		if len(v) > 0 {
			return false
		}
	}
	return true
}
