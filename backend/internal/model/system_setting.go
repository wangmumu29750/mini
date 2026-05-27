package model

type SystemSetting struct {
	BaseModel

	Key         string `gorm:"size:80;not null;uniqueIndex" json:"key"`
	Value       string `gorm:"size:255;not null" json:"value"`
	ValueType   string `gorm:"size:20;not null;default:STRING" json:"valueType"`
	Description string `gorm:"size:255" json:"description"`
}
