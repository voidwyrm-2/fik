# Changelog

## 1.2.0

- (Fix) fictions with the 'Anonymous' author caused a nil pointer panic
- `remove` now supports removing multiple fictions at once
- `add` and `remove` can now use a file as a list of fictions with the `-f/--file` flag
- `add` now supports AO3 links, e.g. `fik add https://archiveofourown.org/works/52921651` or `fik add https://archiveofourown.org/works/59301142/chapters/151244305`
- Added the `clean` command
- `add` now sets the current chapter of a fiction to its first chapter when `--first` is passed if a current chapter isn't specified

## 1.1.0

- Adding fiction favoriting
- `add` now allows specifying the current chapter for each fiction
- Fictions can now have a current chapter set, allowing basic chapter bookmarking
- `list` can now filter fictions with the `-f/--filter` flag
- `list` can now limit the amount of fictions shown with the `-o/--only` flag
- Added the `fav`, `fics`, and `setchap` commands

## 1.0.0

- Basic commands implemented (`add`, `list`, `remove`, `show`, and `version`)
