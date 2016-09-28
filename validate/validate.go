package validate

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const Tag = "validate"

var ErrParamsCountNotMatch = errors.New("parameters count not matched with validators")

const (
	tagLen    = "Len"
	tagChars  = "Chars"
	tagMatch  = "Match"
	tagChoice = "Choice"
)

type Validator func(string) []error
type Valid struct {
	ValidChain []Validator
	errs       []error
}

func (v *Valid) HasError() bool {
	return len(v.errs) > 0
}

func (v *Valid) Errors() []error {
	return v.errs
}

func (v *Valid) Reset() {
	v.errs = make([]error, 0)
}

func (v *Valid) Validate(s string) []error {
	for _, validator := range v.ValidChain {
		if errs := validator(s); errs != nil {
			v.errs = append(v.errs, errs...)
		}
	}
	return v.Errors()
}

func New(validators ...Validator) *Valid {
	return &Valid{
		ValidChain: validators,
		errs:       make([]error, 0),
	}
}

func Use(vc ...Validator) Validator {
	return New(vc...).Validate
}

func getTag(t *reflect.Type, field string, tagName string) (tagVal string, err error) {
	fieldVal, ok := (*t).FieldByName(field)
	if ok {
		tagVal = fieldVal.Tag.Get(tagName)
	} else {
		err = errors.New("no field named: " + field)
	}
	return
}

func getValidateTag(t *reflect.Type, field string) (string, error) {
	return getTag(t, field, Tag)
}

type Field struct {
	Name string
	Val  interface{}
	Tag  string
}

// get all fields taged with `validate`
func getFields(s interface{}) []*Field {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	fields := make([]*Field, 0)
	numField := t.NumField()
	for i := 0; i < numField; i++ {
		name := t.Field(i).Name
		val := v.FieldByName(name).String()
		tag, _ := getValidateTag(&t, name)
		if tag != "" {
			fields = append(fields, &Field{Name: name, Val: val, Tag: tag})
		}
	}
	return fields
}

func getValidators(tag string) []Validator {
	tags := strings.Split(tag, ";")
	vc := make([]Validator, 0)
	for _, t := range tags {
		tagMap := strings.Split(t, ":")
		tagName, tagVal := tagMap[0], tagMap[1]
		switch tagName {
		case tagLen:
			vals := strings.Split(tagVal, ",")
			min, err := strconv.Atoi(vals[0])
			if err != nil {
				min = 0
			}
			max, err := strconv.Atoi(vals[1])
			if err != nil {
				max = 0
			}
			vc = append(vc, Use(ValidLen(min, max)))
		case tagChars:
			vc = append(vc, Use(ValidChars(tagVal, true)))
		case tagMatch:
			reg, ok := MatchMap[tagVal]
			if !ok {
				reg, _ = regexp.Compile(tagVal)
			}
			vc = append(vc, Use(ValidMatch(reg)))
		case tagChoice:
			vc = append(vc, Use(ValidChoice(strings.Split(tagVal, ","), true)))
		default:
		}
	}
	return vc
}

// `validate:"Len:6,32;Choice:TestVal,TestVal2"`
// `validate:"Match:IP"` or Email/URL/Domain/MAC
func Validate(s interface{}) []error {
	fields := getFields(s)
	errs := make([]error, 0)
	for _, f := range fields {
		validator := Use(getValidators(f.Tag)...)
		err := validator(f.Val.(string))
		errs = append(errs, err...)
	}
	return errs
}
