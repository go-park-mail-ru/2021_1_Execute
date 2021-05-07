package models

import (
	"2021_1_Execute/internal/tasks"
	"time"
)

type PostTaskRequest struct {
	RowID    int    `json:"row_id"`
	Name     string `json:"name"`
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
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
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
		Time:   time.Now().UTC().Format(time.RFC3339),
	}
}

type PostChecklistRequest struct {
	TaskID int      `json:"taskId"`
	Name   string   `json:"name"`
	Fields []string `json:"fields,omitempty"`
}

type PostChecklistResponse struct {
	ID int `json:"id"`
}

func PostChecklistToChecklist(input *PostChecklistRequest) tasks.Checklist {
	var fields []tasks.Field
	for _, field := range input.Fields {
		fields = append(fields, tasks.Field{
			Name: field,
			Done: false,
		})
	}
	return tasks.Checklist{
		Name:   input.Name,
		Fields: fields,
	}
}

type PatchChecklistRequest struct {
	Name   string        `json:"name,omitempty"`
	Fields []tasks.Field `json:"fields,omitempty"`
}

func PatchChecklistToChecklist(input *PatchChecklistRequest) tasks.Checklist {
	return tasks.Checklist{
		Name:   input.Name,
		Fields: input.Fields,
	}
}
