package fic

import (
	"strconv"
	"strings"
	"time"

	GoHtml "github.com/udan-jayanith/GoHTML"
)

func handlerRating(f *Fic, ratingNode *GoHtml.Node) {
	ratingText := ratingNode.GetElementByTagName("ul").GetElementByTagName("li").GetElementByTagName("a").GetChildNode().GetText()
	f.Rating = RatingFromString(ratingText)
}

func handlerArchiveWarning(f *Fic, warningNode *GoHtml.Node) {
	warningNodes := warningNode.GetChildNode().GetElementsByTagName("li")
	for n := warningNodes.Next(); n != nil; n = warningNodes.Next() {
		f.ArchiveWarnings |= GetWarningFromString(n.GetChildNode().GetChildNode().GetText())
	}
}

func handlerCategory(f *Fic, categoryNode *GoHtml.Node) {
	categoryNodes := categoryNode.GetChildNode().GetElementsByTagName("li")
	for n := categoryNodes.Next(); n != nil; n = categoryNodes.Next() {
		f.Categories |= GetCategoryFromString(n.GetChildNode().GetChildNode().GetText())
	}
}

func handlerStrings(strs *[]string, stringsNode *GoHtml.Node) {
	stringNodes := stringsNode.GetChildNode().GetElementsByTagName("li")

	*strs = make([]string, 0, stringNodes.Len())

	for n := stringNodes.Next(); n != nil; n = stringNodes.Next() {
		*strs = append(*strs, n.GetChildNode().GetChildNode().GetText())
	}
}

func handlerLanguage(f *Fic, languageNode *GoHtml.Node) {
	f.Language = strings.TrimSpace(languageNode.GetChildNode().GetText())
}

func handlerStats(f *Fic, statsNode *GoHtml.Node) {
	keys := statsNode.GetElementsByTagName("dt")

	for key := keys.Next(); key != nil; key = keys.Next() {
		if cl, _ := key.GetAttribute("class"); cl == "status" {
			switch key.GetChildNode().GetText() {
			case "Updated:":
				f.Status = Incomplete
			case "Completed:":
				f.Status = Completed
			}

			break
		}
	}

	items := statsNode.GetElementsByTagName("dd")

	for item := items.Next(); item != nil; item = items.Next() {
		buf := make([]byte, 0, 9)

		cl, _ := item.GetAttribute("class")
		for _, ch := range []byte(cl) {
			if ch == ' ' {
				break
			}

			buf = append(buf, ch)
		}

		switch string(buf) {
		case "published":
			handlerStatTime(&f.PublishedAt, item)
		case "status":
			handlerStatTime(&f.UpdatedAt, item)
		case "words":
			handlerStatU32(&f.Words, item)
		case "chapters":
			handlerStatChapters(f, item)
		case "comments":
			handlerStatU32(&f.Comments, item)
		case "kudos":
			handlerStatU32(&f.Kudos, item)
		case "bookmarks":
			{
				n, err := strconv.ParseUint(strings.ReplaceAll(strings.TrimSpace(item.GetChildNode().GetChildNode().GetText()), ",", ""), 10, 32)
				if err != nil {
					panic(err.Error())
				}

				f.Bookmarks = uint32(n)
			}
		case "hits":
			handlerStatU32(&f.Hits, item)
		}
	}
}

func handlerStatTime(t *time.Time, timeNode *GoHtml.Node) {
	var err error
	*t, err = time.Parse(time.DateOnly, timeNode.GetChildNode().GetText())
	if err != nil {
		panic(err.Error())
	}
}

func handlerStatU32(u32 *uint32, u32Node *GoHtml.Node) {
	n, err := strconv.ParseUint(strings.ReplaceAll(strings.TrimSpace(u32Node.GetChildNode().GetText()), ",", ""), 10, 32)
	if err != nil {
		panic(err.Error())
	}

	*u32 = uint32(n)
}

func handlerStatChapters(f *Fic, chaptersNode *GoHtml.Node) {
	chapters := chaptersNode.GetChildNode().GetText()
	parts := strings.Split(chapters, "/")

	actual, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		panic(err.Error())
	}

	f.Chapters = uint32(actual)

	if parts[1] == "?" {
		f.MaxChapters = unknownMaxChaptersMask
	} else {
		max, err := strconv.ParseUint(parts[1], 10, 32)
		if err != nil {
			panic(err.Error())
		}

		f.MaxChapters = uint32(max) & unknownMaxChaptersAntiMask
	}
}
