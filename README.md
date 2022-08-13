# Data Quran [![License: CC BY-NC-ND 4.0][cc-badge]][cc-url]

This repository is collection of free dataset for everything related to Quran: from the text, translation, word-by-word, and tafsir. There are several reasons why this repository created:

1. **To centralize all Quran dataset in one place.**

   Currently, to create Quran apps, developers need to gather data from various sources either by downloading, using API, or scraping manually. It would be nice if there is a single repository to gather them all.

2. **To standardize the dataset format.**

   Each dataset source usually has their own format which means the developer need to parse and normalize each of them. It would be nice if all of those dataset only use a single type of formatting.

3. **To archive dataset, in case the original source goes down or unreachable.**

   There are several useful Quranic website that went down after being inactive for several years. There are also cases where governments decided to ban Quranic apps from app stores. Hopefully this repository can be used as archive so those useful data doesn't vanish even after the original websites are gone.

4. **To make sure all dataset are legal and not infringing copyright.**

   In some countries, Quran translations are protected by copyright. Since copyright are recognized by Islam, in this repository all dataset must be collected from valid and legal source, both according to government law and Islamic law.

5. **To give proper attributions and explanation on how the data collected.**

   There are several other repositories that also collects the Quranic data. However, as far as I know all of them doesn't really mention the source and how the data are collected.

## Data Format

### Criteria

When choosing format for this repository, there are several criterias that must be fulfilled:

1. It must be usable across programming languages.
2. It must be platform agnostic, and doesn't require specific app to use.
3. It must supports multi-line text.
4. It must supports rich-text formatting and footnotes.
5. It must be easy to read and write even for non-programmers.
6. It must be usable with Git, and the diff must be easy to read.

### Chosen Format

There are two format that chosen for this repository.

The first is `json` which is the universal data format across all programming language. It's used for all Quranic data where the value are short, i.e. Quran metadata, surah, and word-by-word translation. The reason it's chosen are:

1. Every programming languages support JSON, and most of them include JSON parser and decoder in standard library.
2. It can be opened and edited by every text editor in every operating system.
3. The properly formatted JSON files are easy to read and write, ever for common people.
4. Since it's just a text file, it's trackable using Git.

The only downside for JSON is we can't easily put multi-line or rich-text content as JSON value. While it's doable, it's not really easy to read and common people usually don't know how. This is why we only use it for Quranic data with short values.

---

The second format is `markdown` which is used for all Quranic data where the values are a long or multi line strings. This include the Quran text, translation, transliteration and tafseer. The reason it's chosen are:

1. Most programming languages have third-party library for encoding and decoding markdown languages.
2. It can be opened and edited by every text editor in every operating system.
3. It supports multi-line texts.
4. It supports rich-text formatting, and there is extension to make markdown supports footnotes.
5. It's easy to read and write.
6. It's also a text file, so it's trackable using Git.

The markdown files in this repository are formatted like this:

```md
<!--
Comment block for license or metadata
-->

# [verse-id-1]

The content for this verse.

# [verse-id-2]

The content for this verse.
```

### Why Not \[Other Format\]

There are several other formats that considered to be used in this repository.

The first is **plain text**. In this format, each verse only use one line, which make it compact and easy to read and write.

Pros:

- It can be opened by every text editor.
- It's easy enough to read and write.
- It's trackable with Git.

Cons:

- Since each verse only use one line, multi-line text is not supported.
- It doesn't support rich-text formatting and footnotes.
- We can force it to support multi-line and rich-text format by using HTML tags like `<br>`, `<u>`, and `<b>`. However, by doing so, now it's hard to read and write by common people (which remove the pros of this format).

---

The second candidate is **CSV** format.

Pros:

- It can be opened by every text editor and spreadsheet programs.
- It's easy enough to read and write, especially when edited using spreadsheet programs.
- It supports multi-line texts.
- It's text file which make it trackable with Git.

Cons:

- Default CSV symbols (i.e. separator and quote) differs depending on user locale. This could lead to problem when editing the file.
- It doesn't support rich-text formatting and footnotes.
- We can force it to support multi-line and rich-text format by using HTML tags. However, by doing so, now it's hard to read and write by common people.

---

We also considered using **XML** format, but it's immediately rejected because it's hard to read and write by common people.

## Repository Structure

This repository is composed by several directories:

- [`meta`](meta) contains metadata that used in Quran.
- [`ayah-text`](ayah-text) contains Arabic text that used in Quran.
- [`ayah-transliteration`](ayah-transliteration) contains transliteration from Arabic to Latin scripts for each verse. Useful for those starting to learn how to read Quran.
- [`ayah-translation`](ayah-translation) contains the translations for each verse in Quran.
- [`ayah-tafsir`](ayah-tafsir) contains additional explanation for each verse in Quran.
- [`surah`](surah) contains Arabic name, data and ayah range for each surah in Quran.
- [`surah-info`](surah-info) contains descriptions and additional info for each surah in Quran.
- [`surah-translation`](surah-translation) contains the translation from Arabic name of each surah.
- [`source`](source) contains the explanation on where and how data in this repository collected.

## Contributions

Like other open source projects, we are open to suggestions and corrections. Feel free to submit your issues if there are any error in the dataset. However, there is a special rule for pull requests.

**Every pull requests that modifies data must be done following the terms of use from the original source**. So, there are two cases:

1. **The original source allows data modification.**

   For example, Tanzil released their translations for free and their terms of use doesn't disallow data modification. In this case, PRs for translations from Tanzil are allowed.

2. **The original source prohibits data modification.**

   For example, Tanzil released their Quran texts for free, but they explicitly state that changing the Quran texts are not allowed. In this case, any PRs that modifies the Quran texts from Tanzil are not allowed and will never be accepted.

   If you found an issue but the source doesn't allow data modifications, you should contact the original source and ask them to correct their data. Once they make the correction upstream, we will update the data in this repository.

## License

This repository is available under **CC BY-NC-ND 4.0** license. This means you can use this repository for free under following terms:

- **Attribution**. You must give appropriate credit to this repository and provide a link to the license. Check out Creative Commons [guide][attr-guide] on how to give attribution.

  If possible, please also include the original sources on your attribution. For example:

  > Data is taken from [data-quran] repository which licensed under [CC BY-NC-ND 4.0][cc-url] and collected by [Hablullah team][hablullah] from various sources, e.g. Tanzil, QuranEnc, etc.

- **Non commercial**. You may not use data from this repository for commercial purpose. This includes one-time purchase, subscription, in-app purchase, and in-app advertising. However you are allowed to ask donation for your apps, as long as it's not mandatory.

- **No derivatives**. You are not allowed to publish derivative work from this repository. Derivative here means any modifications including translations, revisions, annotations, elaborations, or any other modifications that based on this repository.

  If you have any modifications or revisions, you must submit it as pull request to this repository. This is done to make sure this repository stays as the single source of truth ([SSOT]) and to prevent confusions between multiple forks.

  However, you are allowed to change data formats to make it suitable for your applications. So, even though this repository publish data in `json` and `markdown` format, you can safely convert it to SQL format. For more details, check out section 2.a.4 in license page and this [FAQ][cc-faq] from Creative Commons.

[attr-guide]: https://creativecommons.org/use-remix/attribution/
[data-quran]: https://github.com/hablullah/data-quran
[hablullah]: https://github.com/orgs/hablullah/people
[cc-badge]: https://i.creativecommons.org/l/by-nc-nd/4.0/80x15.png
[cc-url]: https://creativecommons.org/licenses/by-nc-nd/4.0/
[cc-faq]: https://creativecommons.org/faq/#can-i-take-a-cc-licensed-work-and-use-it-in-a-different-format
[ssot]: https://en.wikipedia.org/wiki/Single_source_of_truth
