package db

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"time"

	"github.com/go-kit/log"
	bolt "go.etcd.io/bbolt" // Fork of the original bolt project, backwards compatible with bolt and actively maintained
)

var (
	ErrUserNotFound = errors.New("User Not Found")
	ErrUserAlreadyExists = errors.New("User Already Exists")
)

// Defines instance of Bolt Auth Database with pointer to memory
type BoltAuthRepository struct {
	usersDB *bolt.DB
}

// Constructor for Bolt Auth Database instance
func NewBoltAuthRepository(dataDir string, mainLogger log.Logger) (*BoltAuthRepository, error) {
	mainLogger.Log("msg", "Opening on-disk BoltDB databases for Auth service")
	dbPath := filepath.Join(dataDir, "users.db")
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		mainLogger.Log("msg", "Error opening users database", "err", err)
		return nil, err
	}

	// db.Update(func(tx *bolt.Tx) error {...}) creates a wrote transaction, run the function, and handles commit/rollback
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		return err
	})
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("username_index"))
		return err
	})

	// check err assignments from above db transactions
	if err != nil {
		mainLogger.Log("msg", "Error creating users bucket", "err", err)
		db.Close()
		return nil, err
	}

	// Returns an initialized BoltAuthRepository, with usersDB field set to db variable
	// nil is returned since no errors occured
	return &BoltAuthRepository{
		usersDB: db, // db is a *bolt.DB pointer created on line 27
	}, nil
}

// --------
// BoltAuthRepository struct Methods
// --------

func (s *BoltAuthRepository) Close(logger log.Logger) error {
	err := s.usersDB.Close()
	if err != nil {
		logger.Log("msg", "Error closing users database", "err", err)
		return err
	}
	return nil

}

// CreateUser adds a new user to the store and updates the username index.
func (s *BoltAuthRepository) CreateUser(user *User) error {
	return s.usersDB.Update(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte("users")) // converts string "users" to a byte slice
		indexBucket := tx.Bucket([]byte("username_index")) // converts string "username_index" to a byte slice

		// Check if user already exists using the username index
		userID := indexBucket.Get([]byte(user.Username))
		if userID != nil {
			return ErrUserAlreadyExists
		}

		// Create user
		encodedUser, err := json.Marshal(user)
		if err != nil {
			return err
		}
		err = usersBucket.Put([]byte(user.UserID), encodedUser)
		if err != nil {
			return err
		}

		// Update username index
		return indexBucket.Put([]byte(user.Username), []byte(user.UserID))
	})
}

// GetUser retrieves a user by username.
func (s *BoltAuthRepository) GetUserByUsername(username string) (*User, error) {
	var user User
	err := s.usersDB.View(func(tx *bolt.Tx) error {
		indexBucket := tx.Bucket([]byte("username_index"))
		usersBucket := tx.Bucket([]byte("users"))

		// Get user ID from username index
		userID := indexBucket.Get([]byte(username))
		if userID == nil {
			return ErrUserNotFound
		}

		// Get user data by ID
		userData := usersBucket.Get(userID)
		if userData == nil {
			return ErrUserNotFound
		}
		return json.Unmarshal(userData, &user)
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID retrieves a user by ID.
func (s *BoltAuthRepository) GetUserByID(userID string) (*User, error) {
	var user User
	err := s.usersDB.View(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte("users"))

		// Get user data by ID
		userData := usersBucket.Get([]byte(userID))
		if userData == nil {
			return ErrUserNotFound
		}
		return json.Unmarshal(userData, &user)
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user.
func (s *BoltAuthRepository) UpdateUser(user *User) error {
	return s.usersDB.Update(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte("users"))
		indexBucket := tx.Bucket([]byte("username_index"))

		// Check if user exists
		userData := usersBucket.Get([]byte(user.UserID))
		if userData == nil {
			return ErrUserNotFound
		}

		// Update user data
		encodedUser, err := json.Marshal(user)
		if err != nil {
			return err
		}
		err = usersBucket.Put([]byte(user.UserID), encodedUser)
		if err != nil {
			return err
		}

		// Update username index if username changed
		var oldUser User
		err = json.Unmarshal(userData, &oldUser)
		if err != nil {
			return err
		}
		if oldUser.Username != user.Username {
			err = indexBucket.Delete([]byte(oldUser.Username))
			if err != nil {
				return err
			}
			err = indexBucket.Put([]byte(user.Username), []byte(user.UserID))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteUser removes a user from the store.
func (s *BoltAuthRepository) DeleteUser(id string) error {
	return s.usersDB.Update(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte("users"))
		indexBucket := tx.Bucket([]byte("username_index"))

		// Check if user exists
		userData := usersBucket.Get([]byte(id))
		if userData == nil {
			return ErrUserNotFound
		}

		// Delete user from users bucket
		err := usersBucket.Delete([]byte(id))
		if err != nil {
			return err
		}

		// Remove username from index
		var user User
		if err != nil {
			return err
		}
		return indexBucket.Delete([]byte(user.Username))
	})
}

// ListUsers returns a slice of all users.
func (s *BoltAuthRepository) ListUsers() ([]User, error) {
	var users []User
	err := s.usersDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		return b.ForEach(func(k, v []byte) error {
			var user User
			if err := json.Unmarshal(v, &user); err != nil {
				return err
			}
			users = append(users, user)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}