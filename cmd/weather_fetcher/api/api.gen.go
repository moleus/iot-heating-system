// Package weather_fetcher provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package weather_fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
	strictgin "github.com/oapi-codegen/runtime/strictmiddleware/gin"
)

// DayForecast defines model for DayForecast.
type DayForecast struct {
	Hours []HourForecast `json:"hours"`
}

// DaysForecast defines model for DaysForecast.
type DaysForecast struct {
	Days []DayForecast `json:"days"`
}

// HourForecast defines model for HourForecast.
type HourForecast struct {
	Temperature float32   `json:"temperature"`
	Time        time.Time `json:"time"`
}

// GetWeatherParams defines parameters for GetWeather.
type GetWeatherParams struct {
	// Lat The latitude for which to fetch the weather forecast
	Lat float32 `form:"lat" json:"lat"`

	// Lon The longitude for which to fetch the weather forecast
	Lon float32 `form:"lon" json:"lon"`

	// Days The number of days to fetch the weather forecast for
	Days int `form:"days" json:"days"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get weather forecast
	// (GET /weather)
	GetWeather(c *gin.Context, params GetWeatherParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetWeather operation middleware
func (siw *ServerInterfaceWrapper) GetWeather(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetWeatherParams

	// ------------- Required query parameter "lat" -------------

	if paramValue := c.Query("lat"); paramValue != "" {

	} else {
		siw.ErrorHandler(c, fmt.Errorf("Query argument lat is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "lat", c.Request.URL.Query(), &params.Lat)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter lat: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Required query parameter "lon" -------------

	if paramValue := c.Query("lon"); paramValue != "" {

	} else {
		siw.ErrorHandler(c, fmt.Errorf("Query argument lon is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "lon", c.Request.URL.Query(), &params.Lon)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter lon: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Required query parameter "days" -------------

	if paramValue := c.Query("days"); paramValue != "" {

	} else {
		siw.ErrorHandler(c, fmt.Errorf("Query argument days is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "days", c.Request.URL.Query(), &params.Days)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter days: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetWeather(c, params)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/weather", wrapper.GetWeather)
}

type GetWeatherRequestObject struct {
	Params GetWeatherParams
}

type GetWeatherResponseObject interface {
	VisitGetWeatherResponse(w http.ResponseWriter) error
}

type GetWeather200JSONResponse DaysForecast

func (response GetWeather200JSONResponse) VisitGetWeatherResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetWeather400Response struct {
}

func (response GetWeather400Response) VisitGetWeatherResponse(w http.ResponseWriter) error {
	w.WriteHeader(400)
	return nil
}

type GetWeather500Response struct {
}

func (response GetWeather500Response) VisitGetWeatherResponse(w http.ResponseWriter) error {
	w.WriteHeader(500)
	return nil
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Get weather forecast
	// (GET /weather)
	GetWeather(ctx context.Context, request GetWeatherRequestObject) (GetWeatherResponseObject, error)
}

type StrictHandlerFunc = strictgin.StrictGinHandlerFunc
type StrictMiddlewareFunc = strictgin.StrictGinMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// GetWeather operation middleware
func (sh *strictHandler) GetWeather(ctx *gin.Context, params GetWeatherParams) {
	var request GetWeatherRequestObject

	request.Params = params

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetWeather(ctx, request.(GetWeatherRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetWeather")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(GetWeatherResponseObject); ok {
		if err := validResponse.VisitGetWeatherResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}
