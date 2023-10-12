package main

import "github.com/aynakeya/deepcolor/transform"

type EpisodeExtractor interface {
	Extract(name string) EpisodeInfo
}

type DefaultEpisodeExtractor struct {
	EpTitle  transform.Translator
	EpNum    transform.Translator
	EpRes    transform.Translator
	EpFormat transform.Translator
	EpSeason transform.Translator
	Tag      transform.Translator
	Ext      transform.Translator
}

func (d *DefaultEpisodeExtractor) Extract(name string) EpisodeInfo {
	return EpisodeInfo{
		Title:      d.EpTitle.MustApply(name).(string),
		Episode:    d.EpNum.MustApply(name).(int),
		Resolution: d.EpRes.MustApply(name).(string),
		Format:     d.EpFormat.MustApply(name).(string),
		Season:     d.EpSeason.MustApply(name).(int),
		ExtraTag:   d.Tag.MustApply(name).([]string),
		Extension:  d.Ext.MustApply(name).(string),
	}
}
