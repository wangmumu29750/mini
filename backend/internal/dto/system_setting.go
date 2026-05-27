package dto

type SystemSettingResponse struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	ValueType   string `json:"valueType"`
	Description string `json:"description"`
}

type UpdateSystemSettingsRequest struct {
	Settings []UpdateSystemSettingItem `json:"settings" binding:"required"`
}

type UpdateSystemSettingItem struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}
