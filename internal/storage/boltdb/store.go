package boltdb

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go_recipe_app/internal/models"

	bolt "go.etcd.io/bbolt"
)

var recipeBucket = []byte("recipes")

type Store struct {
	db     *bolt.DB
	logger *log.Logger
}

// New creates a new BoltDB store
func New(dbPath string) (*Store, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("could not open db: %v", err)
	}

	logger := log.New(os.Stdout, "[BOLTDB] ", log.LstdFlags)
	logger.Printf("Opening database at %s", dbPath)

	// Create recipes bucket if it doesn't exist
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(recipeBucket)
		if err != nil {
			return fmt.Errorf("could not create recipes bucket: %v", err)
		}
		logger.Println("Recipes bucket created/verified")
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Store{
		db:     db,
		logger: logger,
	}, nil
}

// Close closes the database
func (s *Store) Close() error {
	return s.db.Close()
}

// Create stores a new recipe
func (s *Store) Create(recipe models.Recipe) error {
	s.logger.Printf("Creating recipe: ID=%s, Title=%s", recipe.ID, recipe.Title)
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(recipeBucket)

		// Convert recipe to JSON
		buf, err := json.Marshal(recipe)
		if err != nil {
			return fmt.Errorf("could not marshal recipe: %v", err)
		}

		// Store using recipe ID as key
		err = b.Put([]byte(recipe.ID), buf)
		if err != nil {
			return fmt.Errorf("could not store recipe: %v", err)
		}

		s.logger.Printf("Successfully created recipe: %s", recipe.ID)
		return nil
	})
}

// Get reads a recipe from the DB
func (s *Store) Get(id string) (models.Recipe, error) {
	s.logger.Printf("Fetching recipe: %s", id)
	var recipe models.Recipe

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(recipeBucket)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("recipe not found: %s", id)
		}

		if err := json.Unmarshal(data, &recipe); err != nil {
			return fmt.Errorf("could not unmarshal recipe: %v", err)
		}
		return nil
	})

	if err != nil {
		s.logger.Printf("Error fetching recipe: %v", err)
		return models.Recipe{}, err
	}

	s.logger.Printf("Retrieved recipe: %s - %s", recipe.ID, recipe.Title)
	return recipe, nil
}

// Lists all recipes in the DB
func (s *Store) List() ([]models.Recipe, error) {
	s.logger.Println("Listing all recipes")
	var recipes []models.Recipe

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(recipeBucket)

		return b.ForEach(func(k, v []byte) error {
			var recipe models.Recipe
			if err := json.Unmarshal(v, &recipe); err != nil {
				return fmt.Errorf("could not unmarshal recipe: %v", err)
			}
			recipes = append(recipes, recipe)
			return nil
		})
	})

	if err != nil {
		s.logger.Printf("Error listing recipes: %v", err)
		return nil, err
	}

	s.logger.Printf("Found %d recipes", len(recipes))
	return recipes, nil
}

// Updates a recipe
func (s *Store) Update(recipe models.Recipe) error {
	s.logger.Printf("Updating recipe with ID: %s", recipe.ID)

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(recipeBucket)

		// Check if recipe exists
		if existing := b.Get([]byte(recipe.ID)); existing == nil {
			return fmt.Errorf("recipe not found: %s", recipe.ID)
		}

		buf, err := json.Marshal(recipe)
		if err != nil {
			return fmt.Errorf("could not marshal recipe: %v", err)
		}

		if err := b.Put([]byte(recipe.ID), buf); err != nil {
			return fmt.Errorf("could not update recipe: %v", err)
		}

		return nil
	})
}

// Deletes a recipe
func (s *Store) Delete(id string) error {
	s.logger.Printf("Deleting recipe with ID: %s", id)

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(recipeBucket)

		// Check if recipe exists
		if existing := b.Get([]byte(id)); existing == nil {
			return fmt.Errorf("recipe not found: %s", id)
		}

		if err := b.Delete([]byte(id)); err != nil {
			return fmt.Errorf("could not delete recipe: %v", err)
		}

		return nil
	})
}
