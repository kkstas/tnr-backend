package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kkstas/tr-backend/internal/config"
	"github.com/kkstas/tr-backend/internal/handlers/misc"
	"github.com/kkstas/tr-backend/internal/handlers/session"
	"github.com/kkstas/tr-backend/internal/handlers/user"
	mw "github.com/kkstas/tr-backend/internal/middleware"
	"github.com/kkstas/tr-backend/internal/services"
)

func SetupRoutes(
	cfg *config.Config,
	logger *slog.Logger,
	db *sql.DB,
	userService *services.UserService,
	vaultService *services.VaultService,
) http.Handler {
	mux := http.NewServeMux()

	requireAuth := mw.RequireAuth(cfg.JWTSecretKey, logger)
	withUser := mw.WithUser(logger, userService)

	mux.HandleFunc("GET /health-check", misc.HealthCheckHandler)
	mux.HandleFunc("/", misc.NotFoundHandler)

	mux.Handle("POST /login", session.LoginHandler(cfg.JWTSecretKey, logger, userService))
	mux.Handle("POST /register", mw.Enable(cfg.EnableRegister, session.RegisterHandler(logger, userService)))

	mux.Handle("GET /user", requireAuth(withUser(user.GetUserInfo(logger, userService))))

	return mux
}
