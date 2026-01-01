# Architecture Documentation

## Overview

This codebase follows **enterprise-grade design principles** with a clean, modular architecture using multiple design patterns.

## Design Patterns Implemented

### 1. **Repository Pattern**
- **Location**: `repositories/`
- **Purpose**: Abstracts data access layer from business logic
- **Benefits**: 
  - Easy to swap database implementations
  - Testable with mock repositories
  - Single Responsibility Principle

**Structure:**
```
repositories/
├── interfaces.go      # Repository interfaces
├── userRepository.go  # User repository implementation
└── errors.go         # Repository-specific errors
```

### 2. **Service Layer Pattern**
- **Location**: `services/`
- **Purpose**: Contains business logic separate from HTTP handling
- **Benefits**:
  - Reusable business logic
  - Testable without HTTP layer
  - Separation of concerns

**Structure:**
```
services/
├── interfaces.go    # Service interfaces
├── dto.go          # Data Transfer Objects
├── userService.go  # User business logic
├── authService.go  # Authentication business logic
└── errors.go       # Service-specific errors
```

### 3. **Factory Pattern**
- **Location**: `factories/`
- **Purpose**: Creates instances of repositories and services
- **Benefits**:
  - Centralized object creation
  - Easy to extend with new types
  - Consistent initialization

**Structure:**
```
factories/
├── repositoryFactory.go  # Creates repository instances
└── serviceFactory.go      # Creates service instances
```

### 4. **Dependency Injection Container**
- **Location**: `container/`
- **Purpose**: Manages all application dependencies
- **Benefits**:
  - Single source of truth for dependencies
  - Easy to test with mock containers
  - No global state

**Structure:**
```
container/
└── container.go  # DI container with all dependencies
```

## Architecture Layers

```
┌─────────────────────────────────────┐
│         HTTP Layer                  │
│  (Controllers, Routes, Middleware) │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         Service Layer                │
│  (Business Logic, Validation)       │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         Repository Layer             │
│  (Data Access, Database Operations) │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         Database                     │
│  (PostgreSQL via GORM)               │
└─────────────────────────────────────┘
```

## Dependency Flow

1. **main.go** → Creates `Container` using Factory Pattern
2. **Container** → Uses Factories to create Repositories and Services
3. **Routes** → Receives Container, extracts Services
4. **Controllers** → Receive Services (not database directly)
5. **Services** → Use Repositories for data access
6. **Repositories** → Use GORM for database operations

## Key Principles Applied

### SOLID Principles

1. **Single Responsibility Principle (SRP)**
   - Controllers: Handle HTTP requests/responses only
   - Services: Business logic only
   - Repositories: Data access only

2. **Open/Closed Principle (OCP)**
   - Interfaces allow extension without modification
   - New repositories/services can be added easily

3. **Liskov Substitution Principle (LSP)**
   - Any implementation of an interface can be substituted
   - Enables easy mocking for testing

4. **Interface Segregation Principle (ISP)**
   - Small, focused interfaces
   - Services don't depend on unused methods

5. **Dependency Inversion Principle (DIP)**
   - High-level modules depend on abstractions (interfaces)
   - Low-level modules implement interfaces

### Other Patterns

- **Dependency Injection**: All dependencies injected via constructor
- **Interface-Based Design**: Everything depends on interfaces
- **Factory Pattern**: Centralized object creation
- **Repository Pattern**: Data access abstraction
- **Service Layer**: Business logic separation

## Benefits of This Architecture

1. **Testability**: Easy to mock interfaces for unit testing
2. **Maintainability**: Clear separation of concerns
3. **Scalability**: Easy to add new features
4. **Flexibility**: Swap implementations without changing business logic
5. **No Global State**: All dependencies injected
6. **Type Safety**: Interfaces ensure contracts are met

## Example: Adding a New Feature

To add a new feature (e.g., "Products"):

1. **Create Model**: `models/product.go`
2. **Create Repository Interface**: `repositories/productRepository.go`
3. **Implement Repository**: `repositories/productRepository.go`
4. **Add to Factory**: `factories/repositoryFactory.go`
5. **Create Service Interface**: `services/productService.go`
6. **Implement Service**: `services/productService.go`
7. **Add to Factory**: `factories/serviceFactory.go`
8. **Add to Container**: `container/container.go`
9. **Create Controller**: `controllers/productController.go`
10. **Add Routes**: `routes/productRoutes.go`

Each layer is independent and testable!

## Testing Strategy

With this architecture, you can:

1. **Unit Test Services**: Mock repositories
2. **Unit Test Controllers**: Mock services
3. **Integration Test Repositories**: Use test database
4. **Integration Test Services**: Use real repositories with test DB
5. **E2E Tests**: Test full stack

## Migration from Old Architecture

The old architecture had:
- Controllers directly using `*gorm.DB`
- Business logic in controllers
- Global `initializers.DB` variable
- No separation of concerns

The new architecture:
- ✅ Controllers use Services
- ✅ Services contain business logic
- ✅ Repositories handle data access
- ✅ Dependency Injection via Container
- ✅ Factory Pattern for creation
- ✅ Interface-based design

## File Structure

```
goapi/
├── container/          # Dependency Injection Container
├── controllers/       # HTTP handlers (thin layer)
├── factories/          # Factory Pattern implementations
├── middleware/         # HTTP middleware
├── models/             # Data models
├── repositories/       # Data access layer
├── routes/             # Route definitions
├── services/           # Business logic layer
├── initializers/       # App initialization
└── main.go            # Application entry point
```

## Best Practices Followed

1. ✅ **No global state** (except initializers.DB for migrations)
2. ✅ **Interface-based design** for testability
3. ✅ **Dependency injection** throughout
4. ✅ **Factory pattern** for object creation
5. ✅ **Repository pattern** for data access
6. ✅ **Service layer** for business logic
7. ✅ **Error handling** with typed errors
8. ✅ **Separation of concerns** at every level

