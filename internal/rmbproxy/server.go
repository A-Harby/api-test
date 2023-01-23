package rmbproxy

import (
	"bytes"
	"net/http"
	"strconv"

	// swagger configuration
	"github.com/pkg/errors"
	"github.com/threefoldtech/grid_proxy_server/internal/explorer/mw"
	"github.com/threefoldtech/substrate-client"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// sendMessage godoc
// @Summary submit the message
// @Description submit the message
// @Tags RMB
// @Accept  json
// @Produce  json
// @Param msg body Message true "rmb.Message"
// @Param twin_id path int true "twin id"
// @Success 200 {object} MessageIdentifier
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Failure 502 {object} string
// @Router /twin/{twin_id} [post]
func (a *App) sendMessage(r *http.Request) (*http.Response, mw.Response) {
	twinIDString := mux.Vars(r)["twin_id"]

	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(r.Body); err != nil {
		return nil, mw.BadRequest(err)
	}

	twinID, err := strconv.Atoi(twinIDString)
	if err != nil {
		return nil, mw.BadRequest(errors.Wrap(err, "invalid twin_id"))
	}

	c, err := a.resolver.Get(twinID)
	if err != nil {
		log.Error().Err(err).Msg("failed to create twin client")
		return nil, mw.Error(errors.Wrap(err, "failed to create twin client"))
	}

	response, err := c.SubmitMessage(*buffer)
	if err != nil {
		return nil, mw.BadGateway(errors.Wrap(err, "failed to submit message"))
	}
	return response, nil
}

// getResult godoc
// @Summary Get the message result
// @Description Get the message result
// @Tags RMB
// @Accept  json
// @Produce  json
// @Param twin_id path int true "twin id"
// @Param retqueue path string true "message retqueue"
// @Success 200 {array} Message
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Failure 502 {object} string
// @Router /twin/{twin_id}/{retqueue} [get]
func (a *App) getResult(r *http.Request) (*http.Response, mw.Response) {
	twinIDString := mux.Vars(r)["twin_id"]
	retqueue := mux.Vars(r)["retqueue"]

	reqBody := MessageIdentifier{
		Retqueue: retqueue,
	}

	twinID, err := strconv.Atoi(twinIDString)
	if err != nil {
		return nil, mw.BadRequest(errors.Wrap(err, "invalid twin_id"))
	}

	c, err := a.resolver.Get(twinID)
	if err != nil {
		log.Error().Err(err).Msg("failed to create twin client")
		return nil, mw.Error(errors.Wrap(err, "failed to create twin client"))
	}

	response, err := c.GetResult(reqBody)
	if err != nil {
		return nil, mw.BadGateway(errors.Wrap(err, "failed to submit message"))
	}
	return response, nil
}

// ping godoc
// @Summary ping the server
// @Description ping the server to check if it is running
// @Tags ping
// @Accept  json
// @Produce  json
// @Success 200 {object} PingMessage
// @Router /ping [get]
func (a *App) ping(r *http.Request) (interface{}, mw.Response) {
	return PingMessage{Ping: "pong"}, mw.Ok()
}

// Setup : sets rmb routes
// @title RMB proxy API
// @version 1.0
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email soberkoder@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func Setup(router *mux.Router, substrate *substrate.Substrate, rmbTimeout int) error {
	log.Info().Msg("Creating server")

	resolver, err := NewTwinResolver(substrate, rmbTimeout)
	if err != nil {
		return errors.Wrap(err, "couldn't get a client to explorer resolver")
	}

	a := &App{
		resolver: resolver,
	}

	router.HandleFunc("/twin/{twin_id:[0-9]+}", mw.AsProxyHandlerFunc(a.sendMessage))
	router.HandleFunc("/twin/{twin_id:[0-9]+}/{retqueue}", mw.AsProxyHandlerFunc(a.getResult))
	router.HandleFunc("/ping", mw.AsHandlerFunc(a.ping))
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	return nil
}
