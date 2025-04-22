# Refactoring Plan: Layer-Based to Domain-Based Architecture

## Motivation
The current application is structured by technical layers (models, storage, handlers), which works well for small applications but has limitations as complexity grows. This refactoring will transition to a domain-based architecture that organizes code by feature/domain, to:

1. Improve cohesion of related functionality
2. Create clearer boundaries between application domains
3. Prepare for integration with the auth system
4. Align with go-kit service pattern
5. Support future microservice extraction

## Architectural Changes

### Current Architecture (Layer-Based)
```
go_recipe_app/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   │   └── recipe/
│   ├── models/
│   ├── storage/
│   │   ├── boltdb/
│   │   └── memory/
│   └── logging/
├── auth/         (new, not yet integrated)
│   ├── service.go
│   └── db/
└── templates/
```

### Target Architecture (Domain-Based)
```
go_recipe_app/
├── cmd/
│   └── main.go           (composition of services)
├── recipes/              (recipes domain)
│   ├── service.go        (business logic)
│   ├── handler.go        (HTTP handlers)
│   ├── models.go         (domain models)
│   └── db/
│       ├── repository.go (interface)
│       └── boltdb.go     (implementation)
├── auth/                 (auth domain, already structured)
│   ├── service.go
│   ├── handler.go
│   ├── middleware.go
│   ├── models.go
│   └── db/
│       ├── repository.go
│       └── boltdb.go
├── pkg/                  (shared utilities)
│   ├── config/
│   └── logging/
└── templates/
```

## Refactoring Phases

### Phase 1: Recipes Domain Creation

0. **Other stuff
   - [ ] Add user_recipes bucket to map users to recipes
   - [ ] Add public_recipes bucket to map public flagged recipes to ids

1. **Create Basic Structure**
   - [ ] Create `/recipes` directory
   - [ ] Create basic files: `service.go`, `handler.go`, `models.go`
   - [ ] Create `/recipes/db` subdirectory with `repository.go` and `boltdb.go`

2. **Move Recipe Models**
   - [ ] Copy Recipe model to `/recipes/models.go`
   - [ ] Add necessary imports
   - [ ] Update struct tags if needed
   - [ ] Add any domain-specific methods

3. **Create Recipe Repository**
   - [ ] Define repository interface in `/recipes/db/repository.go`
   - [ ] Move BoltDB implementation to `/recipes/db/boltdb.go`
   - [ ] Create constructor function with proper dependency injection

4. **Create Recipe Service**
   - [ ] Define service interface in `/recipes/service.go`
   - [ ] Implement business logic in service implementation
   - [ ] Move validation logic from handlers to service
   - [ ] Add logger dependency

5. **Move HTTP Handlers**
   - [ ] Create handler struct in `/recipes/handler.go`
   - [ ] Move handler methods from existing handlers
   - [ ] Add route registration function
   - [ ] Ensure handlers call service methods instead of directly accessing storage

### Phase 2: Update Main Composition

1. **Update Main Entry Point**
   - [ ] Refactor `main.go` to use new domain structure
   - [ ] Initialize recipe service and dependencies
   - [ ] Create and register recipe handler
   - [ ] Use shared router for registration

2. **Move Config & Logging**
   - [ ] Move config to `/pkg/config`
   - [ ] Move logging to `/pkg/logging`
   - [ ] Update imports across project

3. **Clean Up Legacy Code**
   - [ ] Once everything is working, mark old code as deprecated
   - [ ] Eventually remove old internal directory structure

### Phase 3: Prepare for Auth Integration

1. **Recipe Service Auth Integration**
   - [ ] Add auth service dependency to recipe service
   - [ ] Prepare service methods for permission checks
   - [ ] Add placeholders for ownership verification

2. **Recipe Model Extension**
   - [ ] Add `OwnerID` and `IsPublic` fields to Recipe model (without actual use yet)
   - [ ] Update repository to handle these fields but not enforce permissions

3. **Handler Auth Context**
   - [ ] Modify handlers to extract user info from request context
   - [ ] Prepare handlers for auth integration but don't enforce yet

### Phase 4: Final Cleanup & Testing

1. **Update Tests**
   - [ ] Create test package for recipe domain
   - [ ] Move and adapt existing tests
   - [ ] Add missing tests for new functionality

2. **Documentation**
   - [ ] Update codebase documentation
   - [ ] Add README to each domain directory
   - [ ] Document service interfaces

3. **Validation**
   - [ ] Run full test suite
   - [ ] Perform manual testing
   - [ ] Verify all functionality works as before

## Implementation Details

### Recipe Service Interface

```go
// service.go
package recipes

import (
	"context"
)

type Service interface {
	// Core operations
	ListRecipes(ctx context.Context) ([]Recipe, error)
	GetRecipe(ctx context.Context, id string) (Recipe, error)
	CreateRecipe(ctx context.Context, recipe Recipe) error
	UpdateRecipe(ctx context.Context, recipe Recipe) error
	DeleteRecipe(ctx context.Context, id string) error
	
	// Future methods for ownership/visibility
	// ListPublicRecipes(ctx context.Context) ([]Recipe, error)
	// ListUserRecipes(ctx context.Context, userID string) ([]Recipe, error)
}

type service struct {
	repo   db.Repository
	logger Logger
}

func NewService(repo db.Repository, logger Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}
```

### Recipe Repository Interface

```go
// db/repository.go
package db

import (
	"go_recipe_app/recipes"
)

type Repository interface {
	List() ([]recipes.Recipe, error)
	Get(id string) (recipes.Recipe, error)
	Create(recipe recipes.Recipe) error
	Update(recipe recipes.Recipe) error
	Delete(id string) error
	Close() error
}
```

### Main Composition

```go
// cmd/main.go
func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logger
	logger, err := logging.NewLogger(logging.LogConfig{
		Level:   cfg.LogLevel,
		Format:  cfg.LogFormat,
		LogPath: cfg.LogPath,
	})
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// Initialize templates
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		logger.Error("failed to parse templates", "error", err)
		return
	}

	// Initialize recipe domain
	recipeRepo, err := recipes.NewBoltRepository(cfg.DBPath, logger)
	if err != nil {
		logger.Error("failed to create recipe repository", "error", err)
		return
	}
	defer recipeRepo.Close()

	recipeService := recipes.NewService(recipeRepo, logger)
	recipeHandler := recipes.NewHandler(recipeService, tmpl, logger)

	// Create router
	router := mux.NewRouter()
	
	// Register recipe routes
	recipeHandler.RegisterRoutes(router)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	logger.Info("starting server", "addr", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("server failed", "error", err)
	}
}
```

## Migration Strategy

To ensure a smooth transition with minimal disruption:

1. **Parallel Implementation**: Keep existing code working while building new structure
2. **Incremental Changes**: Refactor one domain at a time, starting with recipes
3. **Feature Toggles**: Use config to switch between old and new implementations for testing
4. **Comprehensive Testing**: Maintain test coverage during transition
5. **Graceful Cleanup**: Remove old code only after new code is proven stable

## Benefits of this Approach

1. **Improved Code Organization**: Related functionality grouped together
2. **Better Testability**: Clearer boundaries make testing simpler
3. **Easier Maintenance**: Domains can evolve independently
4. **Reduced Coupling**: Dependencies clearly defined and injected
5. **Future-Proof Design**: Prepared for auth integration and microservices

## Next Steps After Refactoring

Once this refactoring is complete, we'll be well-positioned to implement the auth integration plan, with clearer interfaces between the recipe and auth domains.
