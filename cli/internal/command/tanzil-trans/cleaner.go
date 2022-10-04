package tanzilTrans

import (
	"data-quran-cli/internal/util"
	"regexp"
	"strings"

	"github.com/go-shiori/dom"
)

type fnCleaner func(TranslationData) TranslationData

var cleanerList = map[string]fnCleaner{
	"en-transliteration": cleanEnTransliteration, // remove html tags
	"id-muntakhab":       cleanIdMuntakhab,       // remove surah info
	"ru-muntahab":        cleanRuMuntahab,        // remove surah info
	"zh-jian":            cleanZhJian,            // remove html tags
}

var rxNewline = regexp.MustCompile(`\n{1,}`)

func cleanEnTransliteration(data TranslationData) TranslationData {
	div := dom.CreateElement("div")
	rxAlefLamShams := regexp.MustCompile(`(?i)al([tsrzdn])`)

	for i, ayah := range data.AyahList {
		// Remove HTML tags
		dom.SetInnerHTML(div, ayah.Translation)
		trans := util.DomTextContent(div)

		// Replace alef lam shamsiyyah
		trans = rxAlefLamShams.ReplaceAllString(trans, "a${1}")

		// Replace 'ayn
		trans = strings.ReplaceAll(trans, "AA", "'")

		// Apply the cleaned value
		ayah.Translation = trans
		data.AyahList[i] = ayah
	}

	return data
}

func cleanIdMuntakhab(data TranslationData) TranslationData {
	rxSurahInfo := regexp.MustCompile(`^\[{2}\d+ ~ ([^\]]+)\]{2}\s*`)
	for i, ayah := range data.AyahList {
		trans := ayah.Translation
		trans = rxSurahInfo.ReplaceAllString(trans, "")
		ayah.Translation = trans
		data.AyahList[i] = ayah
	}
	return data
}

func cleanRuMuntahab(data TranslationData) TranslationData {
	rxSurahInfo1 := regexp.MustCompile(`\s*\[{2}Во имя[^\]]+\]{2}`)
	rxSurahInfo2 := regexp.MustCompile(`\s*\[{2}[^\]]+\]{2}`)
	for i, ayah := range data.AyahList {
		trans := ayah.Translation
		if i+1 == 1236 {
			trans = rxSurahInfo2.ReplaceAllString(trans, "")
		} else {
			trans = rxSurahInfo1.ReplaceAllString(trans, "")
		}

		ayah.Translation = trans
		data.AyahList[i] = ayah
	}
	return data
}

func cleanZhJian(data TranslationData) TranslationData {
	div := dom.CreateElement("div")
	for i, ayah := range data.AyahList {
		// Remove HTML tags
		dom.SetInnerHTML(div, ayah.Translation)
		trans := util.DomTextContent(div)
		trans = rxNewline.ReplaceAllString(trans, "\n\n")
		trans = strings.TrimSpace(trans)
		ayah.Translation = trans
		data.AyahList[i] = ayah
	}
	return data
}
