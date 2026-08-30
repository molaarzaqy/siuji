package scorer

// Contoh slice konversi untuk Listening (Indeks = jumlah benar, Value = skor konversi)
// Perhatikan: Peserta yang benar 0 pun biasanya tetap mendapat skor minimal (misal: 31)
var listeningConversion = []int{
    31, 31, 31, 31, 31, 31, 31, 31, 31, 31, // 0 - 9 benar
    32, 33, 34, 35, 36, 37, 38, 39, 40, 41, // 10 - 19 benar
    42, 43, 44, 45, 46, 47, 48, 49, 50, 51, // 20 - 29 benar
    52, 53, 54, 55, 56, 57, 58, 59, 60, 61, // 30 - 39 benar
    62, 63, 64, 65, 66, 67, 68, 68, 68, 68, 68, // 40 - 50 benar
}

var structureConversion = []int{
    31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 
    32, 33, 35, 37, 38, 40, 41, 42, 43, 44, 
    45, 46, 48, 50, 52, 54, 56, 58, 60, 63, 
    65, 67, 68, 68, 68, 68, 68, 68, 68, 68, 68, // 0 - 40 benar
}

var readingConversion = []int{
    31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 
    32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 
    42, 43, 44, 45, 46, 48, 49, 51, 52, 54, 
    55, 57, 58, 60, 61, 63, 65, 67, 67, 67, 
    67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, // 0 - 50 benar
}

func ConvertListeningScore(correct int) int {
    if correct < 0 { correct = 0 }
    if correct >= len(listeningConversion) { correct = len(listeningConversion) - 1 }
    return listeningConversion[correct]
}

func ConvertStructureScore(correct int) int {
    if correct < 0 { correct = 0 }
    if correct >= len(structureConversion) { correct = len(structureConversion) - 1 }
    return structureConversion[correct]
}

func ConvertReadingScore(correct int) int {
    if correct < 0 { correct = 0 }
    if correct >= len(readingConversion) { correct = len(readingConversion) - 1 }
    return readingConversion[correct]
}

func CalculateFinalTOEFL(listeningCorrect, structureCorrect, readingCorrect int) int {
    lConv := ConvertListeningScore(listeningCorrect)
    sConv := ConvertStructureScore(structureCorrect)
    rConv := ConvertReadingScore(readingCorrect)

    sum := lConv + sConv + rConv
    
    // Rumus resmi TOEFL ITP: (Jumlah Ketiga Skor Konversi / 3) * 10
    finalScore := (float64(sum) / 3.0) * 10
    return int(finalScore)
}