package api

import (
	"github.com/Snehaaa-Kri/job-apply-automator/config"
	"github.com/Snehaaa-Kri/job-apply-automator/internal/api/handlers"
	"github.com/Snehaaa-Kri/job-apply-automator/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// Router holds all handler dependencies.
// As we add more modules, we add more handlers here.
type Router struct {
	engine      *gin.Engine
	cfg         *config.Config
	authHandler *handlers.AuthHandler
}

// NewRouter wires up the Gin engine and all route groups.
// Node/Express analogy:
//   const app = express()
//   app.use('/api/v1/auth', authRoutes)
func NewRouter(
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
) *Router {
	// gin.New() = bare engine with no default middleware
	// gin.Default() = engine with Logger + Recovery middleware pre-attached
	// We use gin.New() and add our own middleware for full control
	engine := gin.New()

	// Recovery middleware: if a handler panics, recover gracefully (like Express's error handler)
	engine.Use(gin.Recovery())

	r := &Router{
		engine:      engine,
		cfg:         cfg,
		authHandler: authHandler,
	}

	r.registerRoutes()
	return r
}

// Engine exposes the underlying gin.Engine so main.go can start the server.
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

// registerRoutes sets up all route groups.
func (r *Router) registerRoutes() {
	// Health check — always useful to have
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// All API routes live under /api/v1
	// Node analogy: app.use('/api/v1', apiRouter)
	v1 := r.engine.Group("/api/v1")

	r.registerAuthRoutes(v1)
	// Future: r.registerJobRoutes(v1)
	// Future: r.registerResumeRoutes(v1)
}

// registerAuthRoutes sets up auth-specific routes.
// Node analogy: authRouter.post('/register', authController.register)
func (r *Router) registerAuthRoutes(v1 *gin.RouterGroup) {
	auth := v1.Group("/auth")
	{
		// Public routes — no JWT needed
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)

		// Protected routes — JWT middleware runs first
		// Node analogy: router.get('/me', protect, authController.getMe)
		auth.GET("/me", middleware.AuthRequired(r.cfg), r.authHandler.GetMe)
	}
}
