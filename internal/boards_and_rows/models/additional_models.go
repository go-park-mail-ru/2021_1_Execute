package models

import (
	"2021_1_Execute/internal/tasks"
	"2021_1_Execute/internal/users"
)

type FullBoardInfo struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Owner       users.User    `json:"-"`
	Rows        []FullRowInfo `json:"rows"`
}
type FullRowInfo struct {
	ID       int          `json:"id"`
	Name     string       `json:"name"`
	Position int          `json:"position"`
	Tasks    []tasks.Task `json:"tasks"`
}
