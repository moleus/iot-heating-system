// Package main provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
	strictgin "github.com/oapi-codegen/runtime/strictmiddleware/gin"
)

// GetAnalyzeParams defines parameters for GetAnalyze.
type GetAnalyzeParams struct {
	RequiredTemp                 float32 `form:"required_temp" json:"required_temp"`
	OutsideTemp                  float32 `form:"outside_temp" json:"outside_temp"`
	Efficiency                   float32 `form:"efficiency" json:"efficiency"`
	SpecificHeatOfCombustionFuel float32 `form:"specific_heat_of_combustion_fuel" json:"specific_heat_of_combustion_fuel"`
	DesignOutsideTemp            float32 `form:"design_outside_temp" json:"design_outside_temp"`
	DesignCapacity               float32 `form:"design_capacity" json:"design_capacity"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Analyze fuel consumption
	// (GET /analyze)
	GetAnalyze(c *gin.Context, params GetAnalyzeParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetAnalyze operation middleware
func (siw *ServerInterfaceWrapper) GetAnalyze(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAnalyzeParams

	// ------------- Required query parameter "required_temp" -------------

	if paramValue := c.Query("required_temp"); paramValue != "" {

	} else {
		siw.ErrorHandler(c, fmt.Errorf("Query argument required_temp is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "required_temp", c.Request.URL.Query(), &params.RequiredTemp)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter required_temp: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Required query parameter "outside_temp" -------------

	if paramValue := c.Query("outside_temp"); paramValue != "" {

	} else {
		siw.ErrorHandler(c, fmt.Errorf("Query argument outside_temp is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "outside_temp", c.Request.URL.Query(), &params.OutsideTemp)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter outside_temp: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Required query parameter "efficiency" -------------

	if paramValue := c.Query("efficiency"); paramValue != "" {

	} else {
		siw.ErrorHandler(c, fmt.Errorf("Query argument efficiency is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "efficiency", c.Request.URL.Query(), &params.Efficiency)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter efficiency: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Required query parameter "specific_heat_of_combustion_fuel" -------------

	if paramValue := c.Query("specific_heat_of_combustion_fuel"); paramValue != "" {

	} else {
		siw.ErrorHandler(c, fmt.Errorf("Query argument specific_heat_of_combustion_fuel is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "specific_heat_of_combustion_fuel", c.Request.URL.Query(), &params.SpecificHeatOfCombustionFuel)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter specific_heat_of_combustion_fuel: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Required query parameter "design_outside_temp" -------------

	if paramValue := c.Query("design_outside_temp"); paramValue != "" {

	} else {
		siw.ErrorHandler(c, fmt.Errorf("Query argument design_outside_temp is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "design_outside_temp", c.Request.URL.Query(), &params.DesignOutsideTemp)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter design_outside_temp: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Required query parameter "design_capacity" -------------

	if paramValue := c.Query("design_capacity"); paramValue != "" {

	} else {
		siw.ErrorHandler(c, fmt.Errorf("Query argument design_capacity is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "design_capacity", c.Request.URL.Query(), &params.DesignCapacity)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter design_capacity: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetAnalyze(c, params)
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

	router.GET(options.BaseURL+"/analyze", wrapper.GetAnalyze)
}

type GetAnalyzeRequestObject struct {
	Params GetAnalyzeParams
}

type GetAnalyzeResponseObject interface {
	VisitGetAnalyzeResponse(w http.ResponseWriter) error
}

type GetAnalyze200JSONResponse struct {
	FuelConsumption *float32 `json:"fuel_consumption,omitempty"`
}

func (response GetAnalyze200JSONResponse) VisitGetAnalyzeResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetAnalyze400Response struct {
}

func (response GetAnalyze400Response) VisitGetAnalyzeResponse(w http.ResponseWriter) error {
	w.WriteHeader(400)
	return nil
}

type GetAnalyze405Response struct {
}

func (response GetAnalyze405Response) VisitGetAnalyzeResponse(w http.ResponseWriter) error {
	w.WriteHeader(405)
	return nil
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Analyze fuel consumption
	// (GET /analyze)
	GetAnalyze(ctx context.Context, request GetAnalyzeRequestObject) (GetAnalyzeResponseObject, error)
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

// GetAnalyze operation middleware
func (sh *strictHandler) GetAnalyze(ctx *gin.Context, params GetAnalyzeParams) {
	var request GetAnalyzeRequestObject

	request.Params = params

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetAnalyze(ctx, request.(GetAnalyzeRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetAnalyze")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(GetAnalyzeResponseObject); ok {
		if err := validResponse.VisitGetAnalyzeResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}