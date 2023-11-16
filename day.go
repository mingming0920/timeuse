package timeuse

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"timeuse/locale"
)

type D struct {
	time     time.Time
	Year     int
	Month    time.Month
	Day      int
	Hour     int
	Minute   int
	Second   int
	Unix     int64
	UnixNano int64
	Weekday  time.Weekday
}

type Unit = int

// day构建在单位类型中，用于Add、Subtract、Set、EndOf、StartOf方法
const (
	Year Unit = iota //  iota 枚举器，它使得每个常量的值逐个递增
	Month
	Day
	Hour
	Minute
	Second
	Weekday
)

var translator locale.Translator

func init() {
	translator = locale.EN
}

func (d *D) fields() {
	time := d.time

	d.Year = time.Year()
	d.Month = time.Month()
	d.Day = time.Day()
	d.Hour = time.Hour()
	d.Minute = time.Minute()
	d.Second = time.Second()
	d.Weekday = time.Weekday()

	unixNano := time.UnixNano()
	d.Unix = unixNano / 1e6
	d.UnixNano = unixNano
}

func createDay(time time.Time) *D {
	d := &D{
		time: time,
	}

	d.fields()
	return d
}

// MonthDay 按年度和月份获取最大天数
func MonthDay(year int, month int) int {
	leapMonth := []int{1, 3, 5, 7, 8, 10, 12}
	if month == 2 {
		if IsLeapYear(year) {
			return 29
		} else {
			return 28
		}
	} else if intInSlice(leapMonth, month) {
		return 31
	} else {
		return 30
	}
}

// IsLeapYear 判断是不是闰年
func IsLeapYear(year int) bool {
	return year%400 == 0 || (year%4 == 0 && year%100 != 0)
}

// 判断是不是31天的月份
func intInSlice(nums []int, num int) bool {
	for _, n := range nums { // 判断num是nums中的某一个不
		if n == num {
			return true
		}
	}
	return false
}

// Locale 设置day的翻译，day/locale具有zh-cn或en，默认值：en
func Locale(t locale.Translator) {
	translator = t
}

func (d *D) Time() time.Time {
	return d.time
}

// 用值或单位更改日期时间
// 值可能为int或-int
func (d *D) change(value int, unit Unit) *D {
	sec := int(time.Second)
	switch unit {
	case Year:
		day := MonthDay(d.Year+value, int(d.Month))
		return createDay(d.time.AddDate(value, 0, day-d.Day))
	case Month:
		month := int(d.Month) + value
		// 感觉第二个d.Minute应该是d.Second----------------------！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！
		return createDay(time.Date(d.Year, time.Month(month), MonthDay(d.Year, month), d.Hour, d.Minute, d.Second, d.SecondAfterUnixNano(), d.time.Location()))
	case Day:
		return createDay(d.time.AddDate(0, 0, value))
	case Hour:
		return createDay(d.time.Add(time.Duration(60 * 60 * value * sec)))
	case Minute:
		return createDay(d.time.Add(time.Duration(60 * value * sec)))
	case Second:
		return createDay(d.time.Add(time.Duration(value * sec)))
	}

	return d
}

// SecondAfterUnixNano 9位UnixNano之后的第二个切片
// 示例：1627637214-376669500=>376669500返回UnixNano
func (d *D) SecondAfterUnixNano() int {
	str := fmt.Sprintf("%v", d.UnixNano)
	ret, _ := strconv.Atoi(str[len(str)-9:])

	return ret
}

// New 使用time.Time创建一个新的day.D
func New(time time.Time) *D {
	return createDay(time)
}

// Now 使用time.Now创建一个新的day.D
func Now() *D {
	return createDay(time.Now())
}

// Parse 解析ISO日期时间字符串，使用解析结果创建新的day.D
func Parse(t string) (*D, error) {
	ret := parseT(t)
	if ret == nil {
		return nil, fmt.Errorf("表单失败：无法分析字符串%s", t)
	}

	var list []int
	for _, r := range ret[1:] { // 将 strings 解析为 ints
		val, _ := strconv.Atoi(r)
		list = append(list, val)
	}

	year, month, Day, hour, minute, second := parseList(list)

	return createDay(time.Date(year, time.Month(month), Day, hour, minute, second, 0, time.Local)), nil
}

// Unix 用毫秒创造新的day.D
func Unix(unix int) (*D, error) {
	unixStr := fmt.Sprintf("%v", unix)
	if len(unixStr) != 13 {
		return nil, errors.New("unix是13位毫秒")
	}

	sec := int64(unix / 1e3)

	nsecStr := fmt.Sprintf("%v", unix)
	nsec, _ := strconv.Atoi(nsecStr[len(nsecStr)-3:])

	return createDay(time.Unix(sec, int64(nsec))), nil
}

// List 用int切片创造新的day.D、 未设置的项目使用默认值
// 参数：list int[year, month, day, hour, minute, second]
// 示例：List（[]int｛2021，8，17｝）=>day.D year: 2021, month: 8, day: 17, hour: 0, minute: 0, second: 0
func List(list []int) *D {
	year, month, Day, hour, minute, second := parseList(list)
	return createDay(time.Date(year, time.Month(month), Day, hour, minute, second, 0, time.Local))
}

// Set 按特定单位设定值
// 示例：D.set(2020, day.Year)）//将年份设置为2020
func (d *D) Set(value int, unit Unit) *D {
	switch unit {
	case Year:
		return d.change(value-d.Year, Year)
	case Month:
		return d.change(value-int(d.Month), Month)
	case Day:
		return d.change(value-d.Day, Day) //createDay(d.time.AddDate(0, 0, value-d.Day))
	case Hour:
		return d.change(value-d.Hour, Hour)
	case Minute:
		return d.change(value-d.Minute, Minute)
	case Second:
		return d.change(value-d.Second, Second)
	case Weekday:
		return d.change(value-int(d.Weekday), Day)
	}

	return d
}

// Add 按特定单位增值
// 示例： List([]int{2020}).Add(2, day.Year)）//2020年添加到2022年
func (d *D) Add(value int, uint Unit) *D {
	return d.change(value, uint)
}

// Subtract 按特定单位增值
// 示例：List([]int{2020}).Subtract(2, day.Year) //2020年sub到2018年
func (d *D) Subtract(value int, uint Unit) *D {
	return d.change(-value, uint)
}

func (d *D) SetYear(value int) *D {
	return d.Set(value, Year)
}

func (d *D) SetMonth(value int) *D {
	return d.Set(value, Month)
}

func (d *D) SetDay(value int) *D {
	return d.Set(value, Day)
}

func (d *D) SetMinute(value int) *D {
	return d.Set(value, Minute)
}

func (d *D) SetHour(value int) *D {
	return d.Set(value, Hour)
}

func (d *D) SetSecond(value int) *D {
	return d.Set(value, Second)
}

func (d *D) SetWeekDay(value int) *D {
	return d.Set(value, Weekday)
}

func fillZero(value int) string {
	if value < 10 { // 考虑需不需要加大于等于0！！！
		return fmt.Sprintf("0%v", value)
	} else {
		return fmt.Sprintf("%v", value)
	}
}

func (d *D) Format(t string) string {
	ret := formatRe.ReplaceAllStringFunc(t, func(substr string) string { // ReplaceAllStringFunc 用于使用提供的函数对正则表达式匹配的字符串进行替换操作。
		switch substr {
		case "YYYY":
			return fmt.Sprintf("%v", d.Year)
		case "YY":
			year := fmt.Sprintf("%v", d.Year)
			return year[len(year)-2:]
		case "M":
			return fmt.Sprintf("%v", int(d.Month))
		case "MM":
			return fillZero(int(d.Month))
		case "MMMM":
			return translator.MT(int(d.Month))
		case "D":
			return fmt.Sprintf("%v", d.Day)
		case "DD":
			return fillZero(d.Day)
		case "h":
			return fmt.Sprintf("%v", d.Hour%12) // 12小时制
		case "hh":
			return fillZero(d.Hour % 12)
		case "H":
			return fmt.Sprintf("%v", d.Hour) // 24小时制
		case "HH":
			return fillZero(d.Hour)
		case "m":
			return fmt.Sprintf("%v", d.Minute)
		case "mm":
			return fillZero(d.Minute)
		case "s":
			return fmt.Sprintf("%v", d.Second)
		case "ss":
			return fillZero(d.Second)
		case "SSS":
			unixStr := fmt.Sprint(d.Unix)
			return unixStr[len(unixStr)-3:]
		case "d":
			return fmt.Sprintf("%v", int(d.Weekday))
		case "dd":
			return translator.WT(int(d.Weekday))
		}

		return substr
	})

	return ret
}

// UTC 返回当前day.D, 从UTC时间开始
func (d *D) UTC() *D {
	return createDay(d.time.UTC())
}

// Local 返回当前day.D, 从当地时间开始
func (d *D) Local() *D {
	return createDay(d.time.Local())
}

// StartOf 设置特殊单位的日期开始时间
func (d *D) StartOf(unit Unit) *D {
	year := d.Year
	month := 1
	day := 1

	var (
		hour,
		minute,
		second int
	)

	fmt.Println(unit)
	if unit >= Month {
		month = int(d.Month)
	}
	if unit >= Day {
		day = d.Day
	}
	if unit >= Hour {
		hour = d.Hour
	}
	if unit >= Minute {
		minute = d.Minute
	}
	if unit >= Second {
		second = d.Second
	}

	return createDay(time.Date(year, time.Month(month), day, hour, minute, second, 0, d.time.Location()))
}

// EndOf 设置特殊单位的日期结束时间
func (d *D) EndOf(unit Unit) *D {
	month := 12
	day := MonthDay(d.Year, month)
	hour := 23
	minute := 59
	second := 59

	if unit >= Month {
		month = int(d.Month)
	}
	if unit >= Day {
		day = d.Day
	}
	if unit >= Hour {
		hour = d.Hour
	}
	if unit >= Minute {
		minute = d.Minute
	}
	if unit >= Second {
		second = d.Second
	}

	return createDay(time.Date(d.Year, time.Month(month), day, hour, minute, second, 999999999, d.time.Location()))
}

// DaysInMonth 每月返回天数
func (d *D) DaysInMonth() int {
	return MonthDay(d.Year, int(d.Month))
}

// From a.从（b）返回一个time.Duration，a.time减去b.time
func (d *D) From(d2 *D) time.Duration {
	return d.time.Sub(d2.time)
}
