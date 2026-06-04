package dto

type AdminDictionaryItemInput struct {
	ItemCode    string `json:"itemCode"`
	DisplayName string `json:"displayName"`
}

type AdminDictionarySaveSchemeRequest struct {
	SchemeID    *int64                     `json:"schemeId"`
	EditionCode string                     `json:"editionCode"`
	SchemeName  string                     `json:"schemeName"`
	Description string                     `json:"description"`
	Items       []AdminDictionaryItemInput `json:"items"`
}

type AdminDictionaryUpdateCurrentRequest struct {
	Items  []AdminDictionaryItemInput `json:"items"`
	Reason string                     `json:"reason"`
}

type AdminDictionaryApplySchemeRequest struct {
	SchemeID *int64 `json:"schemeId" binding:"required"`
}
