package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kit/log"
	bolt "go.etcd.io/bbolt"
)

var (
	ErrRecipeNotFound = errors.New("Recipe Not Found")
)

// This is old model of DB filepath - use new model like Auth
// var recipeBucket = []byte("recipes")

// Defines instance of Bolt Recipe repository with pointer to memory
type BoltRecipeRepository struct {
	recipesDB *bolt.DB
	// logger *log.Logger -> Should logger be in struct or not? Claude says yes
}

// New creates a new BoltDB store
func NewBoltRecipeRepository(dataDir string, mainLogger log.Logger) (*BoltRecipeRepository, error) {
	mainLogger.Log("msg", "Opening on-disk BoltDB databases for Recipe service")
	dbPath := filepath.Join(dataDir, "recipes.db")
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		mainLogger.Log("msg", "Error opening recipe database", "err", err)
		return nil, err
	}

	// Create required buckets if they don't exist
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("recipes"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("user_recipes"))
		return err
	})

	// Check err assignments from above db transaction
	if err != nil {
		mainLogger.Log("msg", "Error creating buckets", "err", err)
		db.Close()
		return nil, err
	}

	// Returns an initialized BoltRecipeRepository
	return &BoltRecipeRepository{
		recipesDB: db,
	}, nil
}

// --------
// DATA HANDLING METHODS
// --------

// Close closes the database
func (s *BoltRecipeRepository) Close(logger log.Logger) error {
	err := s.recipesDB.Close()
	if err != nil {
		logger.Log("msg", "Error closing recipes database", "err", err)
		return err
	}
	return nil
}

// Helper function to update user_recipes index
func (s *BoltRecipeRepository) updateUserRecipeIndex(tx *bolt.Tx, recipeID string, ownerID string, isAdd bool) error {
	// Skip if ownerID is empty
	if ownerID == "" {
		return nil
	}

	indexBucket := tx.Bucket([]byte("user_recipes"))

	// Get current recipe list for user
	recipeList := []string{}
	userRecipes := indexBucket.Get([]byte(ownerID))

	if userRecipes != nil {
		recipeList = strings.Split(string(userRecipes), ",")
	}

	if isAdd {
		// Check if recipe already in list to avoid duplicates
		for _, id := range recipeList {
			if id == recipeID {
				return nil // Recipe already in user's list
			}
		}
		// Add recipe to list
		recipeList = append(recipeList, recipeID)
	} else {
		// Remove recipe from list
		newList := []string{}
		for _, id := range recipeList {
			if id != recipeID {
				newList = append(newList, id)
			}
		}
		recipeList = newList
	}

	// Save updated list
	return indexBucket.Put([]byte(ownerID), []byte(strings.Join(recipeList, ",")))
}

// CreateRecipe stores a new recipe
func (s *BoltRecipeRepository) CreateRecipe(recipe *Recipe, logger log.Logger) error {
	return s.recipesDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("recipes"))

		// Convert recipe to JSON
		buf, err := json.Marshal(recipe)
		if err != nil {
			logger.Log("msg", "Could not marshal recipe", "err", err)
			return err
		}

		// Store using recipe ID as key
		err = b.Put([]byte(recipe.ID), buf)
		if err != nil {
			logger.Log("msg", "Could not store recipe", "err", err)
			return err
		}

		// Update user_recipes index
		if err := s.updateUserRecipeIndex(tx, recipe.ID, recipe.OwnerID, true); err != nil {
			logger.Log("msg", "Could not update user_recipes index", "err", err)
			return err
		}

		logger.Log("msg", "Successfully created recipe", "id", recipe.ID)
		return nil
	})
}

// GetRecipe reads a recipe from the DB
func (s *BoltRecipeRepository) GetRecipe(id string, logger log.Logger) (*Recipe, error) {
	var recipe Recipe

	err := s.recipesDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("recipes"))

		// Get recipe by ID
		data := b.Get([]byte(id))
		if data == nil {
			return ErrRecipeNotFound
		}

		if err := json.Unmarshal(data, &recipe); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logger.Log("msg", "Error fetching recipe", "err", err, "id", id)
		return nil, err
	}

	logger.Log("msg", "Retrieved recipe", "id", recipe.ID, "title", recipe.Title)
	return &recipe, nil
}

// ListRecipes lists all recipes in the DB
func (s *BoltRecipeRepository) ListRecipes(logger log.Logger) ([]Recipe, error) {
	logger.Log("msg", "Listing all recipes")
	var recipes []Recipe

	err := s.recipesDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("recipes"))

		return b.ForEach(func(k, v []byte) error {
			var recipe Recipe
			if err := json.Unmarshal(v, &recipe); err != nil {
				return fmt.Errorf("could not unmarshal recipe: %v", err)
			}
			recipes = append(recipes, recipe)
			return nil
		})
	})

	if err != nil {
		logger.Log("msg", "Error listing recipes", "err", err)
		return nil, err
	}

	logger.Log("msg", "Found recipes", "count", len(recipes))
	return recipes, nil
}

// ListUserRecipes returns all recipes for a specific user
func (s *BoltRecipeRepository) ListUserRecipes(userID string, logger log.Logger) ([]Recipe, error) {
	logger.Log("msg", "Listing recipes for user", "userID", userID)
	var recipes []Recipe

	err := s.recipesDB.View(func(tx *bolt.Tx) error {
		indexBucket := tx.Bucket([]byte("user_recipes"))
		recipesBucket := tx.Bucket([]byte("recipes"))

		// Get list of recipe IDs for this user
		userRecipes := indexBucket.Get([]byte(userID))
		if userRecipes == nil || len(userRecipes) == 0 {
			return nil // No recipes found, not an error
		}

		// Split the comma-separated list
		recipeIDs := strings.Split(string(userRecipes), ",")

		// Fetch each recipe
		for _, id := range recipeIDs {
			data := recipesBucket.Get([]byte(id))
			if data != nil {
				var recipe Recipe
				if err := json.Unmarshal(data, &recipe); err != nil {
					return err
				}
				recipes = append(recipes, recipe)
			}
		}

		return nil
	})

	if err != nil {
		logger.Log("msg", "Error listing user recipes", "err", err, "userID", userID)
		return nil, err
	}

	logger.Log("msg", "Found user recipes", "count", len(recipes), "userID", userID)
	return recipes, nil
}

// UpdateRecipe updates a recipe
func (s *BoltRecipeRepository) UpdateRecipe(recipe Recipe, logger log.Logger) error {
	logger.Log("msg", "Updating recipe", "id", recipe.ID)

	return s.recipesDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("recipes"))

		// Check if recipe exists and get current data
		existingData := b.Get([]byte(recipe.ID))
		if existingData == nil {
			return ErrRecipeNotFound
		}

		// Unmarshal existing recipe to get the current owner
		var existingRecipe Recipe
		if err := json.Unmarshal(existingData, &existingRecipe); err != nil {
			return err
		}

		// Update recipe data
		buf, err := json.Marshal(recipe)
		if err != nil {
			return err
		}

		if err := b.Put([]byte(recipe.ID), buf); err != nil {
			return err
		}

		// Update user index if owner changed
		if existingRecipe.OwnerID != recipe.OwnerID {
			// Remove from old owner's list if there was an owner
			if existingRecipe.OwnerID != "" {
				if err := s.updateUserRecipeIndex(tx, recipe.ID, existingRecipe.OwnerID, false); err != nil {
					return err
				}
			}

			// Add to new owner's list if there is a new owner
			if recipe.OwnerID != "" {
				if err := s.updateUserRecipeIndex(tx, recipe.ID, recipe.OwnerID, true); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// DeleteRecipe deletes a recipe
func (s *BoltRecipeRepository) DeleteRecipe(id string, logger log.Logger) error {
	logger.Log("msg", "Deleting recipe", "id", id)

	return s.recipesDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("recipes"))

		// Check if recipe exists and get current data for index update
		existingData := b.Get([]byte(id))
		if existingData == nil {
			return ErrRecipeNotFound
		}

		// Unmarshal to get owner for index update
		var existingRecipe Recipe
		if err := json.Unmarshal(existingData, &existingRecipe); err != nil {
			return err
		}

		// Delete from recipe bucket
		if err := b.Delete([]byte(id)); err != nil {
			return err
		}

		// Update user_recipes index - remove from owner's list
		if existingRecipe.OwnerID != "" {
			if err := s.updateUserRecipeIndex(tx, id, existingRecipe.OwnerID, false); err != nil {
				return err
			}
		}

		return nil
	})
}
