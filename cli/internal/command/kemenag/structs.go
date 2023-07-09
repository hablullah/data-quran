package kemenag

import "gopkg.in/guregu/null.v4"

type Surah struct {
	ID              int    `json:"id"`
	Arabic          string `json:"arabic"`
	Latin           string `json:"latin"`
	Transliteration string `json:"transliteration"`
	Translation     string `json:"translation"`
	NumAyah         int    `json:"num_ayah"`
	Page            int    `json:"page"`
	Location        string `json:"location"`
}

type Tafsir struct {
	Wajiz              string      `json:"wajiz"`
	Tahlili            string      `json:"tahlili"`
	IntroSurah         null.String `json:"intro_surah"`
	OutroSurah         null.String `json:"outro_surah"`
	MunasabahPrevSurah null.String `json:"munasabah_prev_surah"`
	MunasabahPrevTheme null.String `json:"munasabah_prev_theme"`
	ThemeGroup         null.String `json:"theme_group"`
	Kosakata           string      `json:"kosakata"`
	SababNuzul         null.String `json:"sabab_nuzul"`
	Conclusion         null.String `json:"conclusion"`
}

type Ayah struct {
	ID          int         `json:"id"`
	SurahID     int         `json:"surah_id"`
	Ayah        int         `json:"ayah"`
	Page        int         `json:"page"`
	QuarterHizb float64     `json:"quarter_hizb"`
	Juz         int         `json:"juz"`
	Manzil      int         `json:"manzil"`
	Arabic      string      `json:"arabic"`
	Latin       string      `json:"latin"`
	Translation string      `json:"translation"`
	NoFootnote  null.String `json:"no_footnote"`
	Footnotes   null.String `json:"footnotes"`
	Surah       Surah       `json:"surah"`
	Tafsir      Tafsir      `json:"tafsir"`
}

type RespDownloadTafsir struct {
	Data Ayah `json:"data"`
}

type ListSurahEntry struct {
	Name        string `json:"name"`
	Translation string `json:"translation"`
}
