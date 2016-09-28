package validate

import (
	"fmt"
	"regexp"
	"strings"
)

var matchReg = map[string]string{
	"Email":  `\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,4}`,
	"IP":     `^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$`,
	"Domain": `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}$`,
	"URL":    `((http|ftp|https)://)(([a-zA-Z0-9\._-]+\.[a-zA-Z]{2,6})|([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}))(:[0-9]{1,4})*(/[a-zA-Z0-9\&%_\./-~-]*)?`,
	"MAC":    `^([0-9a-fA-F][0-9a-fA-F]:){5}([0-9a-fA-F][0-9a-fA-F])$`,
}

var MatchMap = genMatchMap(matchReg)

func genMatchMap(m map[string]string) (matchMap map[string]*regexp.Regexp) {
	for k, v := range m {
		compile, _ := regexp.Compile(v)
		matchMap[k] = compile
	}
	return
}

type Len struct {
	Min, Max int
}

func (l Len) Validate(s string) []error {
	errs := make([]error, 0)
	if l.Min >= l.Max {
		errs = append(errs, fmt.Errorf("value min(%d) can't >= max(%d)", l.Min, l.Max))
	}
	size := len(strings.TrimSpace(s))
	if l.Min > 0 && size < l.Min {
		errs = append(errs, fmt.Errorf("The length of `%s` can't less than %d", s, l.Min))
	} else if l.Max > 0 && l.Max < size {
		errs = append(errs, fmt.Errorf("The length of `%s` can't greater than %d", s, l.Max))
	}
	return errs
}

func ValidLen(min, max int) Validator {
	return Len{Min: min, Max: max}.Validate
}

type Chars struct {
	Chars      string
	IgnoreCase bool
}

func (c Chars) Validate(s string) []error {
	errs := make([]error, 0)
	if c.IgnoreCase {
		s = strings.ToLower(s)
		c.Chars = strings.ToLower(c.Chars)
	}
	size := len(s)
	for i := 0; i < size; i++ {
		if !strings.Contains(c.Chars, string(s[i])) {
			errs = append(errs, fmt.Errorf("Chars `%s` do not contains `%s[%d]->%s`", c.Chars, s, i, string(s[i])))
		}
	}
	return errs
}

func ValidChars(chars string, ignoreCase bool) Validator {
	return Chars{Chars: chars, IgnoreCase: ignoreCase}.Validate
}

type Match struct {
	Regexp *regexp.Regexp
}

func (m Match) Validate(s string) []error {
	if !m.Regexp.MatchString(s) {
		return []error{fmt.Errorf("%s Do not match %s", s, m.Regexp.String())}
	}
	return nil
}

func ValidMatch(r *regexp.Regexp) Validator {
	return Match{Regexp: r}.Validate
}

type Choice struct {
	Choice     []string
	IgnoreCase bool
}

func (c Choice) Validate(s string) []error {
	if c.IgnoreCase {
		s = strings.ToLower(s)

	}
	for _, cc := range c.Choice {
		if c.IgnoreCase {
			cc = strings.ToLower(cc)
		}
		if cc == s {
			return nil
		}
	}
	return []error{fmt.Errorf("%s not in %v", s, c.Choice)}
}

func ValidChoice(cc []string, ignoreCase bool) Validator {
	return Choice{Choice: cc, IgnoreCase: ignoreCase}.Validate
}
