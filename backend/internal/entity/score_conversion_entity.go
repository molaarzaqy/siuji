package entity

type ScoreConversion struct {
	ID           uint   `gorm:"column:id;primaryKey"`
	SectionType  string `gorm:"column:section_type"`
	CorrectCount int    `gorm:"column:correct_count"`
	ScaledScore  int    `gorm:"column:scaled_score"`
}

func (ScoreConversion) TableName() string {
	return "score_conversions"
}

