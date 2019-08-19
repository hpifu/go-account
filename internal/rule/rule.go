package rule

import (
	"fmt"
	"regexp"
	"time"
)

var EmailRegex *regexp.Regexp
var PhoneRegex *regexp.Regexp
var CodeRegex *regexp.Regexp

func init() {
	PhoneRegex = regexp.MustCompile(`^1[345789][0-9]{9}$`)
	EmailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	CodeRegex = regexp.MustCompile(`^[0-9]{6}$`)
}

type Rule func(interface{}) error

func Check(vrules map[interface{}][]Rule) error {
	for val, rules := range vrules {
		for _, r := range rules {
			if err := r(val); err != nil {
				return fmt.Errorf("[%v] %v", val, err)
			}
		}
	}

	return nil
}

func In(sets map[interface{}]struct{}) Rule {
	return func(v interface{}) error {
		if _, ok := sets[v]; !ok {
			return fmt.Errorf("%v 必须在 %v 中", v, sets)
		}
		return nil
	}
}

func Required(v interface{}) error {
	if len(v.(string)) == 0 {
		return fmt.Errorf("必要字段")
	}

	return nil
}

func AtLeast8Characters(v interface{}) error {
	if len(v.(string)) < 8 {
		return fmt.Errorf("至少8个字符")
	}

	return nil
}

func AtMost64Characters(v interface{}) error {
	if len(v.(string)) > 64 {
		return fmt.Errorf("至多64个字符")
	}

	return nil
}

func AtMost32Characters(v interface{}) error {
	if len(v.(string)) > 32 {
		return fmt.Errorf("至多32个字符")
	}

	return nil
}

func ValidEmail(v interface{}) error {
	if !EmailRegex.MatchString(v.(string)) {
		return fmt.Errorf("无效的邮箱")
	}

	return nil
}

func ValidPhone(v interface{}) error {
	if !PhoneRegex.MatchString(v.(string)) {
		return fmt.Errorf("无效的电话号码")
	}

	return nil
}

func ValidCode(v interface{}) error {
	if !CodeRegex.MatchString(v.(string)) {
		return fmt.Errorf("无效的验证码")
	}

	return nil
}

func ValidBirthday(v interface{}) error {
	birthday, err := time.Parse("2006-01-02", v.(string))
	if err != nil {
		return fmt.Errorf("日期格式错误")
	}
	if birthday.After(time.Now()) {
		return fmt.Errorf("日期超过范围")
	}
	if time.Now().Sub(birthday) > 100*365*24*time.Hour {
		return fmt.Errorf("日期超过范围")
	}

	return nil
}
