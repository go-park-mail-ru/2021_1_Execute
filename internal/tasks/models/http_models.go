package models

import (
	"2021_1_Execute/internal/tasks"
	"time"
)

type PostTaskRequest struct {
	RowID    int    `json:"row_id"`
	Name     string `json:"name" valid:"name"`
	Position int    `json:"position"`
}

type PostTaskResponse struct {
	ID int `json:"id"`
}

func TaskRequestToTask(req *PostTaskRequest) tasks.Task {
	return tasks.Task{
		Name:     req.Name,
		Position: req.Position,
	}
}

type GetTaskResponse struct {
	Task tasks.Task `json:"task"`
}

type PatchTaskRequest struct {
	Name        string `json:"name,omitempty" valid:"name"`
	Description string `json:"description,omitempty" valid:"description"`
}

func PatchTaskToTask(req *PatchTaskRequest) tasks.Task {
	return tasks.Task{
		Name:        req.Name,
		Description: req.Description,
		Position:    -1,
	}
}

type PostCommentRequest struct {
	TaskID int    `json:"taskId"`
	Text   string `json:"text"`
}

type PostCommentResponse struct {
	ID int `json:"id"`
}

func CommentRequestToComment(input *PostCommentRequest, author int) tasks.Comment {
	return tasks.Comment{
		Text:   input.Text,
		Author: author,
		Time:   time.Now(),
	}
}
