package main

var chnNumbers = map[string]int{
	"零": 0,
	"一": 1,
	"二": 2,
	"三": 3,
	"四": 4,
	"五": 5,
	"六": 6,
	"七": 7,
	"八": 8,
	"九": 9,
	"十": 10,
}

var chnUnit = map[string]int{
	"十": 10,
	"百": 100,
	"千": 1000,
	"万": 10000,
	"亿": 100000000,
}

func ChnNumberToInt(s string) int {
	total := 0
	num := 0
	for _, c := range s {
		ch := string(c)
		if n, ok := chnNumbers[ch]; ok && num == 0 {
			num = n
			if num == 10 {
				total += 10
				num = 0
			}
			continue
		}
		if n, ok := chnUnit[ch]; ok {
			if ch == "万" || ch == "亿" {
				total = (total + num) * n
			} else {
				total = total + num*n
			}
			num = 0
		}
	}
	total = total + num
	return total
}
