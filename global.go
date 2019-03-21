package main

var (
	//此变量表示额外拥有哪些语言
	gLanguages = map[string]string{
		`en`: `English`,
		`jp`: `日本語`,
	}
)

var (
	gMapPerm = map[string]string{
		`0`: `超级管理员`,
		`1`: `普通管理员`,
	}
	gMapEnable = map[string]string{
		`0`: `禁用`,
		`1`: `启用`,
	}
	gMapItemStatus = map[string]string{
		`0`: `下架`,
		`1`: `正常`,
	}
	gMapHeadCount = map[string]string{
		`0`: `若干`,
		`1`: `1-2人`,
		`2`: `2-5人`,
		`3`: `5人以上`,
	}
	gMapTreatment = map[string]string{
		`0`: `面议`,
		`1`: `2000-5000`,
		`2`: `5000-10000`,
		`3`: `10000-20000`,
	}
	//根据名称调用函数
	gMapFunc = map[string]func(string) string{
		`ConvertPerm`:       ConvertPerm,
		`ConvertStatus`:     ConvertStatus,
		`ConvertItemStatus`: ConvertItemStatus,
	}
)

func ConvertPerm(s string) string {
	if value, ok := gMapPerm[s]; ok {
		return value
	} else {
		return ``
	}
}

func ConvertStatus(s string) string {
	if value, ok := gMapEnable[s]; ok {
		return value
	} else {
		return ``
	}
}

func ConvertItemStatus(s string) string {
	if value, ok := gMapItemStatus[s]; ok {
		return value
	} else {
		return ``
	}
}
