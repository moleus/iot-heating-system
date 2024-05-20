package main

// FuelConsumptionAnalyzer - агент анализатор расходов топлива
type FuelConsumptionAnalyzer struct {
	RequiredTemp                 float32
	OutsideTemp                  float32
	DesignOutsideTemp            float32
	DesignCapacity               float32
	SpecificHeatOfCombustionFuel float32
	efficiency                   float32
	FuelConsumption              float32
}

func NewFuelConsumptionAnalyzer(reqTemp, outTemp, designOutTemp float32, designCapacity float32, specHeatOfFuel float32, efficiency float32) *FuelConsumptionAnalyzer {
	return &FuelConsumptionAnalyzer{
		RequiredTemp:                 reqTemp,
		OutsideTemp:                  outTemp,
		DesignOutsideTemp:            designOutTemp,
		DesignCapacity:               designCapacity,
		SpecificHeatOfCombustionFuel: specHeatOfFuel,
		efficiency:                   efficiency,
	}
}

// CalculateHeatConsumption - расчет расхода тепла в ед. времени (МДж/ч = МВт)
func (a *FuelConsumptionAnalyzer) CalculateHeatConsumption() float32 {
	// Формула: Qов = Q`ов * (tвр - tср.от)/(tвр - tнр) * 60 * 60
	// где Q`ов - расход тепла на отопление, tвр - температура внутри помещения, tср.от - средняя температура наружного воздуха за час,
	// tнр - расчетная температура наружного воздуха для проектирования систем отопления
	if a.OutsideTemp > 15 {
		return 0
	}
	temperatureDifference := a.RequiredTemp - a.OutsideTemp
	heatConsumption := a.DesignCapacity * temperatureDifference / (a.RequiredTemp - a.DesignOutsideTemp) * 60 * 60
	return heatConsumption
}

// CalculateFuelConsumption - расчет расхода топлива в ед. времени (тыс.м3/час)
func (a *FuelConsumptionAnalyzer) CalculateFuelConsumption(heatConsumption float32) float32 {
	// Пример расчета расхода топлива
	// Формула: B = Qов / (qсг * кпд)
	// где F - расход топлива, Qов - расход тепла за 1 час, q - удельная теплота сгорания топлива (кДж/м^3), кпд - КПД котла

	fuelConsumption := heatConsumption / (a.SpecificHeatOfCombustionFuel * a.efficiency)

	return fuelConsumption
}

// Analyze - основной метод анализа и расчетов
func (a *FuelConsumptionAnalyzer) Analyze() {
	heatConsumption := a.CalculateHeatConsumption()
	a.FuelConsumption = a.CalculateFuelConsumption(heatConsumption)
}
