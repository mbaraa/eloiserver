package models

type Overlay struct {
	Owner struct {
		Email string `json:"email" xml:"email"`
		Name  string `json:"name" xml:"name"`
		Type  string `json:"type" xml:"type,attr"`
	} `json:"owner" xml:"owner"`

	Source []struct {
		Type string `json:"type" xml:"type,attr"`
		Link string `json:"link" xml:",innerxml"`
	} `json:"source" xml:"source"`

	EbuildGroups map[string]*EbuildGroup `json:"ebuildGroups"`

	Name        string `json:"name" xml:"name"`
	Description string `json:"description" xml:"description"`
	Homepage    string `json:"homepage" xml:"homepage"`
	Feed        string `json:"feed" xml:"feed"`
}
type EbuildGroup struct {
	Ebuilds map[string]*Ebuild `json:"ebuilds"`
	Name    string             `json:"name"`
	Link    string             `json:"link"`
}

type Ebuild struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Homepage     string `json:"homepage"`
	Flags        string `json:"flags"`
	License      string `json:"license"`
	OverlayName  string `json:"overlayName"`
	GroupName    string `json:"groupName"`
	Architecture string `json:"architecture"`
	Description  string `json:"description"`
}
