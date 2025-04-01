package boltdb

// To Run the tests - go test ./internal/storage/boltdb/...
// ... is a go specific wildcard operator that means "test this package and all subpackages"

import (
	"go_recipe_app/internal/models"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) (*Store, string) {
	// Create a temporary directory for the test database
	tempDir, err := os.MkdirTemp("", "recipe-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tempDir, "test.db")
	store, err := New(dbPath)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create test store: %v", err)
	}

	return store, tempDir
}

func cleanupTestDB(store *Store, tempDir string) {
	store.Close()
	os.RemoveAll(tempDir)
}

func createTestRecipe() models.Recipe {
	return models.Recipe{
		ID:          "test-recipe-1",
		Title:       "Test Recipe",
		Description: "A test recipe",
		PrepTime:    15 * time.Minute,
		CookTime:    30 * time.Minute,
		Servings:    4,
		Ingredients: []models.Ingredient{
			{
				ID:       "ing-1",
				Name:     "Test Ingredient",
				Amount:   2,
				Unit:     "cups",
				Position: 0,
			},
		},
		Instructions: []models.Instruction{
			{
				ID:       "step-1",
				Step:     "Test Step",
				Position: 0,
			},
		},
	}
}

func TestCreateAndGet(t *testing.T) {
	store, tempDir := setupTestDB(t)
	defer cleanupTestDB(store, tempDir)

	recipe := createTestRecipe()

	// Test Create
	err := store.Create(recipe)
	if err != nil {
		t.Errorf("Failed to create recipe: %v", err)
	}

	// Test Get
	retrieved, err := store.Get(recipe.ID)
	if err != nil {
		t.Errorf("Failed to get recipe: %v", err)
	}

	if retrieved.ID != recipe.ID {
		t.Errorf("Got wrong recipe. Want %s, got %s", recipe.ID, retrieved.ID)
	}
	if retrieved.Title != recipe.Title {
		t.Errorf("Got wrong title. Want %s, got %s", recipe.Title, retrieved.Title)
	}
}

func TestList(t *testing.T) {
	store, tempDir := setupTestDB(t)
	defer cleanupTestDB(store, tempDir)

	// Create multiple recipes
	recipes := []models.Recipe{
		createTestRecipe(),
		{
			ID:          "test-recipe-2",
			Title:       "Another Test Recipe",
			Description: "Another test recipe",
			PrepTime:    20 * time.Minute,
			CookTime:    45 * time.Minute,
			Servings:    6,
		},
	}

	for _, recipe := range recipes {
		if err := store.Create(recipe); err != nil {
			t.Fatalf("Failed to create recipe: %v", err)
		}
	}

	// Test List
	listed, err := store.List()
	if err != nil {
		t.Errorf("Failed to list recipes: %v", err)
	}

	if len(listed) != len(recipes) {
		t.Errorf("Wrong number of recipes. Want %d, got %d", len(recipes), len(listed))
	}
}

func TestUpdate(t *testing.T) {
	store, tempDir := setupTestDB(t)
	defer cleanupTestDB(store, tempDir)

	recipe := createTestRecipe()

	// Create initial recipe
	if err := store.Create(recipe); err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	// Modify recipe
	recipe.Title = "Updated Test Recipe"
	recipe.Description = "Updated description"

	// Test Update
	err := store.Update(recipe)
	if err != nil {
		t.Errorf("Failed to update recipe: %v", err)
	}

	// Verify update
	updated, err := store.Get(recipe.ID)
	if err != nil {
		t.Errorf("Failed to get updated recipe: %v", err)
	}

	if updated.Title != recipe.Title {
		t.Errorf("Update failed. Want title %s, got %s", recipe.Title, updated.Title)
	}
}

func TestDelete(t *testing.T) {
	store, tempDir := setupTestDB(t)
	defer cleanupTestDB(store, tempDir)

	recipe := createTestRecipe()

	// Create recipe
	if err := store.Create(recipe); err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	// Test Delete
	err := store.Delete(recipe.ID)
	if err != nil {
		t.Errorf("Failed to delete recipe: %v", err)
	}

	// Verify deletion
	_, err = store.Get(recipe.ID)
	if err == nil {
		t.Error("Recipe still exists after deletion")
	}
}

func TestErrorCases(t *testing.T) {
	store, tempDir := setupTestDB(t)
	defer cleanupTestDB(store, tempDir)

	// Test getting non-existent recipe
	_, err := store.Get("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent recipe")
	}

	// Test updating non-existent recipe
	err = store.Update(models.Recipe{ID: "non-existent"})
	if err == nil {
		t.Error("Expected error when updating non-existent recipe")
	}

	// Test deleting non-existent recipe
	err = store.Delete("non-existent")
	if err == nil {
		t.Error("Expected error when deleting non-existent recipe")
	}
}
