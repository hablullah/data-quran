# Tanzil ![Collected: August 2022][badge-collect] ![Legal Status: Not OK][badge-legal]

[Tanzil][tanzil] is an international Quranic project aimed at providing a highly verified precise Quran text in Unicode. It was started back in 2007 and still [active][tanzil-updates] until now. Currently it's the most precise and accurate Quran texts available in internet. If you ever use Quranic apps, chance it uses texts from Tanzil.

Besides the Quran texts, Tanzil also provides Quranic metadata and translations. However, unlike the texts which extracted automatically from Medina Mushaf, the translations are provided by volunteer which means Tanzil cannot guarantee their authenticity and/or accuracy.

## Table of Contents
 
- [Terms of Use](#terms-of-use)
  - [Quran Texts](#quran-texts)
  - [Quran Metadata](#quran-metadata)
  - [Quran Translations](#quran-translations)
- [Legality](#legality)
  - [Quran Texts](#quran-texts-1)
  - [Quran Metadata](#quran-metadata-1)
  - [Quran Translations](#quran-translations-1)
- [Collected Data and Modifications](#collected-data-and-modifications)
  - [Quran Texts](#quran-texts-2)
  - [Quran Metadata](#quran-metadata-2)
  - [Quran Translations](#quran-translations-2)
- [Additional Notes](#additional-notes)

## Terms of Use

### Quran Texts

Quran [texts][tanzil-text] from Tanzil are [released][tanzil-text-license] under [CC BY 3.0][cc-by-3] license. This means we are allowed to use and redistribute it for free with following terms:

- **Attribution**. We must give proper attribution to Tanzil project and provide a link to the license.

There are also several additional rules from Tanzil:

- **No derivatives**. We are not allowed to modify or change the texts.
- **Add link to changelog**. We must provide link back to Tanzil project to enable end users to keep track of changes.

### Quran Metadata

Quran [metadata][tanzil-meta] from Tanzil are also released under [CC BY 3.0][cc-by-3] license which means we need to give proper attribution. However, there are no additional rules from Tanzil.

### Quran Translations

Quran [translations][tanzil-trans] from Tanzil are released without a specific license. However, we are required to:

- **Use it at our own discretion**. As mentioned above, Tanzil tried to provide a set of mostly acceptable Quran translations, however they can't guarantee their authenticity and/or accuracy.
- **Non commercial**. The translations provided by Tanzil are for non-commercial purposes only. If used otherwise, we need to obtain necessary permission from the translator or the publisher.
- **Add link back**. If using more than three of the translations provided by Tanzils, we required to add link to Tanzil [page][tanzil-trans] to enable end users to keep track of changes and use the latest updates.

## Legality

### Quran Texts

- [x] **License**. Tanzil uses CC BY 3.0 which compatible with CC BY-NC-ND 4.0 that used in this repository.
- [x] **Attribution**. We've put links to Tanzil project in this page. We also keep the license block in each text file.
- [x] **No derivatives**. Except changing the format from plain text into markdown, there are no change to the texts.
- [x] **Add link to changelog**. We've done so in this page and in license block of each text file.

### Quran Metadata

- [x] **License**. Metadata from Tanzil uses CC BY 3.0 which compatible with CC BY-NC-ND 4.0 that used in this repository.
- [x] **Attribution**. We've put links to Tanzil metadata in this page. Unfortunately, since we use JSON format for metadata, we can't keep the license block in each metadata file. However we already put a proper attribution in this page, so we hope that's enough.

### Quran Translations

- [x] **Non commercial**. This repository uses CC BY-NC-ND 4.0 which also doesn't allow commercial usage.
- [x] **Add link back**. We've put links to Tanzil in this page. We also keep the version and link bank in comment block in each translation files.

## Collected Data and Modifications

Data from Tanzil are collected automatically using the `cli`, except for surah info for Indonesian and Russian language which separated manually. The collected data from Tanzil are put in files prefixed with `*-tanzil.md`.

### Quran Texts

For the texts, we collected every text type from Tanzil:

- Simple
- Simple (Plain)
- Simple (Minimal)
- Simple (Clean)
- Uthmani
- Uthmani (Minimal)

Every text include pause marks, sajdah signs, rub-el-hizb signs and tatweel below superscript alefs. For Uthmani, we don't include texts with sequential tanweens because some Arabic font doesn't render it nicely.

All collected texts are the latest version, which is **version 1.1** that published in February 2021. The collected texts are located in `ayah-text` directory.

Except for format change from plain text to markdown, there are no modifications occured for the Quran texts.

### Quran Metadata

For the metadata, Tanzil provides a [JS file][tanzil-meta-js] which contains surah, juz, hizb quarter, manzil, ruku, page and sajda. From this JS file, we convert and split it into several JSON files:

- Sura into `surah/surah.json` and `surah-translation/en-tanzil.json`.
- Juz into `meta/juz.json`.
- Hizb quarter into `meta/maqra.json`. For the sake of completion, we also create `meta/hizb.json` from the hizb quarter.
- Manzil into `meta/manzil.json`.
- Ruku into `meta/ruku.json`.
- Page into `meta/page.json`.
- Sajda into `meta/sajda.json`.

All collected metadata are the latest version, which is **version 1.0** that published in February 2008. There are no modifications except for converting JS into JSON.

### Quran Translations

When these data collected, Tanzil provides 115 translations in plain text for 44 languages. After checking the files one-by-one, they are separated into four types:

- Translation
- Transliteration
- Tafsir
- Surah information

There are several modifications done to those translation files:

- Put `TODO:MISSING` to files where there are missing translation entries.
- Put `TODO:DUPLICATE` to files where there are duplicate entries.
- Remove HTML formatting from `en.transliteration`.
- Separate surah info from tafsir `id.muntakhab` and `ru.muntahab`.

## Additional Notes

As mentioned before, there are several translation files from Tanzil that missing some content:

- `cs.hrbek` missing translations for 3 ayah in 3427 (30:18), 4981 (56:2), 4992 (56:13).
- `ku.asan` missing translations for 1 ayah in 6207 (108:3).
- `sq.mehdiu` missing translations for 2 ayah in 2539 (21:56), 5636 (77:14).
- `fa.safavi` missing translations for 1 ayah in 5797 (80:39).

In `id.muntakhab`, there are 41 ayah with duplicate translations:

- Ayah 4273 (42:1) duplicated 1x in 4274 (42:2)
- Ayah 4584 (48:1) duplicated 2x in 4585-4586 (48:2-3)
- Ayah 4631 (50:1) duplicated 1x in 4632 (50:2)
- Ayah 4676 (51:1) duplicated 3x in 4677-4679 (51:2-4)
- Ayah 4736 (52:1) duplicated 5x in 4737-4741 (52:2-6)
- Ayah 4785 (53:1) duplicated 1x in 4786 (53:2)
- Ayah 4902 (55:1) duplicated 1x in 4903 (55:2)
- Ayah 4980 (56:1) duplicated 5x in 4981-4985 (56:2-6)
- Ayah 5324 (69:1) duplicated 1x in 5325 (69:2)
- Ayah 5376 (70:1) duplicated 2x in 5377-5378 (70:2-3)
- Ayah 5448 (72:1) duplicated 1x in 5449 (72:2)
- Ayah 5476 (73:1) duplicated 3x in 5477-5479 (73:2-4)
- Ayah 5496 (74:1) duplicated 3x in 5497-5499 (74:2-4)
- Ayah 5552 (75:1) duplicated 2x in 5553-5554 (75:2-3)
- Ayah 5623 (77:1) duplicated 6x in 5624-5629 (77:2-7)
- Ayah 5849 (83:1) duplicated 2x in 5850-5851 (83:2-3)
- Ayah 6169 (102:1) duplicated 1x in 6170 (102:2)
- Ayah 6194 (106:1) duplicated 1x in 6195 (106:2)

Besides that, `id.muntakhab` also missing info for 2 surah: An-Nisaa and Ibrahim.

While there are some missing translations, we believe majority of Tanzil translations are good enough to use. With that said, please use them at your own discretion.

[tanzil]: https://tanzil.net
[tanzil-updates]: http://tanzil.net/updates/
[tanzil-text]: https://tanzil.net/download/
[tanzil-text-license]: https://tanzil.net/docs/Text_License
[tanzil-meta]: https://tanzil.net/docs/quran_metadata
[tanzil-meta-js]: https://tanzil.net/res/text/metadata/quran-data.js
[tanzil-trans]: https://tanzil.net/trans/
[cc-by-3]: https://creativecommons.org/licenses/by/3.0/
