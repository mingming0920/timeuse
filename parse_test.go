package timeuse

import (
	"fmt"
	"testing"
)

func TestTemplateRe(t *testing.T) {
	failStr1 := `2021-08-02`
	str2 := `2021-08-02T05:53:12`
	str3 := `21210802T055312`
	f := ErrorF("templateRe")

	if !templateRe.MatchString(failStr1) {
		t.Error(f(failStr1))
	}
	if !templateRe.MatchString(str2) {
		t.Error(f(str2))
	}
	if !templateRe.MatchString(str3) {
		t.Error(f(str3))
	}
	fmt.Println(failStr1)
	fmt.Println(templateRe.MatchString(failStr1))
	fmt.Println(str2)
	fmt.Println(templateRe.MatchString(str2))
	fmt.Println(str3)
	fmt.Println(templateRe.MatchString(str3))
}

func TestParseT(t *testing.T) {
	f := ErrorF("ParseT")

	p1 := "2021"
	r1 := parseT(p1)
	fmt.Println(r1)
	if r1 == nil {
		t.Error(f(p1))
	}
	p2 := "2004-05-03T17:30:08:222"
	r2 := parseT(p2)
	fmt.Println(r2)
	if r2 == nil {
		t.Error(f(p2))
	}

	p3 := "2021-08-03T04:44:03+00:00" //"00:00"代表的是时区的偏移量。类似时区的作用
	r3 := parseT(p3)
	fmt.Println(r3)
	if r3 != nil {
		t.Error(f(p3))
	}
}

func TestParseList(t *testing.T) {
	list := []int{2021, 7, 30, 9, 42, 56}
	y, month, d, h, m, s := parseList(list)
	f := ErrorF("ParseList")

	fmt.Println(parseList(list))
	if y != 2021 || month != 7 || d != 30 || h != 9 || m != 42 || s != 56 {
		t.Error(f(list))
	}

	list2 := []int{2021}
	y, month, d, h, m, s = parseList(list2)

	fmt.Println(parseList(list2))
	if y != 2021 || month != 1 || d != 1 || h != 0 || m != 0 || s != 0 {
		t.Error(f(list2))
	}
}
