package echien

import (
	"errors"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var (
	re_tagWithClass = regexp.MustCompile(`^([a-z:-]+)\.(.+)$`)
	re_classOrId    = regexp.MustCompile(`^([.#])(.+)$`)
	re_tagOnly      = regexp.MustCompile(`^[a-z:-]+$`)
	re_fullMonty    = regexp.MustCompile(`^(\.)?([a-z:-]+)\.?([|,a-z:-]+)?@([a-z-]+)([=^$!]?)=(.+)$`)
)

type El struct {
	name  string
	attrs map[string]string
	class map[string]bool
	text  string
}

func (el *El) hasClass(name string) bool {
	_, hasClass := el.class[name]
	return hasClass
}

func (el *El) getAttribute(attr string) string {
	if attr == "text" {
		return el.text
	} else {
		value, hasAttr := el.attrs[attr]
		if hasAttr {
			return value
		}
	}
	return ""
}

type EChien struct {
	elements []*El
	body     *html.Tokenizer
}

func (e *EChien) Open(path string) (int, error) {

	m, _ := regexp.MatchString("^http", path)
	if m {
		res, err := http.Get(path)
		if err != nil {
			return 0, err
		}
		e.body = html.NewTokenizer(res.Body)
	} else {
		file, err := os.Open(path)
		if err != nil {
			return 0, err
		}
		e.body = html.NewTokenizer(file)
	}

	count := e.parse()

	if count > 0 {
		return count, nil
	} else {
		return 0, errors.New("No nodes parsed")
	}

}

func (e *EChien) Find(searchString string) (elementList []*El) {
	funcs, err := getFilters(searchString)
	if err == nil {
		for _, element := range e.elements {
			for i, test := range funcs {
				if !test(element) {
					break
				}
				if i == len(funcs)-1 {
					elementList = append(elementList, element)
				}
			}
		}
	}
	return elementList
}

func (e *EChien) Get() {

}

func (e *EChien) parse() int {

	re_head := regexp.MustCompile(`(?im)^\s*`)
	re_tail := regexp.MustCompile(`(?im)\s*$`)
	re_remove := regexp.MustCompile(`(?im)&nbsp;`)
	re_letters := regexp.MustCompile(`(?im)\w`)

	e.elements = []*El{}

	for {
		tokenType := e.body.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := e.body.Token()
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			tag := token.Data
			tag = strings.ToLower(tag)

			if tag == "b" || tag == "i" {
				continue
			}

			el := El{
				name:  tag,
				attrs: map[string]string{},
				class: map[string]bool{},
				text:  ""}

			for _, attr := range token.Attr {
				el.attrs[attr.Key] = attr.Val
			}

			classes, hasClassAtribute := el.attrs["class"]

			if hasClassAtribute {
				for _, c := range strings.Split(classes, " ") {
					el.class[c] = true
				}
			}
			e.elements = append(e.elements, &el)
		} else if tokenType == html.TextToken {
			text = re_remove.ReplaceAllString(text, "")
			text := re_head.ReplaceAllString(token.Data, "")
			text = re_tail.ReplaceAllString(text, "")
			if text != "" && re_letters.MatchString(text) {
				lastEl := e.elements[len(e.elements)-1]
				lastEl.text = text
			}
		}
	}

	return len(e.elements)

}

func alwaysFalse(el *El) bool {
	return false
}

func elementHasClassesAND(el *El, classes []string) bool {
	for _, class := range classes {
		if !el.hasClass(class) {
			return false
		}
	}
	return true
}

func elementHasClassesOR(el *El, classes []string) bool {
	for _, class := range classes {
		if el.hasClass(class) {
			return true
		}
	}
	return false
}

func tagF(tag string) func(el *El) bool {
	return func(el *El) bool {
		return el.name == tag
	}
}

func classF(class string) func(el *El) bool {
	return func(el *El) bool {
		return el.hasClass(class)
	}
}

func attributeEqual(attrib string, value string) func(el *El) bool {
	return func(el *El) bool {
		v := el.getAttribute(attrib)
		return (v == value)
	}
}

func attributeStartsWith(attrib string, searchString string) func(el *El) bool {
	searchString = "(?im)^" + searchString
	re, err := regexp.Compile(searchString)
	if err != nil {
		return alwaysFalse
	}
	return func(el *El) bool {
		v := el.getAttribute(attrib)
		if v == "" {
			return false
		}
		m := re.FindAllStringSubmatch(v, -1)
		return (m != nil)
	}
}

func attributeEndsWith(attrib string, searchString string) func(el *El) bool {
	searchString = "(?im)" + searchString + "$"
	re, err := regexp.Compile(searchString)
	if err != nil {
		return alwaysFalse
	}
	return func(el *El) bool {
		v := el.getAttribute(attrib)
		if v == "" {
			return false
		}
		m := re.FindAllStringSubmatch(v, -1)
		return (m != nil)
	}
}

func attributeRegex(attrib string, searchString string) func(el *El) bool {
	searchString = "(?im)" + searchString
	re, err := regexp.Compile(searchString)
	if err != nil {
		return alwaysFalse
	}
	return func(el *El) bool {
		v := el.getAttribute(attrib)
		if v == "" {
			return false
		}
		m := re.FindAllStringSubmatch(v, -1)
		return (m != nil)
	}
}

func attributeDoesNotContain(attrib string, searchString string) func(el *El) bool {
	searchString = "(?im)" + searchString
	re, err := regexp.Compile(searchString)
	if err != nil {
		return alwaysFalse
	}
	return func(el *El) bool {
		v := el.getAttribute(attrib)
		if v == "" {
			return true
		}
		m := re.FindAllStringSubmatch(v, -1)
		return (m == nil)
	}
}

func idF(id string) func(el *El) bool {
	return func(el *El) bool {
		v := el.getAttribute("id")
		return (v == id)
	}
}

func classesMultiple(classes string) func(el *El) bool {
	switch {
	case strings.ContainsAny(classes, "|"):
		classlist := strings.Split(classes, "|")
		return func(el *El) bool {
			return elementHasClassesOR(el, classlist)
		}
	case strings.ContainsAny(classes, ","):
		classlist := strings.Split(classes, ",")
		return func(el *El) bool {
			return elementHasClassesAND(el, classlist)
		}
	default:
		return func(el *El) bool {
			return el.hasClass(classes)
		}
	}
}

const (
	TAGONLY      = 1
	TAGWITHCLASS = 2
	CLASSORID    = 3
	FULLMONTY    = 4
)

type Classification struct {
	classification                                     int
	tag, id, class, attribute, matchtype, searchstring string
}

func decompose(search string) (*Classification, error) {
	if m := re_fullMonty.FindAllStringSubmatch(search, -1); m != nil {
		sign := m[0][1]
		identifier := m[0][2]
		classes := m[0][3]
		attribute := m[0][4]
		matchType := m[0][5]
		searchString := m[0][6]
		r := Classification{
			classification: FULLMONTY,
			class:          classes,
			attribute:      attribute,
			matchtype:      matchType,
			searchstring:   searchString}
		switch sign {
		case ".":
			r.class = identifier
		default:
			r.tag = identifier
		}
		return &r, nil
	} else if m = re_classOrId.FindAllStringSubmatch(search, -1); m != nil {
		sign := m[0][1] // Either "." or "#"
		identifier := m[0][2]
		r := Classification{classification: CLASSORID}
		switch sign {
		case ".":
			r.class = identifier
		case "#":
			r.id = identifier
		}
		return &r, nil
	} else if m = re_tagWithClass.FindAllStringSubmatch(search, -1); m != nil {
		tag := m[0][1]
		classes := m[0][2]
		r := Classification{
			classification: TAGWITHCLASS,
			tag:            tag,
			class:          classes}
		return &r, nil
	} else if m = re_tagOnly.FindAllStringSubmatch(search, -1); m != nil {
		r := Classification{
			classification: TAGONLY,
			tag:            m[0][0]}
		return &r, nil
	}
	return &Classification{}, errors.New("")
}

func getFilters(search string) ([]func(*El) bool, error) {

	funcs := []func(*El) bool{}

	r, err := decompose(search)

	if err == nil {
		switch r.classification {
		case TAGONLY:
			funcs = append(funcs, tagF(r.tag))
		case TAGWITHCLASS:
			funcs = append(funcs, tagF(r.tag))
			funcs = append(funcs, classesMultiple(r.class))
		case CLASSORID:
			if r.class != "" {
				funcs = append(funcs, classesMultiple(r.class))
			} else if r.id != "" {
				funcs = append(funcs, idF(r.id))
			}
		case FULLMONTY:
			if r.tag != "" {
				funcs = append(funcs, tagF(r.tag))
			}
			if r.class != "" {
				funcs = append(funcs, classesMultiple(r.class))
			}
			switch r.matchtype {
			case "":
				funcs = append(funcs, attributeRegex(r.attribute, r.searchstring))
			case "=": // Regex match
				funcs = append(funcs, attributeEqual(r.attribute, r.searchstring))
			case "^": // Starts with
				funcs = append(funcs, attributeStartsWith(r.attribute, r.searchstring))
			case "$": // Ends with
				funcs = append(funcs, attributeEndsWith(r.attribute, r.searchstring))
			case "!": // Does not contain
				funcs = append(funcs, attributeDoesNotContain(r.attribute, r.searchstring))
			}
		}
		if len(funcs) > 0 {
			return funcs, nil
		} else {
			return funcs, errors.New("No filter functions")
		}
	}

	return funcs, err

}

/*

































*/
