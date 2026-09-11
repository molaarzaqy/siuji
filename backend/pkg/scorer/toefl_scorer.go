package scorer

import "math"

// CalculateFinalTOEFL computes the official ITP-style final score from the
// three already-converted section scaled scores (from ScoreConversionRepository).
func CalculateFinalTOEFL(listeningScaled, structureScaled, readingScaled int) int {
	sum := listeningScaled + structureScaled + readingScaled
	return int(math.Round((float64(sum) / 3.0) * 10))
}
