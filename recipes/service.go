package recipes

/*
ARCHITECTURAL PSEUDOCODE - RECIPE SERVICE WITH AUTH INTEGRATION

FUNCTION NewRecipeService(repository, logger)
    Return initialized RecipeService with repository and logger
END FUNCTION

FUNCTION GetRecipe(id, requestingUserID)
    Fetch recipe from repository

    IF recipe not found THEN
        Return not found error
    END IF

    // Apply visibility rules
    IF recipe is private AND requestingUserID is not owner THEN
        Return not found error (hide existence)
    END IF

    Return recipe
END FUNCTION

FUNCTION ListRecipes(requestingUserID)
    IF requestingUserID is empty (not logged in) THEN
        Return only public recipes
    ELSE
        Return public recipes AND user's own recipes
    END IF
END FUNCTION

FUNCTION ListUserRecipes(userID)
    Validate userID
    Fetch all recipes owned by this user
    Return recipes list
END FUNCTION

FUNCTION CreateRecipe(recipe, ownerID)
    Validate recipe data
    Set owner ID to current user's ID
    Save recipe to repository
    Return success/error
END FUNCTION

FUNCTION UpdateRecipe(recipe, requestingUserID)
    Fetch existing recipe

    // Check ownership
    IF requestingUserID is not recipe owner THEN
        Return unauthorized error
    END IF

    Update recipe in repository
    Return success/error
END FUNCTION

FUNCTION DeleteRecipe(id, requestingUserID)
    Fetch existing recipe

    // Check ownership
    IF requestingUserID is not recipe owner THEN
        Return unauthorized error
    END IF

    Delete recipe from repository
    Return success/error
END FUNCTION
*/

import (
	"errors"

	"github.com/go-kit/log"
	"github.com/google/uuid"

	"go_recipe_app/recipes/db"
)

var (
	ErrRecipeNotFound = errors.New("recipe not found")
	ErrUnauthorized   = errors.New("unauthorized: you don't own this recipe")
	ErrInvalidData    = errors.New("invalid recipe data")
)

// RecipeService defines the business logic for recipe management
type RecipeService interface {
	GetRecipe(id string, requestingUserID string) (*db.Recipe, error)
	ListPublicRecipes() ([]db.Recipe, error)
	ListUserRecipes(userID string) ([]db.Recipe, error)
	CreateRecipe(recipe *db.Recipe, ownerID string) error
	UpdateRecipe(recipe *db.Recipe, requestingUserID string) error
	DeleteRecipe(id string, requestingUserID string) error
}

type recipeServiceImpl struct {
	repo   db.RecipeRepository
	logger log.Logger
}

// NewRecipeService creates a new recipe service
func NewRecipeService(dataDir string, mainLogger log.Logger, recipeLogger log.Logger) (RecipeService, error) {
	repo, err := db.NewBoltRecipeRepository(dataDir, mainLogger)
	if err != nil {
		mainLogger.Log("msg", "Could not initialize recipes data repository", "err", err)
		return nil, err
	}
	
	return &recipeServiceImpl{
		repo:   repo,
		logger: recipeLogger,
	}, nil
}

// GetRecipe fetches a recipe by ID with visibility rules enforced
func (s *recipeServiceImpl) GetRecipe(id string, requestingUserID string) (*db.Recipe, error) {
	recipe, err := s.repo.GetRecipe(id, s.logger)
	if err != nil {
		s.logger.Log("msg", "Error fetching recipe", "id", id, "err", err)
		return nil, ErrRecipeNotFound
	}

	// If recipe is private, only the owner can see it
	if !recipe.IsPublic && recipe.OwnerID != requestingUserID {
		s.logger.Log("msg", "Unauthorized recipe access attempt", "id", id, "user", requestingUserID)
		return nil, ErrRecipeNotFound // Hide existence from non-owners
	}

	return recipe, nil
}

// ListPublicRecipes returns all public recipes
func (s *recipeServiceImpl) ListPublicRecipes() ([]db.Recipe, error) {
	allRecipes, err := s.repo.ListRecipes(s.logger)
	if err != nil {
		s.logger.Log("msg", "Error listing recipes", "err", err)
		return nil, err
	}

	// Filter out private recipes
	publicRecipes := []db.Recipe{}
	for _, recipe := range allRecipes {
		if recipe.IsPublic {
			publicRecipes = append(publicRecipes, recipe)
		}
	}

	return publicRecipes, nil
}

// ListUserRecipes returns all recipes owned by a specific user
func (s *recipeServiceImpl) ListUserRecipes(userID string) ([]db.Recipe, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	recipes, err := s.repo.ListUserRecipes(userID, s.logger)
	if err != nil {
		s.logger.Log("msg", "Error listing user recipes", "userID", userID, "err", err)
		return nil, err
	}

	return recipes, nil
}

// CreateRecipe creates a new recipe
func (s *recipeServiceImpl) CreateRecipe(recipe *db.Recipe, ownerID string) error {
	// Basic validation
	if recipe.Title == "" {
		return ErrInvalidData
	}

	// Set recipe ID and owner
	recipe.ID = uuid.New().String()
	recipe.OwnerID = ownerID

	err := s.repo.CreateRecipe(recipe, s.logger)
	if err != nil {
		s.logger.Log("msg", "Error creating recipe", "err", err)
		return err
	}

	return nil
}

// UpdateRecipe updates an existing recipe after checking ownership
func (s *recipeServiceImpl) UpdateRecipe(recipe *db.Recipe, requestingUserID string) error {
	// Check if recipe exists and verify ownership
	existingRecipe, err := s.repo.GetRecipe(recipe.ID, s.logger)
	if err != nil {
		s.logger.Log("msg", "Recipe not found for update", "id", recipe.ID, "err", err)
		return ErrRecipeNotFound
	}

	if existingRecipe.OwnerID != requestingUserID {
		s.logger.Log("msg", "Unauthorized update attempt", "id", recipe.ID, "user", requestingUserID)
		return ErrUnauthorized
	}

	// Preserve ownership (prevent ownership transfer)
	recipe.OwnerID = existingRecipe.OwnerID

	err = s.repo.UpdateRecipe(*recipe, s.logger)
	if err != nil {
		s.logger.Log("msg", "Error updating recipe", "id", recipe.ID, "err", err)
		return err
	}

	return nil
}

// DeleteRecipe deletes a recipe after checking ownership
func (s *recipeServiceImpl) DeleteRecipe(id string, requestingUserID string) error {
	// Check if recipe exists and verify ownership
	existingRecipe, err := s.repo.GetRecipe(id, s.logger)
	if err != nil {
		s.logger.Log("msg", "Recipe not found for deletion", "id", id, "err", err)
		return ErrRecipeNotFound
	}

	if existingRecipe.OwnerID != requestingUserID {
		s.logger.Log("msg", "Unauthorized delete attempt", "id", id, "user", requestingUserID)
		return ErrUnauthorized
	}

	err = s.repo.DeleteRecipe(id, s.logger)
	if err != nil {
		s.logger.Log("msg", "Error deleting recipe", "id", id, "err", err)
		return err
	}

	return nil
}
