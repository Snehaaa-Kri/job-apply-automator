package postgres

import (
	"context"
	"errors"

	"github.com/Snehaaa-Kri/job-apply-automator/internal/domain"
	"github.com/Snehaaa-Kri/job-apply-automator/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// userRepo is the PostgreSQL implementation of repository.UserRepository.
// It holds a reference to the DB connection.
//
// Node analogy: class UserRepository { constructor(private db: PgPool) {} }
// The lowercase 'u' means this type is NOT exported — external packages
// must use NewUserRepo() to get an instance. This enforces encapsulation.
type userRepo struct {
	db *gorm.DB
}

// NewUserRepo creates a new userRepo and returns it as the interface type.
// This is the CONSTRUCTOR pattern in Go.
//
// Node analogy: export function createUserRepository(db) { return new UserRepository(db) }
//
// Key trick: return type is `repository.UserRepository` (the interface),
// not `*userRepo` (the concrete struct). This means the caller only ever
// sees the interface — they can't accidentally call internal methods.
func NewUserRepo(db *gorm.DB) repository.UserRepository {
	return &userRepo{db: db}
}

// FindByID fetches a user by UUID.
// Node/Mongoose analogy: User.findById(id)
func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User

	// gorm.ErrRecordNotFound is like Mongoose returning null
	result := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error — caller handles it
		}
		return nil, result.Error
	}

	return &user, nil
}

// FindByEmail fetches a user by email address.
// Node/Mongoose analogy: User.findOne({ email })
func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	result := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

// Create inserts a new user into the database.
// Node/Mongoose analogy: new User({ ...data }).save()
func (r *userRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// Update saves all changed fields of a user.
// Node/Mongoose analogy: user.save() after modifying fields
func (r *userRepo) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// SoftDelete sets deleted_at to NOW() instead of removing the row.
// Node/Mongoose analogy: user.deleteOne() but with { deletedAt: new Date() } instead
func (r *userRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("deleted_at", "NOW()")
	return result.Error
}
