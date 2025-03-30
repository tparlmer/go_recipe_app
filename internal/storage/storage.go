// Interface definition for storage operations

package storage

import "go_recipe_app/internal/models"

// RecipeStore defines the interface for recipe storage
type RecipeStore interface {
	List() ([]models.Recipe, error)
	Get(id string) (models.Recipe, error)
	Create(recipe models.Recipe) error
	Update(recipe models.Recipe) error
	Delete(id string) error
}

