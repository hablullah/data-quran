# Tanzil ![Collected: August 2022][badge-collect] ![Legal Status: Not OK][badge-legal]

[Tanzil][tanzil] is an international Quranic project aimed at providing a highly verified precise Quran text in Unicode. It was started back in 2007 and still [active][tanzil-updates] until now. Currently it's the most precise and accurate Quran texts available in internet. If you ever use Quranic apps, chance it uses texts from Tanzil.

Besides the Quran texts, Tanzil also provides Quranic metadata and translations. However, unlike the texts which extracted automatically from Medina Mushaf, the translations are provided by volunteer which means Tanzil cannot guarantee their authenticity and/or accuracy.

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
- **Ask for permission directly**. This is only required if we are planning to redistribute the list of translations from Tanzil to another website.

## Legality

> **We** in this section refers to Hablullah team.

This section is used to explain whether we are allowed to put data from Tanzil in this repository or not. In short, it seems **we are fine for Quran texts and metadata**, however **we still need direct permission for the translations**.

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
- [ ] **Ask for permission directly**. Since we redistribute the translations here, we still need to ask for direct permission from Tanzil team.

## Collected Data and Modifications

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

When these data collected, Tanzil provides 114 translations in plain text for 44 languages. After checking the files one-by-one, they are separated into four types:

- Translation
- Transliteration
- Tafsir
- Surah information

All translations from Tanzil are prefixed with `*-tanzil.md`. For more detail, here is the table that show where and how the translation files from Tanzil are separated:

| Tanzil ID          | Last Update        | Destination                                                               |
| ------------------ | ------------------ | ------------------------------------------------------------------------- |
| ar.jalalayn        | December 1, 2010   | `ayah-tafsir/ar-jalalayn-tanzil.md`                                       |
| ar.muyassar        | January 7, 2011    | `ayah-tafsir/ar-muyassar-tanzil.md`                                       |
| fa.khorramdel      | March 19, 2011     | `ayah-tafsir/fa-khorramdel-tanzil.md`                                     |
| id.jalalayn        | July 27, 2012      | `ayah-tafsir/id-jalalayn-tanzil.md`                                       |
| id.muntakhab       | April 28, 2011     | `ayah-tafsir/id-muntakhab-tanzil.md`, `surah-info/id-muntakhab-tanzil.md` |
| ru.kuliev-alsaadi  | October 4, 2010    | `ayah-tafsir/ru-kuliev-alsaadi-tanzil.md`                                 |
| ru.muntahab        | May 29, 2011       | `ayah-tafsir/ru-muntahab-tanzil.md`, `surah-info/ru-muntahab-tanzil.md`   |
| uz.sodik           | July 17, 2011      | `ayah-tafsir/uz-sodik-tanzil.md`                                          |
| am.sadiq           | August 3, 2012     | `ayah-translation/am-sadiq-tanzil.md`                                     |
| az.mammadaliyev    | June 4, 2010       | `ayah-translation/az-mammadaliyev-tanzil.md`                              |
| az.musayev         | August 16, 2010    | `ayah-translation/az-musayev-tanzil.md`                                   |
| ber.mensur         | July 24, 2012      | `ayah-translation/ber-mensur-tanzil.md`                                   |
| bg.theophanov      | August 12, 2014    | `ayah-translation/bg-theophanov-tanzil.md`                                |
| bn.bengali         | April 30, 2011     | `ayah-translation/bn-bengali-tanzil.md`                                   |
| bn.hoque           | July 19, 2013      | `ayah-translation/bn-hoque-tanzil.md`                                     |
| bs.korkut          | May 22, 2013       | `ayah-translation/bs-korkut-tanzil.md`                                    |
| bs.mlivo           | August 24, 2010    | `ayah-translation/bs-mlivo-tanzil.md`                                     |
| cs.hrbek           | August 16, 2010    | `ayah-translation/cs-hrbek-tanzil.md`                                     |
| cs.nykl            | August 16, 2010    | `ayah-translation/cs-nykl-tanzil.md`                                      |
| de.aburida         | August 19, 2011    | `ayah-translation/de-aburida-tanzil.md`                                   |
| de.bubenheim       | July 17, 2011      | `ayah-translation/de-bubenheim-tanzil.md`                                 |
| de.khoury          | August 16, 2010    | `ayah-translation/de-khoury-tanzil.md`                                    |
| de.zaidan          | June 4, 2010       | `ayah-translation/de-zaidan-tanzil.md`                                    |
| dv.divehi          | May 3, 2022        | `ayah-translation/dv-divehi-tanzil.md`                                    |
| en.ahmedali        | August 28, 2010    | `ayah-translation/en-ahmedali-tanzil.md`                                  |
| en.ahmedraza       | July 17, 2011      | `ayah-translation/en-ahmedraza-tanzil.md`                                 |
| en.arberry         | July 31, 2011      | `ayah-translation/en-arberry-tanzil.md`                                   |
| en.daryabadi       | August 16, 2010    | `ayah-translation/en-daryabadi-tanzil.md`                                 |
| en.hilali          | December 13, 2010  | `ayah-translation/en-hilali-tanzil.md`                                    |
| en.itani           | July 19, 2013      | `ayah-translation/en-itani-tanzil.md`                                     |
| en.maududi         | May 10, 2011       | `ayah-translation/en-maududi-tanzil.md`                                   |
| en.mubarakpuri     | June 20, 2015      | `ayah-translation/en-mubarakpuri-tanzil.md`                               |
| en.pickthall       | September 4, 2010  | `ayah-translation/en-pickthall-tanzil.md`                                 |
| en.qarai           | April 11, 2022     | `ayah-translation/en-qarai-tanzil.md`                                     |
| en.qaribullah      | June 4, 2010       | `ayah-translation/en-qaribullah-tanzil.md`                                |
| en.sahih           | April 24, 2011     | `ayah-translation/en-sahih-tanzil.md`                                     |
| en.sarwar          | July 30, 2012      | `ayah-translation/en-sarwar-tanzil.md`                                    |
| en.shakir          | June 22, 2010      | `ayah-translation/en-shakir-tanzil.md`                                    |
| en.wahiduddin      | August 25, 2011    | `ayah-translation/en-wahiduddin-tanzil.md`                                |
| en.yusufali        | May 10, 2013       | `ayah-translation/en-yusufali-tanzil.md`                                  |
| es.bornez          | July 15, 2014      | `ayah-translation/es-bornez-tanzil.md`                                    |
| es.cortes          | August 12, 2014    | `ayah-translation/es-cortes-tanzil.md`                                    |
| es.garcia          | July 15, 2014      | `ayah-translation/es-garcia-tanzil.md`                                    |
| fa.ansarian        | July 6, 2011       | `ayah-translation/fa-ansarian-tanzil.md`                                  |
| fa.ayati           | August 1, 2012     | `ayah-translation/fa-ayati-tanzil.md`                                     |
| fa.bahrampour      | November 29, 2013  | `ayah-translation/fa-bahrampour-tanzil.md`                                |
| fa.fooladvand      | May 3, 2022        | `ayah-translation/fa-fooladvand-tanzil.md`                                |
| fa.gharaati        | April 11, 2022     | `ayah-translation/fa-gharaati-tanzil.md`                                  |
| fa.ghomshei        | July 17, 2011      | `ayah-translation/fa-ghomshei-tanzil.md`                                  |
| fa.khorramshahi    | July 27, 2012      | `ayah-translation/fa-khorramshahi-tanzil.md`                              |
| fa.makarem         | January 22, 2014   | `ayah-translation/fa-makarem-tanzil.md`                                   |
| fa.moezzi          | August 24, 2010    | `ayah-translation/fa-moezzi-tanzil.md`                                    |
| fa.mojtabavi       | April 16, 2012     | `ayah-translation/fa-mojtabavi-tanzil.md`                                 |
| fa.sadeqi          | August 1, 2011     | `ayah-translation/fa-sadeqi-tanzil.md`                                    |
| fa.safavi          | April 11, 2022     | `ayah-translation/fa-safavi-tanzil.md`                                    |
| fr.hamidullah      | July 18, 2011      | `ayah-translation/fr-hamidullah-tanzil.md`                                |
| ha.gumi            | August 16, 2010    | `ayah-translation/ha-gumi-tanzil.md`                                      |
| hi.farooq          | March 18, 2011     | `ayah-translation/hi-farooq-tanzil.md`                                    |
| hi.hindi           | January 12, 2011   | `ayah-translation/hi-hindi-tanzil.md`                                     |
| id.indonesian      | June 4, 2010       | `ayah-translation/id-indonesian-tanzil.md`                                |
| it.piccardo        | January 2, 2011    | `ayah-translation/it-piccardo-tanzil.md`                                  |
| ja.japanese        | June 4, 2010       | `ayah-translation/ja-japanese-tanzil.md`                                  |
| ko.korean          | July 15, 2011      | `ayah-translation/ko-korean-tanzil.md`                                    |
| ku.asan            | September 14, 2010 | `ayah-translation/ku-asan-tanzil.md`                                      |
| ml.abdulhameed     | April 2, 2012      | `ayah-translation/ml-abdulhameed-tanzil.md`                               |
| ml.karakunnu       | April 2, 2012      | `ayah-translation/ml-karakunnu-tanzil.md`                                 |
| ms.basmeih         | September 7, 2012  | `ayah-translation/ms-basmeih-tanzil.md`                                   |
| nl.keyzer          | June 4, 2010       | `ayah-translation/nl-keyzer-tanzil.md`                                    |
| nl.leemhuis        | August 5, 2012     | `ayah-translation/nl-leemhuis-tanzil.md`                                  |
| nl.siregar         | August 5, 2012     | `ayah-translation/nl-siregar-tanzil.md`                                   |
| no.berg            | June 4, 2010       | `ayah-translation/no-berg-tanzil.md`                                      |
| pl.bielawskiego    | August 16, 2010    | `ayah-translation/pl-bielawskiego-tanzil.md`                              |
| ps.abdulwali       | June 29, 2016      | `ayah-translation/ps-abdulwali-tanzil.md`                                 |
| pt.elhayek         | June 4, 2010       | `ayah-translation/pt-elhayek-tanzil.md`                                   |
| ro.grigore         | August 16, 2010    | `ayah-translation/ro-grigore-tanzil.md`                                   |
| ru.abuadel         | September 15, 2010 | `ayah-translation/ru-abuadel-tanzil.md`                                   |
| ru.krachkovsky     | September 20, 2010 | `ayah-translation/ru-krachkovsky-tanzil.md`                               |
| ru.kuliev          | May 29, 2011       | `ayah-translation/ru-kuliev-tanzil.md`                                    |
| ru.osmanov         | August 16, 2010    | `ayah-translation/ru-osmanov-tanzil.md`                                   |
| ru.porokhova       | August 16, 2010    | `ayah-translation/ru-porokhova-tanzil.md`                                 |
| ru.sablukov        | October 4, 2010    | `ayah-translation/ru-sablukov-tanzil.md`                                  |
| sd.amroti          | April 9, 2012      | `ayah-translation/sd-amroti-tanzil.md`                                    |
| so.abduh           | August 16, 2010    | `ayah-translation/so-abduh-tanzil.md`                                     |
| sq.ahmeti          | August 16, 2010    | `ayah-translation/sq-ahmeti-tanzil.md`                                    |
| sq.mehdiu          | December 1, 2010   | `ayah-translation/sq-mehdiu-tanzil.md`                                    |
| sq.nahi            | August 16, 2010    | `ayah-translation/sq-nahi-tanzil.md`                                      |
| sv.bernstrom       | July 30, 2012      | `ayah-translation/sv-bernstrom-tanzil.md`                                 |
| sw.barwani         | August 16, 2010    | `ayah-translation/sw-barwani-tanzil.md`                                   |
| ta.tamil           | August 16, 2010    | `ayah-translation/ta-tamil-tanzil.md`                                     |
| tg.ayati           | August 4, 2010     | `ayah-translation/tg-ayati-tanzil.md`                                     |
| th.thai            | October 10, 2011   | `ayah-translation/th-thai-tanzil.md`                                      |
| tr.ates            | June 4, 2010       | `ayah-translation/tr-ates-tanzil.md`                                      |
| tr.bulac           | August 16, 2010    | `ayah-translation/tr-bulac-tanzil.md`                                     |
| tr.diyanet         | December 27, 2011  | `ayah-translation/tr-diyanet-tanzil.md`                                   |
| tr.golpinarli      | September 14, 2010 | `ayah-translation/tr-golpinarli-tanzil.md`                                |
| tr.ozturk          | June 4, 2010       | `ayah-translation/tr-ozturk-tanzil.md`                                    |
| tr.vakfi           | June 4, 2010       | `ayah-translation/tr-vakfi-tanzil.md`                                     |
| tr.yazir           | June 4, 2010       | `ayah-translation/tr-yazir-tanzil.md`                                     |
| tr.yildirim        | September 14, 2010 | `ayah-translation/tr-yildirim-tanzil.md`                                  |
| tr.yuksel          | June 4, 2010       | `ayah-translation/tr-yuksel-tanzil.md`                                    |
| tt.nugman          | August 16, 2010    | `ayah-translation/tt-nugman-tanzil.md`                                    |
| ug.saleh           | June 4, 2010       | `ayah-translation/ug-saleh-tanzil.md`                                     |
| ur.ahmedali        | August 16, 2010    | `ayah-translation/ur-ahmedali-tanzil.md`                                  |
| ur.jalandhry       | December 24, 2010  | `ayah-translation/ur-jalandhry-tanzil.md`                                 |
| ur.jawadi          | December 24, 2010  | `ayah-translation/ur-jawadi-tanzil.md`                                    |
| ur.junagarhi       | April 25, 2011     | `ayah-translation/ur-junagarhi-tanzil.md`                                 |
| ur.kanzuliman      | March 17, 2011     | `ayah-translation/ur-kanzuliman-tanzil.md`                                |
| ur.maududi         | November 15, 2010  | `ayah-translation/ur-maududi-tanzil.md`                                   |
| ur.najafi          | August 19, 2011    | `ayah-translation/ur-najafi-tanzil.md`                                    |
| ur.qadri           | August 16, 2010    | `ayah-translation/ur-qadri-tanzil.md`                                     |
| zh.jian            | March 13, 2011     | `ayah-translation/zh-jian-tanzil.md`                                      |
| zh.majian          | January 7, 2011    | `ayah-translation/zh-majian-tanzil.md`                                    |
| en.transliteration | September 6, 2010  | `ayah-transliteration/en-tanzil.md`                                       |
| tr.transliteration | September 15, 2010 | `ayah-transliteration/tr-tanzil.md`                                       |

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

[badge-collect]: https://img.shields.io/badge/collected-Aug%202022-lightgrey?style=flat-square
[badge-legal]: https://img.shields.io/badge/legality-not--ok-red?style=flat-square
[tanzil]: https://tanzil.net
[tanzil-updates]: http://tanzil.net/updates/
[tanzil-text]: https://tanzil.net/download/
[tanzil-text-license]: https://tanzil.net/docs/Text_License
[tanzil-meta]: https://tanzil.net/docs/quran_metadata
[tanzil-meta-js]: https://tanzil.net/res/text/metadata/quran-data.js
[tanzil-trans]: https://tanzil.net/trans/
[cc-by-3]: https://creativecommons.org/licenses/by/3.0/
