package main

import (
	"fmt"
	api "iot-heating-system/cmd/fuel_analyzer/api"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server - структура для реализации интерфейса ServerInterface.
type Server struct{}

func (s *Server) validateParams(params api.GetAnalyzeParams) (bool, string) {
	if params.DesignCapacity <= 0 {
		return false, "designCapacity must be positive"
	}
	if params.Efficiency < 0.01 || params.Efficiency > 1 {
		return false, "efficiency must be between 0.01 and 1"
	}
	if params.DesignOutsideTemp > 0 || params.DesignOutsideTemp < -50 {
		return false, "designOutsideTemp must be between 0 and -50"
	}
	if params.SpecificHeatOfCombustionFuel <= 0 {
		return false, "specificHeat must be greater than 0"
	}
	return true, ""
}

// GetAnalyze - реализация метода для анализа расхода топлива.
func (s *Server) GetAnalyze(c *gin.Context, params api.GetAnalyzeParams) {
	if valid, errMsg := s.validateParams(params); !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}
	requiredTemp := params.RequiredTemp
	outsideTemp := params.OutsideTemp
	specificHeat := params.SpecificHeatOfCombustionFuel
	designOutsideTemp := params.DesignOutsideTemp
	designCapacity := params.DesignCapacity
	efficiency := params.Efficiency

	analyzer := NewFuelConsumptionAnalyzer(requiredTemp, outsideTemp, designOutsideTemp, designCapacity, specificHeat, efficiency)
	analyzer.Analyze()

	response := api.GetAnalyze200JSONResponse{
		FuelConsumption: analyzer.FuelConsumption,
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
