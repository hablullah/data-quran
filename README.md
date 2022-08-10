# Data Quran [![License: CC BY-NC-ND 4.0][cc-badge]][cc-url]

This repository contains free dataset for everything related to Quran: from the text, translation, word-by-word, and tafseer (explanation). There are several reasons why this repository created:

1. To centralize all Quran dataset in one place.

   Currently, to create Quran apps, developers need to gather data from various sources manually. Not to mention each source usually has their own formatting, which means the developer also need to parse and normalize the dataset. Once this project complete, developers can simply use this repository as the single source for all Quran related data.

2. To archive Quran dataset, in case the original source goes down or unreachable.

   There are several useful Quranic website that went down after being inactive for several years. Hopefully, this repository can be used as archive so those useful data doesn't vanish even after the original websites are gone.

3. To make sure all dataset are legal and not infringing copyright.

   I'm not sure how is it in other countries, but in Indonesia some of Quran translations are protected by copyright. Since copyright are recognized by Islam, in this repository all dataset must be collected legally, both by government law and Islamic law.

4. To give proper attributions and explanation on how the data collected.

   There are several other repositories that collect the Quranic data, however as far as I know, all of them doesn't really mention why, where, when, and how the data are collected.

## Formatting

There are two kind of format that used in this repository.

The first is `json` which is the universal data format across all programming language. It's used for all Quranic data where the value are short, i.e. Quran metadata and word-by-word translation.

The second is `markdown` which is used for all Quranic data where the values are a long or multi line strings. This include the Quran text, translation, transliteration and tafseer. Markdown is chosen because it supports formatting and footnotes (which often used in tafseer) while still simple enough for ordinary people to use. Another advantage for using Markdown is the diff support is great, so we can easily track where the change occured.

The markdown files in this repository are formatted like this:

```md
# [verse-id-1]

The content for this verse.

# [verse-id-2]

The content for this verse.
```

For example, here is how the markdown file for Al-Fatihah translation looks like:

```md
# 1

In the Name of Allâh, the Most Gracious, the Most Merciful

# 2

All the praises and thanks be to Allâh, the Lord[^1] of the ‘Âlamîn (mankind, jinn and all that exists).[^2]

[^1]: Lord: The actual word used in the Qur’ân is Rabb There is no proper equivalent for Rabb in English language. It means the One and the Only Lord for all the universe, its Creator, Owner, Organizer, Provider, Master, Planner, Sustainer, Cherisher, and Giver of security. Rabb is also one of the Names of Allâh. We have used the word "Lord" as the nearest to Rabb . All occurrences of "Lord" in the interpretation of the meanings of the Noble Qur’ân actually mean Rabb and should be understood as such.
[^2]: Narrated Abu Sa‘îd bin Al-Mu‘alla: While I was praying in the mosque, Allâh’s Messenger صلى الله عليه وسلم called me but I did not respond to him. Later I said, "O Allâh’s Messenger, I was praying." He said, "Didn’t Allâh say - Answer Allâh (by obeying Him) and His Messenger when he صلى الله عليه وسلم calls you." (V.8:24). He then said to me, "I will teach you a Sûrah which is the greatest Sû rah in the Qur’ân, before you leave the mosque." Then he got hold of my hand, and when he intended to leave (the mosque), I said to him, "Didn’t you say to me, "I will teach you a Sûrah which is the greatest Sûrah in the Qur’an?" He said, " Al-Hamdu lillahi Rabbil-‘âlamîn [i.e. all the praises and thanks be to Allâh, the Lord of the ‘Âlamîn (mankind, jinn and all that exists)], Sûrat Al-Fâtihah which is As-Sab‘ Al-Mathâni (i.e. the seven repeatedly recited Verses) and the Grand Qur’ân which has been given to me." (Sahih Al-Bukhâri, Vol.6, Hadîth No. 1).
```

# License

This repository is available under **CC BY-NC-ND 4.0** license. This means you can use this repository for free under these following terms:

- **Attribution**. You must give appropriate credit to this repository and provide a link to the license. Check out Creative Commons [guide][attr-guide] on how to give attribution.

  If possible, please also include the original sources on your attribution. For example:

  > Data is taken from [data-quran] repository which licensed under [CC BY-NC-ND 4.0][cc-url] and collected by [Hablullah team][hablullah] from various sources, e.g. Tanzil, QuranEnc, etc.

- **NonCommercial**. You may not use data from this repository for commercial purpose. This includes one-time purchase, subscription, in-app purchase, and in-app advertising. However you are allowed to ask donation for your apps, as long as it's not mandatory.

- **NoDerivatives**. You are not allowed to publish derivative work from this repository. Derivative here means any modifications including translations, revisions, annotations, elaborations, or any other modifications that based on this repository.

  If you have any modifications or revisions, you must submit it as pull request to this repository. This is done to make sure this repository stays as the single source of truth ([SSOT]) and to prevent confusions between multiple forks.

  However, you are allowed to change data formats to make it suitable for your applications. So, even though this repository publish data in `json` and `markdown` format, you can safely convert it to SQL format. For more details, check out section 2.a.4 in license page and this [FAQ][cc-faq] from Creative Commons.

[attr-guide]: https://creativecommons.org/use-remix/attribution/
[data-quran]: https://github.com/hablullah/data-quran
[hablullah]: https://github.com/orgs/hablullah/people
[cc-badge]: https://i.creativecommons.org/l/by-nc-nd/4.0/80x15.png
[cc-url]: https://creativecommons.org/licenses/by-nc-nd/4.0/
[cc-faq]: https://creativecommons.org/faq/#can-i-take-a-cc-licensed-work-and-use-it-in-a-different-format
[ssot]: https://en.wikipedia.org/wiki/Single_source_of_truth
