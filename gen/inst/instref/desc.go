package instref

import (
	"encoding/xml"
	"os"
	"strings"
)

type NameDesc struct {
	Name, Desc string
}

func Descriptions() (map[string][]NameDesc, error) {
	desc := make(map[string][]NameDesc, 800)

	type IForm struct {
		Heading string `xml:"heading,attr"`
		Desc    string `xml:",chardata"`
	}

	type IForms struct {
		Title string  `xml:"title,attr"`
		List  []IForm `xml:"iform"`
	}

	type File struct {
		XMLName xml.Name `xml:"file"`
		IForms  IForms   `xml:"alphaindex>iforms"`
	}

	type All struct {
		XMLName xml.Name `xml:"allinstrs"`
		Files   []File   `xml:"file"`
	}

	// ISA_A64_xml_A_profile-2022-06 -> onebigfile.xml
	raw, err := os.ReadFile("gen/inst/instref/onebigfile.xml")
	if err != nil {
		return nil, err
	}

	var x All
	if err = xml.Unmarshal(raw, &x); err != nil {
		return nil, err
	}

	splitHeading := func(h string) []string {
		var list []string
		var current string
		var inParen bool
		for i := 0; i < len(h); i++ {
			switch c := h[i]; c {
			case ',':
				if !inParen {
					if current != "" {
						list, current = append(list, current), ""
					}
				} else {
					current += ","
				}
			case '(':
				current += "("
				inParen = true
			case ')':
				current += ")"
				inParen = false
			default:
				current += string(c)
			}
		}
		if current != "" {
			list = append(list, current)
		}
		return list
	}

	for _, f := range x.Files {
		if strings.Contains(f.IForms.Title, "SVE") {
			continue
		}
		for _, iform := range f.IForms.List {
			for _, name := range splitHeading(iform.Heading) {
				name = strings.TrimSpace(name)
				mnemonic := name
				if space := strings.IndexByte(mnemonic, ' '); space > 0 {
					mnemonic = mnemonic[:space]
				}
				if dot := strings.IndexByte(mnemonic, '.'); dot > 0 {
					mnemonic = mnemonic[:dot]
				}
				desc[mnemonic] = append(desc[mnemonic], NameDesc{name, iform.Desc})
			}
		}
	}

	for name, list := range desc {
		clean := make([]NameDesc, 0, len(list))
		dedupe := make(map[NameDesc]struct{}, len(list))
		for _, nameDesc := range list {
			if _, seen := dedupe[nameDesc]; !seen {
				clean = append(clean, nameDesc)
				dedupe[nameDesc] = struct{}{}
			}
		}
		desc[name] = clean
	}

	if len(desc) == 0 {
		panic("empty")
	}

	return desc, nil
}
