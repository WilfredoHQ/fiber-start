package models

type Error struct {
	Detail string `json:"detail"`
} // @Name Error

type ValidationError struct {
	Detail map[string]interface{} `json:"detail"`
} // @Name ValidationError
