package instref

import (
	"encoding/xml"
	"errors"
	"os"
	"strings"
)

const AliasPrefix = "an alias of "

type NameDesc struct {
	Name, Desc string
	AliasOf    []AliasName
}

type AliasName struct {
	Mnemonic string
	Name     string
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

	for _, f := range x.Files {
		if title := f.IForms.Title; strings.Contains(title, "SVE") || strings.Contains(title, "SME") {
			continue
		}
		for _, iform := range f.IForms.List {
			for _, name := range splitNameList(iform.Heading) {
				mnemonic := getMnemonic(name)
				nameDesc := NameDesc{
					Name: name,
					Desc: iform.Desc,
				}
				if aliasStart := strings.Index(iform.Desc, AliasPrefix); aliasStart >= 0 {
					nameDesc.Desc = iform.Desc[:aliasStart]
					aliasNamesStart := aliasStart + len(AliasPrefix)
					aliasNames := strings.TrimRight(iform.Desc[aliasNamesStart:], ".")
					if len(aliasNames) == 0 {
						return nil, errors.New("Empty alias names for " + name)
					}
					aliases := splitNameList(aliasNames)
					for _, other := range aliases {
						if other == "" {
							return nil, errors.New("Empty alias name listed for " + name)
						}
						nameDesc.AliasOf = append(nameDesc.AliasOf, AliasName{
							Mnemonic: getMnemonic(other),
							Name:     strings.TrimSpace(other),
						})
					}
				}
				desc[mnemonic] = append(desc[mnemonic], nameDesc)
			}
		}
	}

	for name, list := range desc {
		clean := make([]NameDesc, 0, len(list))
		dedupe := make(map[[2]string]struct{}, len(list))
		for _, nameDesc := range list {
			kv := [2]string{nameDesc.Name, nameDesc.Desc}
			if _, seen := dedupe[kv]; !seen {
				clean = append(clean, nameDesc)
				dedupe[kv] = struct{}{}
			}
		}
		desc[name] = clean
	}

	if len(desc) == 0 {
		panic("empty")
	}

	return desc, nil
}

func getMnemonic(name string) string {
	if space := strings.IndexByte(name, ' '); space > 0 {
		name = name[:space]
	}
	if dot := strings.IndexByte(name, '.'); dot > 0 {
		name = name[:dot]
	}
	return strings.TrimSpace(strings.ToUpper(name))
}

func splitNameList(names string) []string {
	var list []string
	var current string
	var inParen bool
	for i := 0; i < len(names); i++ {
		switch c := names[i]; c {
		case ',':
			if !inParen {
				if current != "" {
					list, current = append(list, strings.TrimSpace(current)), ""
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
		list = append(list, strings.TrimSpace(current))
	}
	return list
}
