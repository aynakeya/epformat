package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

type EpisodeInfo struct {
	Title      string
	Season     int
	Episode    int
	Resolution string
	Format     string
	ExtraTag   []string
	Extension  string
}

var templateMaker = template.New("epformat").Funcs(
	template.FuncMap{},
)

func (e EpisodeInfo) FormatInfo(format string) (string, error) {
	formatter, err := templateMaker.Parse(format)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = formatter.Execute(&buf, e)
	return buf.String() + e.Extension, err
}

var MainExtractor EpisodeExtractor = &DefaultEpisodeExtractor{
	EpTitle:  RemoveTagTranslator,
	EpNum:    EpNumTranslator,
	EpRes:    EpResTranslator,
	EpFormat: EpFormatTranslator,
	EpSeason: EpSeasonTranslator,
	Tag:      TagTranslator,
	Ext:      ExtTranslator,
}

var DefaultFormat = `{{.Title}} S{{printf "%02d" .Season}} E{{printf "%02d" .Episode}}`

func main() {
	if err := createRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
