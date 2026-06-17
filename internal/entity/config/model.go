package config

import (
	"delicias-da-lu-service.com/mod/internal/entity/cakebuilder"
	"delicias-da-lu-service.com/mod/internal/entity/contact"
	"delicias-da-lu-service.com/mod/internal/entity/menu"
)

type SiteConfig struct {
	Home        *HomeConfig       `json:"home,omitempty" firestore:"home"`
	Menu        MenuConfig        `json:"menu" firestore:"menu"`
	CakeBuilder CakeBuilderConfig `json:"cakeBuilder" firestore:"cakeBuilder"`
	Contacts    contact.Contact   `json:"contacts" firestore:"contacts"`
}

type MenuConfig struct {
	Items          []menu.MenuItem   `json:"items" firestore:"items"`
	SectionLabels  map[string]string `json:"sectionLabels,omitempty" firestore:"sectionLabels"`
	CustomSections []string          `json:"customSections,omitempty" firestore:"customSections"`
}

type CakeBuilderConfig struct {
	Massas     []cakebuilder.CakeBuilderComponent `json:"massas" firestore:"massas"`
	Recheios   []cakebuilder.CakeBuilderComponent `json:"recheios" firestore:"recheios"`
	Coberturas []cakebuilder.CakeBuilderComponent `json:"coberturas" firestore:"coberturas"`
	Decoracoes []cakebuilder.CakeBuilderComponent `json:"decoracoes" firestore:"decoracoes"`
}
