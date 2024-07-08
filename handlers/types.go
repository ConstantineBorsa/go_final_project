package handlers

type Task struct {
	ID string `db:"id" json:"id"`

	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}
