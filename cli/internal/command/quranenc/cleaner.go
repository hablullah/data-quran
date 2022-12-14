package quranenc

import (
	"data-quran-cli/internal/util"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/zyedidia/generic/mapset"
)

type fnCleaner func(FlattenedData) FlattenedData

var cleanerList = map[string]fnCleaner{
	"afar_hamza":            cleanAfarHamza,
	"assamese_rafeeq":       cleanAssameseRafeeq,
	"bosnian_rwwad":         cleanBosnianRwwad,
	"english_hilali_khan":   cleanEnglishHilaliKhan,
	"english_saheeh":        cleanEnglishSaheeh,
	"spanish_garcia":        cleanSpanishGarcia,
	"spanish_montada_eu":    cleanMontada,
	"spanish_montada_latin": cleanMontada,
	"french_montada":        cleanMontada,
	"french_hameedullah":    cleanFrenchHameedullah,
	"french_rashid":         cleanFrenchRashid,
	"hausa_gummi":           cleanHausaGummi,
	"hindi_omari":           cleanHindiOmari,
	"indonesian_affairs":    cleanIndonesianAffairs,
	"indonesian_complex":    cleanIndonesianComplex,
	"indonesian_sabiq":      cleanIndonesianSabiq,
	"japanese_saeedsato":    cleanJapaneseSaeedsato,
	"malayalam_kunhi":       cleanMalayalamKunhi,
	"oromo_ababor":          cleanOromoAbabor,
	"kinyarwanda_assoc":     cleanKinyarwandaAssoc,
	"albanian_nahi":         cleanAlbanianNahi,
	"swahili_barawani":      cleanSwahiliBarawani,
	"tagalog_rwwad":         cleanTagalogRwwad,
	"urdu_junagarhi":        cleanUrduJunagarhi,
	"tamil_baqavi":          cleanTamilBaqavi,
	"uyghur_saleh":          cleanUyghurSaleh,
	"uzbek_mansour":         cleanUzbekMansour,
	"vietnamese_rwwad":      cleanVietnameseRwwad,
	"yoruba_mikail":         cleanYorubaMikail,
	"gujarati_omari":        cleanGujaratiOmari,
	"somali_yacob":          cleanSomaliYacob,
	"tamil_omar":            cleanTamilOmar,
	"lingala_zakaria":       cleanLingalaZakaria,
}

var rxNewline = regexp.MustCompile(`\n{1,}`)
var rxNewlines = regexp.MustCompile(`\n{2,}`)
var rxFootnoteNumberSplitter = regexp.MustCompile(`(\[\^\d+\]:\s*)`)

func cleanAfarHamza(data FlattenedData) FlattenedData {
	rxAyahNumber := regexp.MustCompile(`^\d+\\?.\s*`)
	data = removeAyahNumber(data, rxAyahNumber)
	return data
}

func cleanAssameseRafeeq(data FlattenedData) FlattenedData {
	data = normalizeFootnoteNumber(data, nil, nil)
	return data
}

func cleanBosnianRwwad(data FlattenedData) FlattenedData {
	data = normalizeFootnoteNumber(data, nil, nil)
	return data
}

func cleanEnglishHilaliKhan(data FlattenedData) FlattenedData {
	// Normalize data
	rxAyahNumber := regexp.MustCompile(`^\d+\\?.\s*`)
	rxFootFn := regexp.MustCompile(`^\\\[(\d+)\\\]\s*`)
	rxTransFn := regexp.MustCompile(`\s*\\\[(\d+)\\\](\s*)`)
	data = removeAyahNumber(data, rxAyahNumber)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)

	// Remove ayah number in footnote
	rxFnCleaner := regexp.MustCompile(`(?i)(\[\^(?:\d+|\*)\]:)(?:\s|\n)*`)
	rxFnAyahNumber := regexp.MustCompile(`(?i)(\[\^(?:\d+|\*)\]:)\s*\(V\.\d+:\d+\)[.,:\-]?(\s*)`)
	for i, ayah := range data.AyahList {
		fn := ayah.Footnotes
		fn = rxFnAyahNumber.ReplaceAllString(fn, "$1$2")
		fn = rxFnCleaner.ReplaceAllString(fn, "$1 ")
		data.AyahList[i].Footnotes = fn
	}

	return data
}

func cleanEnglishSaheeh(data FlattenedData) FlattenedData {
	rxAyahNumber := regexp.MustCompile(`^\(\d+\)\s*`)
	rxFootFn := regexp.MustCompile(`^\\\[(\d+)\\\]\s*-?(\s*)`)
	rxTransFn := regexp.MustCompile(`\\\[(\d+)\\\]\s*-?(\s*)`)
	data = removeAyahNumber(data, rxAyahNumber)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	return data
}

func cleanSpanishGarcia(data FlattenedData) FlattenedData {
	rxFnKeyword := regexp.MustCompile(`^([^:]+):`)

	for i, ayah := range data.AyahList {
		translation := ayah.Translation

		// Process footnote first by splitting it line by line
		var footnoteLines []string
		var footnoteNumber int

		for _, line := range strings.Split(ayah.Footnotes, "\n") {
			// Make sure this line not empty
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Create footnote number
			footnoteNumber++
			strNumber := fmt.Sprintf("[^%d]", footnoteNumber)

			// Check if keyword detected in this line
			if rxFnKeyword.MatchString(line) {
				// Extract keyword
				keyword := rxFnKeyword.FindStringSubmatch(line)[1]
				keyword = strings.TrimSpace(keyword)

				// Remove keyword from the line
				line = rxFnKeyword.ReplaceAllString(line, "")
				line = strings.TrimSpace(line)

				// If keyword exist in translation, put the marker
				if strings.Contains(translation, keyword) {
					translation = strings.Replace(translation, keyword, keyword+strNumber, 1)
					footnoteLines = append(footnoteLines, strNumber+": "+line)
					continue
				}
			}

			// If keyword not found, just put the footnote numbers at the end of translation
			translation += strNumber
			footnoteLines = append(footnoteLines, strNumber+": "+line)
		}

		// Apply normalized data
		footnotes := strings.Join(footnoteLines, "\n\n")
		footnotes = strings.TrimSpace(footnotes)

		ayah.Footnotes = footnotes
		ayah.Translation = translation
		data.AyahList[i] = ayah
	}

	return data
}

func cleanMontada(data FlattenedData) FlattenedData {
	rxAyahNumber := regexp.MustCompile(`^\d+\\?\.\s*`)
	rxTransFn := regexp.MustCompile(`\\\[(\d+)\\\](\s*)`)
	rxFootFn := regexp.MustCompile(`^\\\[(\d+)\\\](\s*)`)
	data = removeAyahNumber(data, rxAyahNumber)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	return data
}

func cleanFrenchHameedullah(data FlattenedData) FlattenedData {
	// Remove HTML tags
	div := dom.CreateElement("div")
	for i, ayah := range data.AyahList {
		dom.SetInnerHTML(div, ayah.Footnotes)
		fns := util.DomTextContent(div)
		fns = rxNewline.ReplaceAllString(fns, "\n\n")
		fns = strings.TrimSpace(fns)
		ayah.Footnotes = fns
		data.AyahList[i] = ayah
	}

	// Normalize footnote
	rxTransFn := regexp.MustCompile(`\\\[(\d+)\\\](\s*)`)
	rxFootFn := regexp.MustCompile(`^\\\[(\d+)\\\](\s*)`)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)

	return data
}

func cleanFrenchRashid(data FlattenedData) FlattenedData {
	rxAyahNumber := regexp.MustCompile(`^\d+\s*`)
	rxTransFn := regexp.MustCompile(`\\\[(\d+)\\\](\s*)`)
	rxFootFn := regexp.MustCompile(`^\\\[(\d+)\\\](\s*)`)
	data = removeAyahNumber(data, rxAyahNumber)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	return data
}

func cleanHausaGummi(data FlattenedData) FlattenedData {
	// Convert stars to marker number
	rxStars := regexp.MustCompile(`(\\*\*)+`)
	for i, ayah := range data.AyahList {
		ayah.Translation = rxStars.ReplaceAllStringFunc(ayah.Translation, func(s string) string {
			nStar := strings.Count(s, "*")
			return fmt.Sprintf("[%d]", nStar)
		})

		ayah.Footnotes = rxStars.ReplaceAllStringFunc(ayah.Footnotes, func(s string) string {
			nStar := strings.Count(s, "*")
			return fmt.Sprintf("\n[%d]", nStar)
		})

		data.AyahList[i] = ayah
	}

	// Normalize number
	rxTransFn := regexp.MustCompile(`\[(\d+)\](\s*)`)
	rxFootFn := regexp.MustCompile(`^\[(\d+)\](\s*)`)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	return data
}

func cleanHindiOmari(data FlattenedData) FlattenedData {
	rxTransFn := regexp.MustCompile(`\\\[(\d+)\\\](\s*)`)
	rxFootFn := regexp.MustCompile(`(\d+)\\*\.(\s*)`)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	data = splitFootnotesByNumber(data)
	return data
}

func cleanIndonesianAffairs(data FlattenedData) FlattenedData {
	rxTransFn := regexp.MustCompile(`(\d+)\s*\)(\s*)`)
	rxFootFn := regexp.MustCompile(`(?:\\*\*|^|\\*\.\s*)(\d+)\s*\)(\s*)`)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	data = splitFootnotesByNumber(data)
	return data
}

func cleanIndonesianComplex(data FlattenedData) FlattenedData {
	rxTransFn := regexp.MustCompile(`(\d+)(\s*)`)
	rxFootFn := regexp.MustCompile(`(?:^|\\*\.\s+)(\d+)\s*\\*\.?(\s*)`)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	return data
}

func cleanIndonesianSabiq(data FlattenedData) FlattenedData {
	rxAyahNumber := regexp.MustCompile(`^\*?\d+\\*\.\s*`)
	rxTransFn := regexp.MustCompile(`\\*\*{1,}\((\d+)\)(\s*)`)
	rxFootFn := regexp.MustCompile(`\\*\*{1,}(\d+)\)\s*\\*\.?(\s*)`)

	data = removeAyahNumber(data, rxAyahNumber)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	data = splitFootnotesByNumber(data)
	return data
}

func cleanJapaneseSaeedsato(data FlattenedData) FlattenedData {
	rxFn := regexp.MustCompile(`(\d+)(\s*)`)
	for i, ayah := range data.AyahList {
		// Extract fn numbers in translation
		transNumbers := mapset.New[string]()
		matches := rxFn.FindAllStringSubmatch(ayah.Translation, -1)
		for _, m := range matches {
			transNumbers.Put(m[1])
		}

		// Replace fn numbers in footnotes
		footNumbers := mapset.New[string]()
		ayah.Footnotes = rxFn.ReplaceAllStringFunc(ayah.Footnotes, func(s string) string {
			parts := rxFn.FindStringSubmatch(s)
			if transNumbers.Has(parts[1]) && !footNumbers.Has(parts[1]) {
				footNumbers.Put(parts[1])
				return fmt.Sprintf("\n[%s]%s", parts[1], parts[2])
			} else {
				return s
			}
		})

		// Replace fn numbers in transation
		ayah.Translation = rxFn.ReplaceAllStringFunc(ayah.Translation, func(s string) string {
			parts := rxFn.FindStringSubmatch(s)
			if footNumbers.Has(parts[1]) {
				return fmt.Sprintf("[%s]%s", parts[1], parts[2])
			} else {
				return s
			}
		})

		data.AyahList[i] = ayah
	}

	// Normalize number
	rxTransFn := regexp.MustCompile(`\[(\d+)\](\s*)`)
	rxFootFn := regexp.MustCompile(`^\[(\d+)\](\s*)`)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	return data
}

func cleanMalayalamKunhi(data FlattenedData) FlattenedData {
	rxTransFn := regexp.MustCompile(`\((\d+)\)(\s*)`)
	rxFootFn := regexp.MustCompile(`^(\d+)\s*\)?\\*\.?(\s*)`)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	return data
}

func cleanOromoAbabor(data FlattenedData) FlattenedData {
	rxTransFn := regexp.MustCompile(`\\\[\\\](\s*)`)
	rxFootFn := regexp.MustCompile(`^\\\[(\d+)\\\](\s*)`)

	for i, ayah := range data.AyahList {
		// Remove number from the footnote
		footnotes := rxFootFn.ReplaceAllString(ayah.Footnotes, "")
		footnotes = strings.TrimSpace(footnotes)
		if footnotes != "" {
			footnotes = "[^1]: " + footnotes
		}

		// If there are no footnotes, remove marker in translation.
		// If there are, fix or add the marker.
		translation := ayah.Translation
		if footnotes == "" {
			translation = rxTransFn.ReplaceAllString(translation, "$1")
		} else {
			if rxTransFn.MatchString(translation) {
				translation = rxTransFn.ReplaceAllString(translation, "[^1]$1")
			} else {
				translation += "[^1]"
			}
		}

		// Apply normalized data
		ayah.Footnotes = footnotes
		ayah.Translation = translation
		data.AyahList[i] = ayah
	}

	return data
}

func cleanKinyarwandaAssoc(data FlattenedData) FlattenedData {
	rxTransFn := regexp.MustCompile(`\\\[(\d+)\\\](\s*)`)
	rxFootFn := regexp.MustCompile(`\\\[(\d+)\\\](\s*)`)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	data = splitFootnotesByNumber(data)
	return data
}

func cleanAlbanianNahi(data FlattenedData) FlattenedData {
	rxTransFn := regexp.MustCompile(`\\\[(\d+)\\\](\s*)`)
	rxFootFn := regexp.MustCompile(`\\\[(\d+)\\\](\s*)`)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	return data
}

func cleanSwahiliBarawani(data FlattenedData) FlattenedData {
	data = normalizeFootnoteNumber(data, nil, nil)
	return data
}

func cleanTagalogRwwad(data FlattenedData) FlattenedData {
	data = normalizeFootnoteNumber(data, nil, nil)
	return data
}

func cleanUrduJunagarhi(data FlattenedData) FlattenedData {
	rxTransFn := regexp.MustCompile(`\((\d+)\)(\s*)`)
	rxFootFn := regexp.MustCompile(`^\((\d+)\)(\s*)`)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	return data
}

func cleanTamilBaqavi(data FlattenedData) FlattenedData {
	rxAyahNumber := regexp.MustCompile(`^[\d\-, ]+\\*\.\s*`)
	data = removeAyahNumber(data, rxAyahNumber)
	return data
}

func cleanUyghurSaleh(data FlattenedData) FlattenedData {
	rxAyahNumber := regexp.MustCompile(`\\\[[\d\-ـ ]+\\\]\.?`)
	data = removeAyahNumber(data, rxAyahNumber)
	return data
}

func cleanUzbekMansour(data FlattenedData) FlattenedData {
	// Remove ayah number
	rxAyahNumber := regexp.MustCompile(`^\d+\\*\.\s*`)
	data = removeAyahNumber(data, rxAyahNumber)

	// Add footnote marker
	rxFootFn := regexp.MustCompile(`^\s*И з о ҳ.\s*(\s*)`)
	for i, ayah := range data.AyahList {
		if ayah.Footnotes != "" {
			ayah.Translation += "[^*]"
			ayah.Footnotes = rxFootFn.ReplaceAllString(ayah.Footnotes, "[^*]: ")
			data.AyahList[i] = ayah
		}
	}

	return data
}

func cleanVietnameseRwwad(data FlattenedData) FlattenedData {
	rxAyahNumber := regexp.MustCompile(`#\\*_\d+\\*_\d+\s*`)
	for i, ayah := range data.AyahList {
		ayah.Footnotes = rxAyahNumber.ReplaceAllString(ayah.Footnotes, "")
		data.AyahList[i] = ayah
	}

	rxTransFn := regexp.MustCompile(`\((\d+)\)(\s*)`)
	rxFootFn := regexp.MustCompile(`^\((\d+)\)(\s*)`)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	return data
}

func cleanYorubaMikail(data FlattenedData) FlattenedData {
	rxTransFn := regexp.MustCompile(`(\d+)(\s*)`)
	rxFootFn := regexp.MustCompile(`^(\d+)\s*\\*\.?(\s*)`)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	return data
}

func cleanGujaratiOmari(data FlattenedData) FlattenedData {
	data = normalizeFootnoteNumber(data, nil, nil)
	return data
}

func cleanSomaliYacob(data FlattenedData) FlattenedData {
	rxAyahNumber := regexp.MustCompile(`^\d+\\*\.?\s*`)
	rxTransFn := regexp.MustCompile(`(\d+)(\s*)`)
	rxFootFn := regexp.MustCompile(`^(\d+)\s*\\*\.?(\s*)`)
	data = removeAyahNumber(data, rxAyahNumber)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	return data
}

func cleanTamilOmar(data FlattenedData) FlattenedData {
	rxTransFn := regexp.MustCompile(`\\+\*+`)
	rxFootFn := regexp.MustCompile(`^\\+\*+\s*`)

	for i, ayah := range data.AyahList {
		ayah.Translation = rxTransFn.ReplaceAllString(ayah.Translation, "[^*]")
		ayah.Footnotes = rxFootFn.ReplaceAllString(ayah.Footnotes, "[^*]: ")
		data.AyahList[i] = ayah
	}

	return data
}

func cleanLingalaZakaria(data FlattenedData) FlattenedData {
	rxAyahNumber := regexp.MustCompile(`^\d+\\*\.\s*`)
	rxTransFn := regexp.MustCompile(`\((\d+)\)(\s*)`)
	rxFootFn := regexp.MustCompile(`^(\d+)(\s*)`)
	data = removeAyahNumber(data, rxAyahNumber)
	data = normalizeFootnoteNumber(data, rxTransFn, rxFootFn)
	return data
}

func removeAyahNumber(data FlattenedData, rxAyahNumber *regexp.Regexp) FlattenedData {
	if rxAyahNumber == nil {
		return data
	}

	for i, ayah := range data.AyahList {
		ayah.Translation = rxAyahNumber.ReplaceAllString(ayah.Translation, "")
		data.AyahList[i] = ayah
	}

	return data
}

func normalizeFootnoteNumber(data FlattenedData, rxTransFn, rxFootFn *regexp.Regexp) FlattenedData {
	for i, ayah := range data.AyahList {
		// Process footnote first
		// Here we split footnotes line by line and extract the numbers
		var footnoteLines []string
		footnoteNumbers := mapset.New[string]()

		for _, line := range strings.Split(ayah.Footnotes, "\n") {
			// Make sure this line not empty
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// If this line contain footnote number, we will normalize it into `[^%d]:` form.
			// If not and there is already a footnote number that found before, we assume this
			// line as continuation of the previous footnote lines.
			if rxFootFn != nil {
				if rxFootFn.MatchString(line) {
					line = rxFootFn.ReplaceAllStringFunc(line, func(match string) string {
						parts := rxFootFn.FindStringSubmatch(match)
						fnNumber := parts[1]

						footnoteNumbers.Put(fnNumber)
						return fmt.Sprintf("[^%s]: ", fnNumber)
					})
				} else if footnoteNumbers.Size() > 0 {
					line = "    " + line
				}
			}

			// Save this line to footnotes
			footnoteLines = append(footnoteLines, line)
		}

		// Merge all footnote lines into one.
		footnotes := strings.Join(footnoteLines, "\n\n")
		footnotes = strings.TrimSpace(footnotes)

		// If there are footnotes but there are no number, we will use asterisk as the marker.
		if footnoteNumbers.Size() == 0 && len(footnoteLines) > 0 {
			for i, line := range footnoteLines {
				if i == 0 {
					footnoteLines[i] = "[^*]: " + line
				} else {
					footnoteLines[i] = "    " + line
				}
			}

			footnoteNumbers.Put("*")
			footnotes = strings.Join(footnoteLines, "\n\n")
			footnotes = strings.TrimSpace(footnotes)
		}

		// Normalize footnote numbers in translation
		translation := ayah.Translation

		// Next we look using the provided patterns
		if rxTransFn != nil {
			translation = rxTransFn.ReplaceAllStringFunc(translation, func(match string) string {
				// Extract the footnote number
				parts := rxTransFn.FindStringSubmatch(match)
				nParts := len(parts)
				fnNumber := parts[1]

				// If this number doesn't exist in footnote, we remove it
				if !footnoteNumbers.Has(fnNumber) {
					if nParts > 2 {
						var replacementArgs []any
						for _, part := range parts[2:] {
							replacementArgs = append(replacementArgs, part)
						}

						replacementPattern := strings.Repeat("%s", nParts-2)
						return fmt.Sprintf(replacementPattern, replacementArgs...)
					} else {
						return ""
					}
				}

				// If this number do exist, we normalize it
				footnoteNumbers.Remove(fnNumber)

				var replacementArgs []any
				for _, part := range parts[1:] {
					replacementArgs = append(replacementArgs, part)
				}

				replacementPattern := "[^%s]"
				replacementPattern += strings.Repeat("%s", nParts-2)
				return fmt.Sprintf(replacementPattern, replacementArgs...)
			})
		}

		// Put the leftover footnote number at the end of translation.
		// Here we sort it to make sure each run give identical result.
		var fnNumbers []int
		footnoteNumbers.Each(func(k string) {
			if fnNumber, err := strconv.Atoi(k); err == nil {
				fnNumbers = append(fnNumbers, fnNumber)
			}
		})

		sort.Ints(fnNumbers)
		for _, fnNumber := range fnNumbers {
			translation += fmt.Sprintf("[^%d]", fnNumber)
		}

		// If asterisk used, put it at the very end
		if footnoteNumbers.Has("*") {
			translation += "[^*]"
		}

		// Apply normalized data
		ayah.Footnotes = footnotes
		ayah.Translation = translation
		data.AyahList[i] = ayah
	}

	return data
}

func splitFootnotesByNumber(data FlattenedData) FlattenedData {
	for i, ayah := range data.AyahList {
		fns := ayah.Footnotes
		fns = rxFootnoteNumberSplitter.ReplaceAllString(fns, "\n\n$1")
		fns = rxNewlines.ReplaceAllString(fns, "\n\n")
		fns = strings.TrimSpace(fns)
		data.AyahList[i].Footnotes = fns
	}

	return data
}

func noFootnote(data FlattenedData) FlattenedData {
	for i := range data.AyahList {
		data.AyahList[i].Footnotes = ""
	}

	return data
}
