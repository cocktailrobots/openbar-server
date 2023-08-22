package wire

import (
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/db/openbardb"
	"strings"
)

type Menu struct {
	Name        string   `json:"name"`
	Ingredients []string `json:"ingredients"`
	RecipeIds   []string `json:"items"`
}

type Menus []Menu

func (ms Menus) ToDbMenus() []*openbardb.Menu {
	menus := make([]*openbardb.Menu, len(ms))
	for i := range ms {
		menus[i] = &openbardb.Menu{
			Name:        ms[i].Name,
			Ingredients: ms[i].Ingredients,
			RecipeIds:   ms[i].RecipeIds,
		}
	}

	return menus
}

func (ms Menus) Validate() error {
	for i := range ms {
		ms[i].Name = strings.TrimSpace(ms[i].Name)
		if len(ms[i].Name) == 0 {
			return fmt.Errorf("menu name cannot be empty")
		}
	}

	return nil
}

func FromDbMenus(menus []*openbardb.Menu) Menus {
	ms := make(Menus, len(menus))
	for i := range menus {
		ms[i] = Menu{
			Name:        menus[i].Name,
			Ingredients: menus[i].Ingredients,
			RecipeIds:   menus[i].RecipeIds,
		}
	}

	return ms
}
