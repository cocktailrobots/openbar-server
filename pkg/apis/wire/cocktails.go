package wire

import (
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/db/cocktailsdb"
	"strings"
)

func ptr[T any](t T) *T {
	return &t
}

// Cocktail is a cocktail as it will be written to the wire in HTTP responses.
type Cocktail struct {
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Description *string `json:"description"`
}

// Cocktails is a slice of cocktails.
type Cocktails []Cocktail

// ToDbCocktails converts a list of Cocktails to a list of cocktailsdb.Cocktails.
func (c Cocktails) ToDbCocktails() []cocktailsdb.Cocktail {
	cocktails := make([]cocktailsdb.Cocktail, len(c))
	for i := range c {
		cocktails[i] = cocktailsdb.Cocktail{
			Name:        c[i].Name,
			DisplayName: c[i].DisplayName,
			Description: c[i].Description,
		}
	}

	return cocktails
}

func (c Cocktails) Validate() error {
	for i := range c {
		c[i].Name = strings.TrimSpace(c[i].Name)
		if len(c[i].Name) == 0 {
			return fmt.Errorf("cocktail name cannot be empty")
		}

		c[i].DisplayName = strings.TrimSpace(c[i].DisplayName)
		if len(c[i].DisplayName) == 0 {
			return fmt.Errorf("cocktail display name cannot be empty")
		}
	}

	return nil
}

// FromDbCocktails converts a list of cocktailsdb.Cocktails to a list of Cocktails.
func FromDbCocktails(cocktails []cocktailsdb.Cocktail) Cocktails {
	c := make(Cocktails, len(cocktails))
	for i := range cocktails {
		c[i] = Cocktail{
			Name:        cocktails[i].Name,
			DisplayName: cocktails[i].DisplayName,
			Description: cocktails[i].Description,
		}
	}

	return c
}
