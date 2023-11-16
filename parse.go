package timeuse

import (
	"regexp"
)

/*
该正则表达式的模式是用于匹配日期时间字符串的
这个正则表达式可以匹配以下格式的日期时间字符串：

	YYYY-MM-DD
	YYYY/MM/DD
	YYYY-MM-DDTHH:MM:SS
	YYYY/MM/DD HH:MM:SS
	YYYY-MM-DDTHH:MM:SS.sss
	YYYY/MM/DD HH:MM:SS.sss
*/
var templateRe = regexp.MustCompile(`^(\d{4})[-\/]?(\d{1,2})?[-\/]?(\d{0,2})[Tt\s]*(\d{1,2})?:?(\d{1,2})?:?(\d{1,2})?[.:]?(\d+)?$`)

/*
该正则表达式用于匹配日期时间格式模板中的各个部分
该正则表达式可以匹配以下格式的字符串：

	[...]：方括号中的内容，例如 [Year]。
	Y{1,4}：1到4个大写字母 Y，表示年份。
	M{1,4}：1到4个大写字母 M，表示月份。
	D{1,2}：1到2个大写字母 D，表示日期。
	d{1,2}：1到2个小写字母 d，表示周几。
	H{1,2}：1到2个大写字母 H，表示小时（24小时制）。
	h{1,2}：1到2个小写字母 h，表示小时（12小时制）。
	m{1,2}：1到2个小写字母 m，表示分钟。
	s{1,2}：1到2个小写字母 s，表示秒。
	Z{1,2}：1到2个大写字母 Z，表示时区偏移。
	SSS：3个大写字母 S，表示毫秒。
*/
var formatRe = regexp.MustCompile(`\[([^\]]+)]|Y{1,4}|M{1,4}|D{1,2}|d{1,2}|H{1,2}|h{1,2}|m{1,2}|s{1,2}|Z{1,2}|SSS`)

///*
//该正则表达式的模式是用于匹配日期时间格式模板中的占位符
//该正则表达式可以匹配以下格式的字符串：
//
//	\[...]：以方括号括起来的内容，例如 [Year]。
//	Y：大写字母 Y，表示年份。
//	M：大写字母 M，表示月份。
//	D：大写字母 D，表示日期。
//	H：大写字母 H，表示小时。
//	m：小写字母 m，表示分钟。
//	s：小写字母 s，表示秒。
//*/
//var duraFormatRe = regexp.MustCompile(`/\[([^\]]+)]|Y|M|D|H|m|s/g`)

func parseT(t string) []string {
	ret := templateRe.FindStringSubmatch(t) // 正则表达式的方法，用于在字符串中查找匹配项，并返回匹配项的子匹配结果

	// 长度正常为8，第一个元素是完整的匹配项，后续的7个元素对应于年、月、日、小时、分钟、秒和毫秒这7个子匹配项。没输入为空
	if len(ret) <= 1 { // 长度为1时：什么情况下可以找到完整的匹配项，却找不到任何其他子匹配项呢？！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！
		return nil
	}

	return ret
}

func parseList(list []int) (year, month, day, hour, minute, second int) {
	l := len(list)
	if l > 0 {
		year = list[0]
	}

	if l > 1 {
		month = list[1]
	} else {
		month = 1
	}

	if l > 2 {
		day = list[2]
	} else {
		day = 1
	}

	if l > 3 {
		hour = list[3]
	}
	if l > 4 {
		minute = list[4]
	}
	if l > 5 {
		second = list[5]
	}

	return year, month, day, hour, minute, second
}
