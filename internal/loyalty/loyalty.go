package loyalty

import (
	bhandler "github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/balance/handler"
	bservice "github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/balance/service"
	ohandler "github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/orders/handler"
	oservice "github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/orders/service"
	uhandler "github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/user/handler"
	uservice "github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/user/service"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/middleware"
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/db"
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/jwt"
	"github.com/go-chi/chi/v5"
)

// Service has all services of the loyalty program.
type Service struct {
	orders  *oservice.OrdersService
	balance *bservice.BalanceService
	user    *uservice.UserService
}

// NewService creates a new loyalty service.
func NewService(db db.DB, jwtGen *jwt.Generator) *Service {
	var (
		orderService   = oservice.New(db)
		balanceService = bservice.New(db)
		userService    = uservice.New(db, jwtGen)
	)
	return &Service{
		orders:  &orderService,
		balance: &balanceService,
		user:    &userService,
	}
}

// Handler has all handlers of the loyalty program.
type Handler struct {
	orders  *ohandler.Handler
	balance *bhandler.Handler
	user    *uhandler.Handler
	service *Service
}

// NewHandler creates a new loyalty handler.
func NewHandler(service *Service, withJwt middleware.Middleware) *Handler {
	var (
		ordersHandler  = ohandler.New(service.orders, withJwt)
		balanceHandler = bhandler.New(service.balance, withJwt)
		userHandler    = uhandler.New(service.user)
	)
	return &Handler{
		orders:  ordersHandler,
		balance: balanceHandler,
		user:    userHandler,
		service: service,
	}
}

// RegisterRoutes registers all loyalty program routes.
func (h *Handler) RegisterRoutes(router chi.Router) {
	h.user.RegisterRoutes(router)
	h.balance.RegisterRoutes(router)
	// h.orders.RegisterRoutes(router)
}
