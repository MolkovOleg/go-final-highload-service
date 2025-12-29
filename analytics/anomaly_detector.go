package analytics

import (
	"math"
	"sync"
)

type AnomalyDetector struct {
	windowSize int
	values     []float64
	threshold  float64
	mu         sync.RWMutex
}

// Constructor
func NewAnomalyDetector(windowSize int, threshold float64) *AnomalyDetector {

	return &AnomalyDetector{
		windowSize: windowSize,
		values:     make([]float64, 0, windowSize),
		threshold:  threshold,
	}
}

func (ad *AnomalyDetector) Add(value float64) bool {
	ad.mu.Lock()
	defer ad.mu.Unlock()

	isAnomaly := false
	if len(ad.values) >= ad.windowSize {
		mean, stdDev := ad.CalculateStats()
		if stdDev > 0 {
			zScore := math.Abs((value - mean) / stdDev)
			isAnomaly = zScore > ad.threshold
		}
	}

	ad.values = append(ad.values, value)
	if len(ad.values) > ad.windowSize {
		ad.values = ad.values[1:]
	}

	return isAnomaly
}

func (ad *AnomalyDetector) CalculateStats() (mean, stdDev float64) {
	if len(ad.values) == 0 {
		return 0.0, 0.0
	}

	sum := 0.0
	for _, v := range ad.values {
		sum += v
	}

	mean = sum / float64(len(ad.values))

	varience := 0.0
	for _, v := range ad.values {
		varience += math.Pow(v - mean, 2)
	}
	varience /= float64(len(ad.values))
	stdDev = math.Sqrt(varience)

	return mean, stdDev
}

func (ad *AnomalyDetector) GetStatus() (mean, stdDev float64, count int) {
	ad.mu.RLock()
	defer ad.mu.RUnlock()
	mean, stdDev = ad.CalculateStats()
	
	return mean, stdDev, len(ad.values)
}

func (ad *AnomalyDetector) Reset() {
	ad.mu.Lock()
	defer ad.mu.Unlock()
	ad.values = make([]float64, 0, ad.windowSize)
}