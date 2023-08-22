package wire

import (
	"github.com/cocktailrobots/openbar-server/pkg/db"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMenus(t *testing.T) {
	menus := []*db.Menu{
		{
			Name: "Margaritas",
			Ingredients: []string{
				"tequila",
				"sweetened lime juice",
				"countreau", "sweetened dragonfruit and lime juice",
				"sweetened blood orange and lime juice",
			},
			RecipeIds: []string{"0", "1", "2", "3"},
		},
		{
			Name: "Negronis",
			Ingredients: []string{
				"gin",
				"sweet vermouth",
				"campari",
				"borollo chinato",
				"cold brew coffee",
			},
			RecipeIds: []string{"4", "5", "6", "7", "8"},
		},
	}

	menusWire := FromDbMenus(menus)
	menus2 := menusWire.ToDbMenus()
	require.Equal(t, menus, menus2)
}
