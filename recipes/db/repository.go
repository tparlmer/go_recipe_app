package db

import (
	"time"
)

// ---------
// RECIPE MODEL
// ---------

// Recipe represents a recipe in the system
// TODO: Move recipe model to it's own model.go file
type Ingredient struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Amount float64 `json:"amount"`
	Unit string `json:"unit"`
	Position int `json:"position"` // For ordering ingredients
}

type Instruction struct {
	ID string `json:"id"`
	Step string `json:"step"`
	Position int `json:"position"` // For ordering steps
}

type Recipe struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	PrepTime time.Duration `json:"prep_time"`
	CookTime time.Duration `json:"cook_time"`
	Servings int32 `json:"servings"`
	Ingredients []Ingredient `json:"ingredients"`
	Instructions []Instruction `json:"instructions"`
}

// ----------
// DATA ACCESS INTERFACE
// ----------

// RecipeStore defines the interface for recipe storage
type RecipeStore interface {
	List() ([]Recipe, error)
	Get(id string) (Recipe, error)
	Create(recipe Recipe) error
	Update(recipe Recipe) error
	Delete(id string) error
}