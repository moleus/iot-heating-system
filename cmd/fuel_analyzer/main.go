package main

import (
	"fmt"
	"iot-heating-system/cmd/fuel_analyzer/api"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server - структура для реализации интерфейса ServerInterface.
type Server struct{}

// GetAnalyze - реализация метода для анализа расхода топлива.
func (s *Server) GetAnalyze(c *gin.Context, params api.GetAnalyzeParams) {
	requiredTemp := params.RequiredTemp
	outsideTemp := params.OutsideTemp
	specificHeat := params.SpecificHeatOfCombustionFuel
	designOutsideTemp := params.DesignOutsideTemp
	designCapacity := params.DesignCapacity
	efficiency := params.Efficiency

	analyzer := NewFuelConsumptionAnalyzer(requiredTemp, outsideTemp, designOutsideTemp, designCapacity, specificHeat, efficiency)
	analyzer.Analyze()

	response := api.GetAnalyze200JSONResponse{
		FuelConsumption: &analyzer.FuelConsumption,
	}

	c.JSON(http.StatusOK, response)
}

func main() {
	r := gin.Default()

	si := &Server{}

	api.RegisterHandlers(r, si)

	fmt.Println("Starting server at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
