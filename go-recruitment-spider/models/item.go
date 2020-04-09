package models

type ErrItem struct {
	Id  string `json:"id" bson:"_id"`
	Err string `json:"err" bson:"err"`
	LastErr string `json:"lastErr" bson:"lastErr"`
}


