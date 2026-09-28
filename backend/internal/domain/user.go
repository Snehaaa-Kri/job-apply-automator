package domain

// domains in go == models in node 

import (
	"time"

	"github.com/google/uuid"
)

// UserRole defines what kind of user this is.
// Node analogy: like an enum string union type in TypeScript — 'student' | 'professional'
type UserRole string

const (
	RoleStudent      UserRole = "student"
	RoleProfessional UserRole = "professional"
)

// AuthProvider tells us how the user signed up.
// 'local' = email+password, 'google' = Google OAuth
type AuthProvider string

const (
	AuthProviderLocal    AuthProvider = "local"
	AuthProviderGoogle   AuthProvider = "google"
	AuthProviderMicrosoft AuthProvider = "microsoft"
)

// User is the core user entity.
// Node/Mongoose analogy:
//   const UserSchema = new Schema({ email: { type: String, required: true, unique: true } })
//
// Struct tags explained:
//   `json:"email"`           → what key to use when sending JSON to frontend
//   `gorm:"uniqueIndex"`     → creates a UNIQUE constraint in the DB
//   `gorm:"not null"`        → NOT NULL constraint
//   `gorm:"-"`               → IGNORE this field (don't store in DB)
//   `json:"-"`               → NEVER send this field in JSON responses (e.g. password)
type User struct {
	ID        uuid.UUID  `json:"id"         gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time  `json:"createdAt"  gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt"  gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"-"          gorm:"index"` // soft delete — null means active

	// Identity
	Email    string `json:"email"    gorm:"uniqueIndex;not null"`
	Name     string `json:"name"     gorm:"not null"`
	Avatar   string `json:"avatar"`  // URL to profile picture

	// Auth
	// json:"-" means this field is NEVER included in API responses. Critical for security.
	PasswordHash string       `json:"-"          gorm:"column:password_hash"`
	Provider     AuthProvider `json:"provider"   gorm:"default:'local'"`
	ProviderID   string       `json:"-"          gorm:"column:provider_id"` // OAuth user ID

	// Profile
	Role          UserRole `json:"role"          gorm:"default:'student'"`
	IsVerified    bool     `json:"isVerified"    gorm:"default:false"`
	IsActive      bool     `json:"isActive"      gorm:"default:true"`

	// Career preferences (used for job matching later)
	PreferredRoles     []string `json:"preferredRoles"     gorm:"serializer:json"`
	PreferredLocations []string `json:"preferredLocations" gorm:"serializer:json"`
	MinSalary          int      `json:"minSalary"          gorm:"default:0"`
}

// TableName tells GORM to use "users" as the table name.
// Node analogy: mongoose.model('User', UserSchema) — sets the collection name.
func (User) TableName() string {
	return "users"
}
