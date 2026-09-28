package repository

import (
	"context"

	"github.com/Snehaaa-Kri/job-apply-automator/internal/domain"
	"github.com/google/uuid"
)

// UserRepository defines what ANY user data store must be able to do.
// Node/TypeScript analogy:
//   interface IUserRepository {
//     findByID(id: string): Promise<User | null>
//     findByEmail(email: string): Promise<User | null>
//     create(user: User): Promise<User>
//     update(user: User): Promise<User>
//   }
//
// The service layer only knows about this interface.
// The Postgres implementation (user_repo.go) satisfies it.
// A mock implementation can satisfy it too — for testing.
//
// context.Context = Node's AbortSignal / request cancellation.
// Every DB/network call should accept a context so it can be cancelled
// if the HTTP request is cancelled (user closed browser, timeout, etc.)
type UserRepository interface {
	// FindByID fetches a user by their UUID primary key.
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// FindByEmail fetches a user by email. Used during login.
	// Returns (nil, nil) if not found — not an error, just absence.
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	// Create inserts a new user. Returns the created user with ID populated.
	Create(ctx context.Context, user *domain.User) (*domain.User, error)

	// Update saves changes to an existing user.
	Update(ctx context.Context, user *domain.User) (*domain.User, error)

	// SoftDelete marks a user as deleted without removing the row.
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
