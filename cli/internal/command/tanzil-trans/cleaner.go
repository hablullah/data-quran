package tanzilTrans

import "regexp"

type fnCleaner func(TranslationData) TranslationData

var cleanerList = map[string]fnCleaner{
	"id-muntakhab": cleanIdMuntakhab, // remove surah info
	"ru-muntahab":  cleanRuMuntahab,  // remove surah info
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
