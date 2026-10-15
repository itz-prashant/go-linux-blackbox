package analytics

import "math"

func (t *ColumnarTable) StdDevCpu() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	n := len(t.cpuValues)
	if n < 2 {
		return 0.0
	}

	// Mean
	var sum float64
	for _, v := range t.cpuValues {
		sum += v
	}

	mean := sum / float64(n)

	// 2. Variance (Differences ka square):
	var sumSquares float64
	for _, v := range t.cpuValues {
		diff := v - mean
		sumSquares += diff * diff
	}
	variance := sumSquares / float64(n)
	// 3. Standard Deviation (Square Root):
	return math.Sqrt(variance)
}

// ZScoreCPU current CPU ka Z-Score calculate karta hai
func (t *ColumnarTable) ZScoreCPU(currentCPU float64) float64 {
	maen := t.AvgCPU()
	stdDev := t.StdDevCpu()

	if stdDev < 0.001{
		return 0.0
	}

	return (currentCPU - maen) / stdDev
}

// DetectAnomaly check karta hai ki kya current CPU statistically abnormal hai (Z > threshold)
func (t *ColumnarTable) DetectAnomaly(currentCPU float64, threshold float64) (isAnomaly bool, zScore float64) {
	if t.Count() < 5 {
		return false , 0.0
	}
	z := t.ZScoreCPU(currentCPU)
	if z > threshold {
		return true, z
	}

	return false, z
}