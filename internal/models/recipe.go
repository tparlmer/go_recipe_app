// Recipe struct and its methods

package models

import (
	"time"
)

// To store our Recipe struct we need to serialize it into json
// Go structs are typed collections of fields - they are useful for grouping data together to form records
// below add fields with types
// Use PascalCase for exported fiels that should be accessible outside the package
// Use camelCase for unexported fields accesible within the package

// Ingredient represents a recipe ingredient
type Ingredient struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Unit     string  `json:"unit"`
	Position int     `json:"position"` // For ordering ingredients
}

// Instruction represents a recipe step
type Instruction struct {
	ID       string `json:"id"`
	Step     string `json:"step"`
	Position int    `json:"position"` // For ordering steps
}

// Update Recipe struct to include ingredients and instructions
type Recipe struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	PrepTime     time.Duration `json:"prep_time"`
	CookTime     time.Duration `json:"cook_time"`
	Servings     int32         `json:"servings"`
	Ingredients  []Ingredient  `json:"ingredients"`
	Instructions []Instruction `json:"instructions"`
}
