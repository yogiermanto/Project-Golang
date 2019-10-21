package model

type Customer struct {
	NoAntrian int
	Nama      string
	NoHp      int
	Status    bool
}

type Meja struct {
	IdMeja int
	Status bool
}

type ResAdmin struct {
	Meja            []Meja
	Total           int
	AntrianSekarang int
}
