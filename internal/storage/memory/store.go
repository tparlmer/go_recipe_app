// in-memory storage implementation

package memory

import (
	"fmt"
	"go_recipe_app/internal/models"
	"sync"
)

// Store implements storage.RecipeStore interface
type Store struct {
	mu      sync.RWMutex // For safe concurrent access
	recipes map[string]models.Recipe
}

// New creates a new in-memory store
func New() *Store {
	return &Store{
		recipes: make(map[string]models.Recipe),
	}
}

//List returns all recipes
func (s *Store) List() ([]models.Recipe, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	recipes := make([]models.Recipe, 0, len(s.recipes))
	for _, recipe := range s.recipes {
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

// Get returns a single recipe by ID
func (s *Store) Get(id string) (models.Recipe, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	recipe, exists := s.recipes[id]
	if !exists {
		return models.Recipe{}, fmt.Errorf("recipe not found: %s", id)
	}
	return recipe, nil
}

// Create adds a new recipe
func (s *Store) Create(recipe models.Recipe) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.recipes[recipe.ID]; exists {
		return fmt.Errorf("recipe already exists: %s", recipe.ID)
	}

	s.recipes[recipe.ID] = recipe
	return nil
}

// Update modifies an existing recipe
func (s *Store) Update(recipe models.Recipe) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.recipes[recipe.ID]; !exists {
		return fmt.Errorf("recipe not found: %s", recipe.ID)
	}

	s.recipes[recipe.ID] = recipe
	return nil
}

// Delete removes a recipe
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.recipes[id]; !exists {
		return fmt.Errorf("recipe not found: %s", id)
	}

	delete(s.recipes, id)
	return nil
}
