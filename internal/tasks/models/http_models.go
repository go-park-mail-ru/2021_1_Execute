package models

import "2021_1_Execute/internal/tasks"

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
