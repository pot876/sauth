package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pot876/sauth/internal/chain"
	"github.com/pot876/sauth/internal/config"
	"github.com/pot876/sauth/internal/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

type Api struct {
	auth     chain.IServiseAuth
	validate chain.IServiseValidate

	metrics *ApiMetrics
}

func NewApi(ctx context.Context, cfg *config.Config) (*Api, error) {
	authService, err := NewAuth(ctx, cfg)
	if err != nil {
		log.Error().Err(err).Caller().Send()
		return nil, err
	}

	validateService, err := NewValidate(ctx, cfg)
	if err != nil {
		log.Error().Err(err).Caller().Send()
		return nil, err
	}

	api := &Api{
		auth:     authService,
		validate: validateService,
		metrics:  &ApiMetrics{},
	}

	return api, nil
}

func (a *Api) Login(c *gin.Context) {
	status := 0

	t0 := time.Now()
	defer func() {
		switch status / 100 {
		case 5:
			a.metrics.loginCounter5XX.Add(1.)
			a.metrics.loginDurationsXXX.Observe(float64(time.Since(t0).Seconds()))
		case 4:
			a.metrics.loginCounter4XX.Add(1.)
			a.metrics.loginDurationsXXX.Observe(float64(time.Since(t0).Seconds()))
		case 2:
			a.metrics.loginCounter2XX.Add(1.)
			a.metrics.loginDurations2XX.Observe(float64(time.Since(t0).Seconds()))
		default:
			a.metrics.loginDurationsXXX.Observe(float64(time.Since(t0).Seconds()))
		}
	}()

	a.login(c)
	status = c.Writer.Status()
}

func (a *Api) login(c *gin.Context) {
	ctx := c.Request.Context()

	req := struct {
		RealmID  string `json:"realm_id"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}
	err := c.BindJSON(&req)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid request format")
		return
	}

	realmID, err := uuid.Parse(req.RealmID)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid realm_id: value must be a valid UUID")
		return
	}
	if req.Login == "" {
		c.String(http.StatusBadRequest, "invalid login: value must not be empty")
		return
	}
	if req.Password == "" {
		c.String(http.StatusBadRequest, "invalid password: value must not be empty")
		return
	}

	resp, err := a.auth.Login(ctx, realmID, []byte(req.Login), []byte(req.Password))
	if err != nil {
		switch err {
		case chain.ErrNotFound, chain.ErrBadPassword:
			c.String(http.StatusUnauthorized, "invalid login or password")
			return
		}

		log.Error().Err(err).Caller().Send()
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(200, resp)
}

func (a *Api) Refresh(c *gin.Context) {
	status := 0

	t0 := time.Now()
	defer func() {
		switch status / 100 {
		case 5:
			a.metrics.refreshCounter5XX.Add(1.)
			a.metrics.refreshDurationsXXX.Observe(float64(time.Since(t0).Seconds()))
		case 4:
			a.metrics.refreshCounter4XX.Add(1.)
			a.metrics.refreshDurationsXXX.Observe(float64(time.Since(t0).Seconds()))
		case 2:
			a.metrics.refreshCounter2XX.Add(1.)
			a.metrics.refreshDurations2XX.Observe(float64(time.Since(t0).Seconds()))
		default:
			a.metrics.refreshDurationsXXX.Observe(float64(time.Since(t0).Seconds()))
		}
	}()

	a.refresh(c)
	status = c.Writer.Status()
}

func (a *Api) refresh(c *gin.Context) {
	ctx := c.Request.Context()

	req := struct {
		Token string `json:"token"`
	}{}
	err := c.BindJSON(&req)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid request format")
		return
	}

	resp, err := a.auth.Refresh(ctx, []byte(req.Token))
	if err != nil {
		switch err {
		case chain.ErrUnexpectedSigningMethod, chain.ErrInvalidToken, chain.ErrKeyNotFound:
			c.String(http.StatusUnauthorized, "invalid token")
			return
		}

		log.Error().Err(err).Caller().Send()
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(200, resp)
}

func (a *Api) Logout(c *gin.Context) {
	req := struct {
		Token string `json:"token"`
	}{}
	err := c.BindJSON(&req)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid request format")
		return
	}

	c.String(500, "not implemented")
}

func (a *Api) Validate(c *gin.Context) {
	ctx := c.Request.Context()

	req := struct {
		Token string `json:"token"`
	}{}
	err := c.BindJSON(&req)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid request format")
		return
	}

	err = a.validate.Validate(ctx, req.Token)
	if err != nil {
		log.Error().Err(err).Caller().Send()
		c.Status(http.StatusInternalServerError)
		return
	}
}

func (a *Api) Info(c *gin.Context) {
	ctx := c.Request.Context()

	req := struct {
		Token string `json:"token"`
	}{}
	err := c.BindJSON(&req)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid request format")
		return
	}

	_, err = a.validate.Info(ctx, req.Token)
	if err != nil {
		log.Error().Err(err).Caller().Send()
		c.Status(http.StatusInternalServerError)
		return
	}

}

func (a *Api) GetSessions(c *gin.Context) {
	c.String(500, "not implemented")
}

func (a *Api) GetSession(c *gin.Context) {
	c.String(500, "not implemented")
}

func (a *Api) DeleteSession(c *gin.Context) {
	c.String(500, "not implemented")
}

func (a *Api) RegisterEndpoints(cfg *config.Config, r *gin.Engine) {
	if cfg.HTTPEndpointLogin != "" {
		r.POST(cfg.HTTPEndpointPrefix+cfg.HTTPEndpointLogin, a.Login)
	}
	if cfg.HTTPEndpointRefresh != "" {
		r.POST(cfg.HTTPEndpointPrefix+cfg.HTTPEndpointRefresh, a.Refresh)
	}
}

func (a *Api) RegisterMetrics(r prometheus.Registerer, prefix string) {
	a.metrics.registerMetrics(r, prefix)
	log.Info().Msgf("metrics registered for http")

	if util.RegisterMetrics(a.auth, r, prefix) {
		log.Info().Msgf("metrics registered for auth component")
	}
	if util.RegisterMetrics(a.validate, r, prefix) {
		log.Info().Msgf("metrics registered for validate component")
	}
}
