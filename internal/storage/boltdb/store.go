package boltdb

import (
	"encoding/json"
	"fmt"
	"time"

	"go_recipe_app/internal/models"

	bolt "go.etcd.io/bbolt"
)

var recipeBucket = []byte("recipes")

type Store struct {
	db *bolt.DB
}

// New creates a new BoltDB store
func New(dbPath string) (*Store, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("could not open db: %v", err)
	}

	// Create recipes bucket if it doesn't exist
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(recipeBucket)
		if err != nil {
			return fmt.Errorf("could not create recipes bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

// Close closes the database
func (s *Store) Close() error {
	return s.db.Close()
}

// Create stores a new recipe
func (s *Store) Create(recipe models.Recipe) error {
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

		return nil
	})
}
