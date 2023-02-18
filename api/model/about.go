package model

type UpdateAbout struct {
	FIO         string
	Description string
	Contact     string
}

type GetAbout struct {
	FIO         string
	Description string
	Contact     string
}

type AboutObject struct {
	FIO         string
	Description string
	Contact     string
}
