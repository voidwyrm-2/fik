package fic

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	pairPat = regexp.MustCompile(`^([0-9]{1,15})(?:,([0-9]{1,15}))?`)
	linkPat = regexp.MustCompile(`^(?:https://)?(?:archiveofourown\.org/works/)([0-9]{1,15})(?:/chapters/([0-9]{1,15})(?:#workskin)?)?(?:\?[a-z_]+=[a-z]+)*`)
)

type Id uint32

func ParseId(str string) (Id, error) {
	n, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("'%s' is not a valid AO3 fiction ID", str)
	}

	return Id(n), nil
}

func ParseIds(idStr, chapterIdStr string) (id, chapterId Id, err error) {
	id, err = ParseId(idStr)
	if err != nil {
		err = fmt.Errorf("'%s' is not a valid AO3 fiction ID", idStr)
		return
	}

	if len(chapterIdStr) > 0 {
		chapterId, err = ParseId(chapterIdStr)
		if err != nil {
			err = fmt.Errorf("'%s' is not a valid AO3 fiction chapter ID", chapterIdStr)
			return
		}
	}

	return
}

func ParseFicEntry(str string) (Id, Id, error) {
	var idStr, chapterIdStr string

	if m := pairPat.FindStringSubmatch(str); m != nil {
		idStr = m[1]
		chapterIdStr = m[2]
	} else if m := linkPat.FindStringSubmatch(str); m != nil {
		idStr = m[1]
		chapterIdStr = m[2]
	} else {
		fmt.Println(m, len(str))
		return 0, 0, fmt.Errorf("'%s' is neither a valid ID pair nor a valid AO3 link", str)
	}

	return ParseIds(idStr, chapterIdStr)
}
