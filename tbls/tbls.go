package tbls

import (
	"math"
	"strconv"
)

type Tbls struct {
	PageNo    int `gorm:"-"` // 忽略这个字段  分页时，此字段供外部读取
	PageCount int `gorm:"-"` // 忽略这个字段  分页时，此字段供外部读取
	RowCount  int `gorm:"-"` // 忽略这个字段  分页时，此字段供外部读取
	IgnoreMe  int `gorm:"-"` // 忽略这个字段
}

//计算PageNo、PageCount
func (t *Tbls) Calclulate(pageNo string, pageRow int) {
	//获取当前页码和总页码
	t.PageNo = 1
	if t.RowCount <= 0 {
		t.PageCount = 1
	} else {
		t.PageCount = int(math.Ceil(float64(t.RowCount) / float64(pageRow)))
	}

	if len(pageNo) == 0 {
		pageNo = `1`
	}

	value, err := strconv.Atoi(pageNo)
	if err != nil || value <= 0 {
		t.PageNo = 1
	} else {
		t.PageNo = value
	}
	if t.PageNo > t.PageCount {
		t.PageNo = t.PageCount
	}
}
