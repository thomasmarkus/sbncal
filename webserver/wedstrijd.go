package main

type Wedstrijd struct {
	Datum           string   `json:"Datum"`
	Plaats          string   `json:"Plaats"`
	LSR             bool     `json:"LSR"`
	MSR             bool     `json:"MSR"`
	KSR             bool     `json:"KSR"`
	Jeugd           bool     `json:"Jeugd"`
	BSR             bool     `json:"BSR"`
	Kwalificatie    bool     `json:"Kwalificatie"`
	Afstanden       []string `json:"Afstanden"`
	MinLeeftijd     string   `json:"min.leeftijd"`
	Organisatorlink []string `json:"Organisator"`
	Inschrijflink   []string `json:"Inschrijflink"`
	Uitslaglink     []string `json:"Uitslag"`
	ONK             bool     `json:"ONK"`
	Koppel          bool     `json:"Koppel"`
	Afgelast        bool     `json:"Afgelast"`
	Coords          []string `json:"coords"`
}
