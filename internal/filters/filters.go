package filters

import (
	"fmt"
	"strings"

	"github.com/voidwyrm-2/fik/internal/fic"
)

func Filter(fics []fic.Fic, filterStrs []string) (filtered []*fic.Fic, filtersUsed []string, err error) {
	if len(filterStrs) == 0 {
		filtered = make([]*fic.Fic, 0, len(fics))

		for i := range fics {
			filtered = append(filtered, &fics[i])
		}

		return
	}

	filters := []func(*fic.Fic) bool{}

	for _, str := range filterStrs {
		var name, arg string

		parts := strings.Split(str, ":")
		name = strings.ToLower(strings.TrimSpace(parts[0]))
		if len(parts) > 1 {
			arg = strings.TrimSpace(strings.Join(parts[1:], ":"))
		}

		switch name {
		case "favorites":
			filters = append(filters, func(f *fic.Fic) bool { return f.Favorite })
			filtersUsed = append(filtersUsed, name)
		case "rating":
			rating := fic.RatingFromString(arg)
			filters = append(filters, func(f *fic.Fic) bool { return f.Rating == rating })
			filtersUsed = append(filtersUsed, fmt.Sprintf("%s with arg '%s'", name, arg))
		case "author":
			filters = append(filters, func(f *fic.Fic) bool { return f.Author == arg })
			filtersUsed = append(filtersUsed, fmt.Sprintf("%s with arg '%s'", name, arg))
		default:
			return nil, nil, fmt.Errorf("Unknown filter '%s'", name)
		}
	}

	for i := range fics {
		fic := &fics[i]
		add := true

		for _, f := range filters {
			if !f(fic) {
				add = false
				break
			}
		}

		if add {
			filtered = append(filtered, fic)
		}
	}

	return
}
