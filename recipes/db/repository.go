package db

import (
	"time"

	"github.com/go-kit/log"
)

// ---------
// RECIPE MODEL
// ---------

// Recipe represents a recipe in the system
// TODO: Move recipe model to it's own model.go file
type Ingredient struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Unit     string  `json:"unit"`
	Position int     `json:"position"` // For ordering ingredients
}

type Instruction struct {
	ID       string `json:"id"`
	Step     string `json:"step"`
	Position int    `json:"position"` // For ordering steps
}

type Recipe struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	PrepTime     time.Duration `json:"prep_time"`
	CookTime     time.Duration `json:"cook_time"`
	Servings     int32         `json:"servings"`
	Ingredients  []Ingredient  `json:"ingredients"`
	Instructions []Instruction `json:"instructions"`
	OwnerID      string        `json:"ownerid"`  // Used to associate users with recipes
	IsPublic     bool          `json:"ispublic"` // Used to associate users with recipes
}

// ----------
// DATA ACCESS INTERFACE
// ----------

// RecipeRepository defines the interface for recipe repository
type RecipeRepository interface {
	ListRecipes(logger log.Logger) ([]Recipe, error)
	ListUserRecipes(userID string, logger log.Logger) ([]Recipe, error)
	GetRecipe(id string, logger log.Logger) (*Recipe, error)
	CreateRecipe(recipe *Recipe, logger log.Logger) error
	UpdateRecipe(recipe Recipe, logger log.Logger) error
	DeleteRecipe(id string, logger log.Logger) error
	Close(logger log.Logger) error
}
