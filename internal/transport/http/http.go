package http

import (
	_ "avito_intern/docs"
	"avito_intern/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/rs/cors"
	"go.uber.org/zap"
	// swaggerFiles "github.com/swaggo/files".
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
	decoder  = schema.NewDecoder()
)

type transport struct {
	s service.Service
	l *zap.SugaredLogger
}

func New(t service.Service, l *zap.SugaredLogger) http.Handler {
	r := &transport{t, l}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /user_banner", r.GetBannerForUser)
	mux.HandleFunc("GET /banner", r.GetBanners)
	mux.HandleFunc("GET /banner_history/{id}", r.GetBannerWithHistory)

	mux.HandleFunc("POST /banner", r.CreateBanner)
	mux.HandleFunc("DELETE /banner/{id}", r.DeleteBannerByID)
	mux.HandleFunc("DELETE /banner", r.DeleteBannerByTagAndFeature)
	mux.HandleFunc("PATCH /banner/{id}", r.ChangeBanner)

	mux.HandleFunc("POST /login", r.Login)

	mux.HandleFunc("GET /swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"), //The url pointing to API definition
	))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:8000/"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)
	return handler
}
