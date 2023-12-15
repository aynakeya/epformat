package main

import (
	"errors"
	"github.com/aynakeya/deepcolor/transform"
	"github.com/aynakeya/deepcolor/transform/translators"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	transform.RegisterTranslator(NewNumToInt())
}

type NumToInt struct {
	transform.BaseTranslator
}

func NewNumToInt() transform.Translator {
	t := &NumToInt{
		BaseTranslator: transform.BaseTranslator{
			Type: "NumToInt",
		},
	}
	return t
}

func (c *NumToInt) Apply(value interface{}) (interface{}, error) {
	s, ok := value.(string)
	num := 0
	if !ok {
		return 0, errors.New("not a number string")
	}
	if n, err := strconv.Atoi(s); err == nil {
		num = n
	} else {
		num = ChnNumberToInt(s)
	}
	return num, nil
}

func (c *NumToInt) MustApply(value interface{}) interface{} {
	v, _ := c.Apply(value)
	return v
}

var stringTrim = transform.WrapTranslator("trimspace", func(value interface{}) (interface{}, error) {
	return strings.TrimSpace(value.(string)), nil
})

// ((\d+)(\.\d+)?)(v\d+)?
var EpNumTranslator = translators.NewPipeline(
	translators.NewSwitcher(
		translators.NewRegExpFindFirst(regexp.MustCompile(`(?i)第((\d+)(\.\d+)?)(v\d+)?集`), 1),
		translators.NewRegExpFindFirst(regexp.MustCompile(`(?i)e((\d+)(\.\d+)?)(v\d+)?`), 1),
		translators.NewRegExpFindFirst(regexp.MustCompile(`(?i)ep((\d+)(\.\d+)?)(v\d+)?`), 1),
		translators.NewRegExpFindFirst(regexp.MustCompile(`(?i)第(.){1,2}集`), 1),
		translators.NewRegExpFindFirst(regexp.MustCompile(`(?i)[\[ ]((\d+)(\.\d+)?)(v\d+)?[\] ]`), 1),
		translators.NewValue("1"),
	),
	NewNumToInt(),
)

var EpSeasonTranslator = translators.NewPipeline(
	translators.NewSwitcher(
		translators.NewRegExpFindFirst(regexp.MustCompile(`(?i)第(\d+)季`), 1),
		translators.NewRegExpFindFirst(regexp.MustCompile(`(?i)s(\d+)`), 1),
		translators.NewRegExpFindFirst(regexp.MustCompile(`(?i)season(\d+)`), 1),
		translators.NewRegExpFindFirst(regexp.MustCompile(`(?i)第(.){1,2}季`), 1),
		translators.NewRegExpFindFirst(regexp.MustCompile(`(?i) (\d)(rd|nd|th) season`), 1),
		translators.NewValue("1"),
	),
	NewNumToInt(),
)

var EpResTranslator = translators.NewPipeline(
	translators.NewRegExpFindFirst(
		regexp.MustCompile(`(?i)\d{3,4}p\+?|\d{3,4}x\d{3,4}|4k|2k|8k`),
		0),
	translators.NewStrCase(true),
)

var EpFormatTranslator = translators.NewPipeline(
	translators.NewRegExpFindFirst(
		regexp.MustCompile(`(?i)mp4|mkv|rmvb`),
		0),
	translators.NewStrCase(true),
	stringTrim,
)

var TagTranslator = translators.NewPipeline(
	translators.NewRegExpReplacer(regexp.MustCompile("【(.*)】"), "[$1]"),
	translators.NewRegExpFindAll(
		regexp.MustCompile(`\[([^\[\]]*)]`), 1),
)

var ExtTranslator = translators.NewSwitcher(
	translators.NewRegExpFindFirst(regexp.MustCompile(`\.\w+$`), 0),
	translators.NewValue(""),
)

var RemoveTagTranslator = translators.NewPipeline(
	translators.NewRegExpReplacer(regexp.MustCompile(`\[[^\[\]]*\]|【[^【】]*】`), ""),
	stringTrim)
