package quranenc

import "regexp"

type fnCleaner func(FlattenedData) FlattenedData

var cleanerList = map[string]fnCleaner{
	"afar_hamza":            cleanAfarHamza,
	"english_hilali_khan":   cleanEnglishHilaliKhan,
	"english_saheeh":        cleanEnglishSaheeh,
	"spanish_montada_eu":    cleanSpanishMontada,
	"spanish_montada_latin": cleanSpanishMontada,
	"french_montada":        cleanFrenchMontada,
	"french_rashid":         cleanFrenchRashid,
	"indonesian_sabiq":      cleanIndonesianSabiq,
	"tamil_baqavi":          cleanTamilBaqavi,
	"uyghur_saleh":          cleanUyghurSaleh,
	"uzbek_mansour":         cleanUzbekMansour,
}

func cleanAfarHamza(data FlattenedData) FlattenedData {
	rxNumber := regexp.MustCompile(`^\d+.\s*`)

	for i, ayah := range data.AyahList {
		ayah.Translation = rxNumber.ReplaceAllString(ayah.Translation, "")
		data.AyahList[i] = ayah
	}

	return data
}

func cleanEnglishHilaliKhan(data FlattenedData) FlattenedData {
	rxNumber := regexp.MustCompile(`^\d+.\s*`)

	for i, ayah := range data.AyahList {
		ayah.Translation = rxNumber.ReplaceAllString(ayah.Translation, "")
		data.AyahList[i] = ayah
	}

	return data
}

func cleanEnglishSaheeh(data FlattenedData) FlattenedData {
	rxNumber := regexp.MustCompile(`^\(\d+\)\s*`)

	for i, ayah := range data.AyahList {
		ayah.Translation = rxNumber.ReplaceAllString(ayah.Translation, "")
		data.AyahList[i] = ayah
	}

	return data
}

func cleanSpanishMontada(data FlattenedData) FlattenedData {
	rxNumber := regexp.MustCompile(`^\d+.\s*`)

	for i, ayah := range data.AyahList {
		ayah.Translation = rxNumber.ReplaceAllString(ayah.Translation, "")
		data.AyahList[i] = ayah
	}

	return data
}

func cleanFrenchMontada(data FlattenedData) FlattenedData {
	rxNumber := regexp.MustCompile(`^\d+.\s*`)

	for i, ayah := range data.AyahList {
		ayah.Translation = rxNumber.ReplaceAllString(ayah.Translation, "")
		data.AyahList[i] = ayah
	}

	return data
}

func cleanFrenchRashid(data FlattenedData) FlattenedData {
	rxNumber := regexp.MustCompile(`^\d+\s*`)

	for i, ayah := range data.AyahList {
		ayah.Translation = rxNumber.ReplaceAllString(ayah.Translation, "")
		data.AyahList[i] = ayah
	}

	return data
}

func cleanIndonesianSabiq(data FlattenedData) FlattenedData {
	rxNumber := regexp.MustCompile(`^\*?\d+\.\s*`)

	for i, ayah := range data.AyahList {
		ayah.Translation = rxNumber.ReplaceAllString(ayah.Translation, "")
		data.AyahList[i] = ayah
	}

	return data
}

func cleanTamilBaqavi(data FlattenedData) FlattenedData {
	rxNumber := regexp.MustCompile(`^[\d\-, ]+\.\s*`)

	for i, ayah := range data.AyahList {
		ayah.Translation = rxNumber.ReplaceAllString(ayah.Translation, "")
		data.AyahList[i] = ayah
	}

	return data
}

func cleanUyghurSaleh(data FlattenedData) FlattenedData {
	rxNumber := regexp.MustCompile(`\[[\d\-ـ ]+\]\.?`)

	for i, ayah := range data.AyahList {
		ayah.Translation = rxNumber.ReplaceAllString(ayah.Translation, "")
		data.AyahList[i] = ayah
	}

	return data
}

func cleanUzbekMansour(data FlattenedData) FlattenedData {
	rxNumber := regexp.MustCompile(`^\d+\.\s*`)

	for i, ayah := range data.AyahList {
		ayah.Translation = rxNumber.ReplaceAllString(ayah.Translation, "")
		data.AyahList[i] = ayah
	}

	return data
}

func noFootnote(data FlattenedData) FlattenedData {
	for i, ayah := range data.AyahList {
		ayah.Footnotes = ""
		data.AyahList[i] = ayah
	}

	return data
}
