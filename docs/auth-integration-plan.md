# Auth Integration Plan

## Architecture Overview

### Current Architecture
- Recipe-centric application with no user concept
- BoltDB storage for recipes in a single bucket
- Basic CRUD operations for recipes
- No authentication or authorization
- Single handler for all recipe operations

### Target Architecture
- User authentication via JWT tokens
- User-recipe ownership model
- Public/private recipe visibility
- Role-based permissions (admin, regular users)
- Saved recipes and collections functionality
- Clean separation between services

### Components to Add/Modify
1. **User Authentication Service**
   - JWT token generation and validation
   - User registration and login
   - Password hashing and verification

2. **User Data Storage**
   - BoltDB implementation for user data
   - Username indexing for lookups
   - Relations between users and recipes

3. **Authorization Middleware**
   - Token validation on protected routes
   - Role checking for admin functions
   - Ownership verification for recipe operations

4. **Extended Models**
   - User model with roles and profile data
   - Recipe model with owner ID and visibility
   - Collections model for grouping recipes

## Integration Tasks

### Phase 1: Core Authentication (MVP)

1. **Database Setup**
   - [ ] Configure auth database path in environment variables
   - [ ] Implement user bucket and username index in BoltDB
   - [ ] Create migration strategy for existing recipes (assign to admin)

2. **User Model & Storage**
   - [ ] Complete User model implementation
   - [ ] Finish BoltAuthRepository implementation
   - [ ] Add tests for user storage operations

3. **Authentication Service**
   - [ ] Implement JWT token generation
   - [ ] Add password hashing and verification
   - [ ] Create login and registration endpoints
   - [ ] Develop token validation mechanism

4. **Recipe Ownership**
   - [ ] Extend Recipe model with OwnerID field
   - [ ] Update RecipeStore interface to filter by owner
   - [ ] Modify recipe handlers to set owner on creation
   - [ ] Add ownership checks for update/delete operations

5. **Integration in Main App**
   - [ ] Initialize auth service in main.go
   - [ ] Register auth routes
   - [ ] Add authentication middleware
   - [ ] Update templates to include login/register forms

### Phase 2: Enhanced Features

1. **Recipe Visibility**
   - [ ] Add IsPublic field to Recipe model
   - [ ] Update List operation to filter by visibility
   - [ ] Add UI controls for setting recipe visibility
   - [ ] Modify templates to show/hide private recipes

2. **User Profiles**
   - [ ] Create user profile view
   - [ ] Add profile editing functionality
   - [ ] Implement password change feature
   - [ ] Add user avatar support (optional)

3. **Admin Functionality**
   - [ ] Add roles management
   - [ ] Create admin dashboard
   - [ ] Implement user management for admins
   - [ ] Add metrics and monitoring features

### Phase 3: Social Features

1. **Saved Recipes**
   - [ ] Create relationship model for saved recipes
   - [ ] Add save/unsave functionality
   - [ ] Implement saved recipes view
   - [ ] Update UI to show save status

2. **Recipe Collections**
   - [ ] Design collections model
   - [ ] Add CRUD operations for collections
   - [ ] Create collection membership functionality
   - [ ] Develop collection views and UI

## Implementation Notes

### Auth Middleware Implementation
```go
// Pseudocode for auth middleware
func AuthMiddleware(authService AuthService) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract token from Authorization header
            // Validate token using auth service
            // Add user details to request context
            // For public routes or valid tokens, proceed to next handler
            // For protected routes with invalid tokens, return 401
        })
    }
}
```

### Recipe-User Relationship
```go
// Extended Recipe model
type Recipe struct {
    // Existing fields
    ID          string
    Title       string
    // ...
    
    // New fields
    OwnerID     string  `json:"owner_id"`
    IsPublic    bool    `json:"is_public"`
}
```

### Main App Integration
```go
// Pseudocode for main.go changes
func main() {
    // Existing code...
    
    // Initialize auth components
    authRepo := auth.NewBoltAuthRepository(cfg.AuthDBPath, logger)
    authService := auth.NewAuthService(authRepo, cfg.JWTSecret, logger)
    
    // Create router with auth middleware
    router := mux.NewRouter()
    router.Use(auth.AuthMiddleware(authService))
    
    // Register routes
    authHandler := auth.NewHandler(authService, tmpl, logger)
    authHandler.RegisterRoutes(router.PathPrefix("/auth").Subrouter())
    
    // Update recipe handler to use auth
    recipeHandler := recipe.New(tmpl, store, logger, authService)
    recipeHandler.RegisterRoutes(router.PathPrefix("/recipes").Subrouter())
    
    // Serve application
    http.ListenAndServe(addr, router)
}
```

## Testing Strategy

1. **Unit Tests**
   - Test each auth service method in isolation
   - Mock dependencies for service tests
   - Test repository operations with test database

2. **Integration Tests**
   - Test auth middleware with various token scenarios
   - Verify recipe ownership enforcement
   - Test public/private recipe visibility

3. **End-to-End Tests**
   - Test user registration and login flow
   - Verify recipe creation with ownership
   - Test saved recipes and collections

## Future Considerations

1. **Performance Optimization**
   - Consider caching frequently accessed user data
   - Optimize BoltDB queries for user lookups
   - Add pagination for recipe and collection listings

2. **Security Enhancements**
   - Add rate limiting for login attempts
   - Implement CSRF protection
   - Add secure HTTP headers
   - Consider OAuth integration for social login

3. **Scalability**
   - Prepare for potential microservice separation
   - Design clean interfaces between components
   - Consider message passing for async operations
