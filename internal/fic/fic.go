package fic

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/udan-jayanith/GoHTML"
)

type Id uint32

func ParseId(str string) (Id, error) {
	n, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("'%s' is not a valid AO3 fiction ID", str)
	}

	return Id(n), nil
}

func ParseIdPair(str string) (id, chapterId Id, err error) {
	var strId, strChapterId string

	parts := strings.Split(str, ",")
	strId = strings.TrimSpace(parts[0])
	if len(parts) > 1 {
		strChapterId = strings.TrimSpace(parts[1])
	}

	id, err = ParseId(strId)
	if err != nil {
		err = fmt.Errorf("'%s' is not a valid AO3 fiction ID", str)
		return
	}

	if len(strChapterId) > 0 {
		chapterId, err = ParseId(strChapterId)
		if err != nil {
			err = fmt.Errorf("'%s' is not a valid AO3 fiction chapter ID", str)
			return
		}
	}

	return
}

type Rating uint8

const (
	Unrated Rating = iota
	General
	Teen
	Mature
	Explicit
)

func RatingFromString(str string) Rating {
	str = strings.ToLower(str)

	if strings.Contains(str, "general") {
		return General
	}

	if strings.Contains(str, "teen") {
		return Teen
	}

	if strings.Contains(str, "mature") {
		return Mature
	}

	if strings.Contains(str, "explicit") {
		return Explicit
	}

	return Unrated
}

func (r Rating) String() string {
	switch r {
	case General:
		return "General"
	case Teen:
		return "Teen and Up"
	case Mature:
		return "Mature"
	case Explicit:
		return "Explicit"
	default:
		return "Not Rated"
	}
}

type Category uint8

const (
	Unpaired Category = 0
	Gen               = 1 << iota
	Other
	FM
	FF
	MM
)

const (
	_unpairedString = "None"
	_genString      = "Gen"
	_otherString    = "Other"
	_fmString       = "F/M"
	_ffString       = "F/F"
	_mmString       = "M/M"
)

func GetCategoryFromString(str string) Category {
	switch strings.TrimSpace(str) {
	case "Multi":
		return Unpaired
	case "Gen":
		return Gen
	case "F/M":
		return FM
	case "F/F":
		return FF
	case "M/M":
		return MM
	default:
		return Other
	}
}

func (c Category) String() string {
	if c == Unpaired {
		return _unpairedString
	} else if c&Gen != 0 {
		return _genString
	}

	categories := []string{}

	if c&Other != 0 {
		categories = append(categories, _otherString)
	}

	if c&FM != 0 {
		categories = append(categories, _fmString)
	}

	if c&FF != 0 {
		categories = append(categories, _ffString)
	}

	if c&MM != 0 {
		categories = append(categories, _mmString)
	}

	return strings.Join(categories, ", ")
}

type ArchiveWarning uint8

const (
	NoWarnings             ArchiveWarning = 0
	CreatorChoseNoWarnings                = 1 << iota
	GraphicViolence
	MajorCharacterDeath
)

const (
	_noWarningsString             = "No Archive Warnings Apply"
	_creatorChoseNoWarningsString = "Creator Chose Not To Use Archive Warnings"
	_graphicViolenceString        = "Graphic Depictions Of Violence"
	_majorCharacterDeathString    = "Major Character Death"
)

func GetWarningFromString(str string) ArchiveWarning {
	switch strings.TrimSpace(str) {
	case _creatorChoseNoWarningsString:
		return CreatorChoseNoWarnings
	case _graphicViolenceString:
		return GraphicViolence
	case _majorCharacterDeathString:
		return MajorCharacterDeath
	default:
		return NoWarnings
	}
}

func (aw ArchiveWarning) String() string {
	if aw == NoWarnings {
		return _noWarningsString
	} else if aw&CreatorChoseNoWarnings != 0 {
		return _creatorChoseNoWarningsString
	}

	warnings := []string{}

	if aw&GraphicViolence != 0 {
		warnings = append(warnings, _graphicViolenceString)
	}

	if aw&MajorCharacterDeath != 0 {
		warnings = append(warnings, _majorCharacterDeathString)
	}

	return strings.Join(warnings, ", ")
}

type Status uint8

const (
	Unknown Status = iota
	Incomplete
	Completed
)

func (s Status) Format(t time.Time) string {
	var status string

	switch s {
	case Incomplete:
		status = "Updated"
	case Completed:
		status = "Completed"
	default:
		return "Completion Status: Unknown"
	}

	return fmt.Sprintf("%s: %s", status, t.Format(time.DateOnly))
}

const (
	unknownMaxChaptersMask     uint32 = 0b1000000 << 3
	unknownMaxChaptersAntiMask uint32 = ^unknownMaxChaptersMask
)

type Fic struct {
	ChapterInfo struct {
		Title string
		Id    Id
		Num   uint32
	}
	Fandoms, Relationships, Characters, Tags                       []string
	Title, Author, Summary, Language                               string
	PublishedAt, UpdatedAt                                         time.Time
	Id                                                             Id
	Words, Chapters, MaxChapters, Comments, Kudos, Bookmarks, Hits uint32
	Status                                                         Status
	Rating                                                         Rating
	Categories                                                     Category
	ArchiveWarnings                                                ArchiveWarning
	Favorite                                                       bool
}

func GetFicFromId(id, chapterId Id) (f Fic, err error) {
	f.Id = id
	f.ChapterInfo.Id = chapterId

	res, err := http.Get(fmt.Sprintf("https://archiveofourown.org/works/%d?view_adult=true?hide_banner=true", id))
	if err != nil {
		return
	}

	defer res.Body.Close()

	root, err := GoHtml.Decode(res.Body)
	if err != nil {
		return
	}

	main := root.GetElementByTagName("body").GetElementByID("outer").GetElementByID("inner").GetElementByID("main")
	work := main.GetElementByClassName("work")
	workskin := work.GetElementByID("workskin")

	if f.ChapterInfo.Id != 0 {
		err = f.GetCurrentChapterInfo()
		if err != nil {
			return
		}
	}

	{
		headings := main.GetElementsByTagName("h2")

		for n := headings.Next(); n != nil; n = headings.Next() {
			if cl, _ := n.GetAttribute("class"); cl == "heading" && strings.TrimSpace(n.GetInnerText()) == "Error 404" {
				err = fmt.Errorf("Fiction with ID %d does not exist", id)
				return
			}
		}
	}

	info := work.GetElementByClassName("wrapper").GetElementByTagName("dl")
	items := info.GetElementsByTagName("dd")

	for item := items.Next(); item != nil; item = items.Next() {
		buf := make([]byte, 0, 12)

		cl, _ := item.GetAttribute("class")
		for _, ch := range []byte(cl) {
			if ch == ' ' {
				break
			}

			buf = append(buf, ch)
		}

		var handler func(*Fic, *GoHtml.Node)

		switch string(buf) {
		case "rating":
			handler = handlerRating
		case "warning":
			handler = handlerArchiveWarning
		case "category":
			handler = handlerCategory
		case "fandom":
			handlerStrings(&f.Fandoms, item)
		case "relationship":
			handlerStrings(&f.Relationships, item)
		case "character":
			handlerStrings(&f.Characters, item)
		case "freeform":
			handlerStrings(&f.Tags, item)
		case "language":
			handler = handlerLanguage
		case "stats":
			handler = handlerStats
		}

		if handler != nil {
			handler(&f, item)
		}
	}

	preface := workskin.GetElementByClassName("preface")

	f.Title = strings.TrimSpace(preface.GetElementByClassName("title").GetChildNode().GetText())
	f.Author = preface.GetElementByClassName("byline").GetChildNode().GetChildNode().GetText()

	summaryNodes := preface.GetElementByClassName("summary").GetElementByTagName("blockquote").GetElementsByTagName("p")

	parts := make([]string, 0, summaryNodes.Len())

	for n := summaryNodes.Next(); n != nil; n = summaryNodes.Next() {
		parts = append(parts, n.GetInnerText())
	}

	f.Summary = strings.Join(parts, "\n\n ")

	return
}

func (f *Fic) GetCurrentChapterInfo() error {
	res, err := http.Get(fmt.Sprintf("https://archiveofourown.org/works/%d/chapters/%d?view_adult=true?hide_banner=true", f.Id, f.ChapterInfo.Id))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	root, err := GoHtml.Decode(res.Body)
	if err != nil {
		return err
	}

	main := root.GetElementByTagName("body").GetElementByID("outer").GetElementByID("inner").GetElementByID("main")
	work := main.GetElementByClassName("work")
	workskin := work.GetElementByID("workskin")

	chapters := workskin.GetElementByID("chapters")
	chapter := chapters.GetChildNode()

	chns, _ := chapter.GetAttribute("id")

	parts := strings.Split(chns, "-")

	chn, err := strconv.ParseUint(strings.TrimSpace(parts[len(parts)-1]), 10, 32)
	if err != nil {
		return err
	}

	f.ChapterInfo.Num = uint32(chn)

	preface := chapter.GetElementsByClassName("chapter")
	titleParent := preface.Next()
	titleNode := titleParent.GetElementByClassName("title")

	if titleNode == nil {
		titleParent = preface.Next()
		titleNode = titleParent.GetElementByClassName("title")
	}

	f.ChapterInfo.Title = strings.TrimSpace(strings.Join(strings.Split(titleNode.GetInnerText(), ":")[1:], ":"))

	return nil
}

func (f *Fic) FormatSmall() string {
	return fmt.Sprintf("%s by %s (%d, %s, %s)", f.Title, f.Author, f.Id, f.Rating, f.Fandoms[0])
}

func (f *Fic) FormatSmallColor() string {
	return fmt.Sprintf("%s by %s (%d, %s, %s)", f.Title, f.Author, f.Id, f.Rating, f.Fandoms[0])
}

func (f *Fic) String() string {
	var maxChapters string

	if f.MaxChapters&unknownMaxChaptersMask != 0 {
		maxChapters = "?"
	} else {
		maxChapters = fmt.Sprint(f.MaxChapters)
	}

	var status string

	if f.MaxChapters > 1 {
		status = "\n" + f.Status.Format(f.UpdatedAt)
	}

	var currentChapter string

	if f.ChapterInfo.Id != 0 {
		currentChapter = fmt.Sprintf("\nCurrent Chapter: %d, `%s` (id %d)", f.ChapterInfo.Num, f.ChapterInfo.Title, f.ChapterInfo.Id)
	}

	return fmt.Sprintf(`Title: %s
Author: %s
Rating: %s
Archive Warnings: %s
Categories: %s
Fandoms: %s
Relationships: %s
Characters: %s
Tags: %s
Language: %s
Published At: %s%s
Chapters: %d/%s%s
Comments: %d
Kudos: %d
Bookmarks: %d
Hits: %d
Summary:
 %s`,
		f.Title,
		f.Author,
		f.Rating,
		f.ArchiveWarnings,
		f.Categories,
		strings.Join(f.Fandoms, ", "),
		strings.Join(f.Relationships, ", "),
		strings.Join(f.Characters, ", "),
		strings.Join(f.Tags, ", "),
		f.Language,
		f.PublishedAt.Format(time.DateOnly),
		status,
		f.Chapters,
		maxChapters,
		currentChapter,
		f.Comments,
		f.Kudos,
		f.Bookmarks,
		f.Hits,
		f.Summary,
	)
}
