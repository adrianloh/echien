package echien

import "testing"

var text = "likes it up the ass"

var el = El{
	name: "img",
	attrs: map[string]string{
		"id":       "niki",
		"position": "doggystyle",
		"src":      "http://dummy/image.jpg"},
	class: map[string]bool{
		"cock":   true,
		"sucker": true},
	text: text}

var searches = map[string]Classification{
	"img":                                     Classification{classification: TAGONLY, tag: "img"},
	"img.banner":                              Classification{classification: TAGWITHCLASS, tag: "img", class: "banner"},
	"img.banner,big":                          Classification{classification: TAGWITHCLASS, tag: "img", class: "banner,big"},
	"img.banner|thumbnail":                    Classification{classification: TAGWITHCLASS, tag: "img", class: "banner|thumbnail"},
	".banner":                                 Classification{classification: CLASSORID, id: "", class: "banner"},
	"#superific":                              Classification{classification: CLASSORID, id: "superific", class: ""},
	"img@attribute=protopolis":                Classification{classification: FULLMONTY, tag: "img", class: "", attribute: "attribute", matchtype: "", searchstring: "protopolis"},
	"img.killjoy@attribute=protopolis":        Classification{classification: FULLMONTY, tag: "img", class: "killjoy", attribute: "attribute", matchtype: "", searchstring: "protopolis"},
	"img.killjoy,elroy@attribute=protopolis":  Classification{classification: FULLMONTY, tag: "img", class: "killjoy,elroy", attribute: "attribute", matchtype: "", searchstring: "protopolis"},
	"img.killjoy|elroy@attribute==protopolis": Classification{classification: FULLMONTY, tag: "img", class: "killjoy|elroy", attribute: "attribute", matchtype: "=", searchstring: "protopolis"},
	"img.killjoy|elroy@attribute^=protopolis": Classification{classification: FULLMONTY, tag: "img", class: "killjoy|elroy", attribute: "attribute", matchtype: "^", searchstring: "protopolis"},
	"img.killjoy|elroy@attribute$=protopolis": Classification{classification: FULLMONTY, tag: "img", class: "killjoy|elroy", attribute: "attribute", matchtype: "$", searchstring: "protopolis"},
	"img.killjoy|elroy@attribute!=protopolis": Classification{classification: FULLMONTY, tag: "img", class: "killjoy|elroy", attribute: "attribute", matchtype: "!", searchstring: "protopolis"},
	".banner@attribute=protopolis":            Classification{classification: FULLMONTY, tag: "", class: "banner", attribute: "attribute", matchtype: "", searchstring: "protopolis"}}

func Test_Has_Class(t *testing.T) {

	if el.hasClass("cock") != true {
		t.Errorf("Expected class not found")
	}

	if el.hasClass("bullshit") != false {
		t.Errorf("Found unexpected class")
	}

}

func Test_Has_Attribute(t *testing.T) {

	if el.GetAttribute("text") != text {
		t.Errorf("Expected innerText not found")
	}

	if el.GetAttribute("bullshit") != "" {
		t.Errorf("Found unexpected innerText: %s", el.GetAttribute("bullshit"))
	}

	if el.GetAttribute("position") != "doggystyle" {
		t.Errorf("Expected to get attribute")
	}

}

func Test_Always_False(t *testing.T) {

	if alwaysFalse(&el) != false {
		t.Errorf("Always false didn't return false")
	}

}

func Test_Has_Class_BOOLEAN_AND(t *testing.T) {

	if elementHasClassesAND(&el, []string{"cock", "sucker"}) != true {
		t.Errorf("Element has two classes")
	}

	if elementHasClassesAND(&el, []string{"cock", "rabbbit"}) != false {
		t.Errorf("Element does NOT have two classes")
	}

}

func Test_Has_Class_BOOLEAN_OR(t *testing.T) {

	if elementHasClassesOR(&el, []string{"cock", "nohole"}) != true {
		t.Errorf("Element has two classes")
	}

	if elementHasClassesOR(&el, []string{"fucker", "nohole"}) != false {
		t.Errorf("Element does NOT have two classes")
	}

}

func Test_Match_Tag_Function(t *testing.T) {

	f := tagF("img")

	if f(&el) != true {
		t.Errorf("Expected to match tag")
	}

	f = tagF("a")

	if f(&el) != false {
		t.Errorf("Did NOT expect to match tag")
	}

}

func Test_Match_Class_Function(t *testing.T) {

	f := classF("cock")

	if f(&el) != true {
		t.Errorf("Expected to match class")
	}

	f = classF("nohole")

	if f(&el) != false {
		t.Errorf("Did NOT expect to match class")
	}

}

func Test_Match_ID_Function(t *testing.T) {

	f := idF("bejesus")

	if f(&el) != false {
		t.Errorf("Did NOT expect to match id")
	}

	f = idF("niki")

	if f(&el) != true {
		t.Errorf("Expected to match id")
	}

}

func Test_Attribute_Equal(t *testing.T) {

	f := attributeEqual("dude", "ranch")

	if f(&el) != false {
		t.Errorf("Did not expect to find attribute")
	}

	f = attributeEqual("position", "doggystyle")

	if f(&el) != true {
		t.Errorf("Expected to match attribute")
	}

}

func Test_Attribute_Starts_With(t *testing.T) {

	f := attributeStartsWith("dude", "ranch")

	if f(&el) != false {
		t.Errorf("Did not expect to find attribute")
	}

	f = attributeStartsWith("src", "https")

	if f(&el) != false {
		t.Errorf("Did not expect to match attribute")
	}

	f = attributeStartsWith("src", "http")

	if f(&el) != true {
		t.Errorf("Expected to match attribute")
	}

}

func Test_Attribute_Ends_With(t *testing.T) {

	f := attributeEndsWith("dude", "ranch")

	if f(&el) != false {
		t.Errorf("Did not expect to find attribute")
	}

	f = attributeEndsWith("src", "png")

	if f(&el) != false {
		t.Errorf("Did not expect to match attribute")
	}

	f = attributeEndsWith("src", "jpg")

	if f(&el) != true {
		t.Errorf("Expected to match attribute")
	}

}

func Test_Attribute_Regex(t *testing.T) {

	f := attributeRegex("dude", "ranch")

	if f(&el) != false {
		t.Errorf("Did not expect to find attribute")
	}

	f = attributeRegex("src", "^[a-z|][]")

	if f(&el) != false {
		t.Errorf("Expected bad regex to fail")
	}

	f = attributeRegex("src", `\/image\.(png|bmp)$`)

	if f(&el) != false {
		t.Errorf("Did not expect regex to match attribute")
	}

	f = attributeRegex("src", `^http.+\/image\.(png|bmp|jpg)$`)

	if f(&el) != true {
		t.Errorf("Expected regex to match attribute")
	}

}

func Test_Attribute_Does_Not_Contain(t *testing.T) {

	f := attributeDoesNotContain("dude", "ranch")

	if f(&el) != true {
		t.Errorf("Did not expect to find attribute")
	}

	f = attributeDoesNotContain("src", "^[a-z|][]")

	if f(&el) != false {
		t.Errorf("Expected bad regex to fail")
	}

	f = attributeDoesNotContain("src", `\/image\.(png|bmp)$`)

	if f(&el) != true {
		t.Errorf("Expected regex to NOT match attribute")
	}

	f = attributeDoesNotContain("src", `^http.+\/image\.(png|bmp|jpg)$`)

	if f(&el) != false {
		t.Errorf("Expected regex to NOT match attribute")
	}

}

func Test_Multiple_Classes_String(t *testing.T) {

	f := classesMultiple(`cock`)

	if f(&el) != true {
		t.Errorf("Expected to match one class")
	}

	f = classesMultiple(`bingo`)

	if f(&el) != false {
		t.Errorf("Did not expect to match one class")
	}

	f = classesMultiple(`cock,sucker`)

	if f(&el) != true {
		t.Errorf("Expected to match multiple AND classes")
	}

	f = classesMultiple(`cock,donut`)

	if f(&el) != false {
		t.Errorf("Did not expect to match multiple AND classes")
	}

	f = classesMultiple(`cock|donut`)

	if f(&el) != true {
		t.Errorf("Expected to match multiple OR classes")
	}

	f = classesMultiple(`fucker|donut`)

	if f(&el) != false {
		t.Errorf("Did not expect to match multiple OR classes")
	}

}

func Test_Decompose(t *testing.T) {

	for str, D := range searches {
		r, err := decompose(str)
		if err != nil || r.classification != D.classification {
			t.Errorf(`%s | expected %d got %d`, str, D.classification, r.classification)
		} else {
			switch r.classification {
			case TAGONLY:
				if r.tag != D.tag {
					t.Errorf(`%d | expected "%s" got "%s"`, D.classification, D.tag, r.tag)
				}
			case TAGWITHCLASS:
				if r.tag != D.tag {
					t.Errorf(`TAGWITHCLASS | expected "%s" got "%s"`, D.tag, r.tag)
				}
				if r.class != D.class {
					t.Errorf(`TAGWITHCLASS | expected "%s" got "%s"`, D.class, r.class)
				}
			case CLASSORID:
				if r.class != D.class {
					t.Errorf(`CLASSORID | expected "%s" got "%s"`, D.class, r.class)
				}
				if r.id != D.id {
					t.Errorf(`CLASSORID | expected "%s" got "%s"`, D.id, r.id)
				}
			case FULLMONTY:
				if r.tag != D.tag {
					t.Errorf(`FULLMONTY | expected "%s" got "%s"`, D.tag, r.tag)
				}
				if r.class != D.class {
					t.Errorf(`FULLMONTY | expected "%s" got "%s"`, D.class, r.class)
				}
				if r.attribute != D.attribute {
					t.Errorf(`FULLMONTY | expected "%s" got "%s"`, D.attribute, r.attribute)
				}
				if r.matchtype != D.matchtype {
					t.Errorf(`FULLMONTY | expected "%s" got "%s"`, D.matchtype, r.matchtype)
				}
				if r.searchstring != D.searchstring {
					t.Errorf(`FULLMONTY | expected "%s" got "%s"`, D.searchstring, r.searchstring)
				}
			}
		}
	}

}

func Test_Local_Document(t *testing.T) {

	e, err := Open("https://raw.githubusercontent.com/adrianloh/echien/master/test.html")

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	searches := map[string]int{
		`img`:                 3,
		`link`:                7,
		`.headerImg`:          3,
		`div.headerImg,small`: 1,
		`.headerImg,big`:      2,
		`.headerImg,small`:    1,
		`.big|small`:          3,
		`script@src==ajax`:    0,
		`script@src=ajax`:     1,
		`title@text=mo`:       1,
		`a@text==Combined`:    1,
		`img@src^=https`:      3,
		`img@src$=jpg`:        1,
		`link@href!=^https`:   4,
		`img@src!=png`:        1,
		`div@class!=big`:      1}

	for search, expected := range searches {
		els := e.Find(search)
		if len(els) != expected {
			t.Errorf(`"%s" | Expected %d got %d`, search, expected, len(els))
		}
	}

}

func Test_Avaxhome(t *testing.T) {

	e, err := Open("http://avaxhm.com/girls/Nicki-Minaj-Roberto-Cavalli-Spring-Summer-2015.html")

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	searches := map[string]int{
		`a@href=nitroflare`: 1,
		`a@text=nitroflare`: 1}

	for search, expected := range searches {
		els := e.Find(search)
		if len(els) != expected {
			t.Errorf(`"%s" | Expected %d got %d`, search, expected, len(els))
		}
	}

}

/*

















































*/
