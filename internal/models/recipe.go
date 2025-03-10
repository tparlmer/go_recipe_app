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
type Recipe struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description,omitempty"` // use omitempty when field is optional
	PrepTime     time.Duration `json:"prep_time,omitempty"`
	CookTime     time.Duration `json:"cook_time,omitempty"`
	Servings     int32         `json:"servings,omitempty"`
	Ingredients  []Ingredient  `json:"ingredients"`
	Instructions []Instruction `json:"instructions"`
}

type Ingredient struct {
	Quantity float32 `json:"quantity"`
	Unit     string  `json:"unit"`
	Name     string  `json:"name"`
}

type Instruction struct {
	Step             uint   `json:"step"`
	StepInstructions string `json:"instructions"`
}
