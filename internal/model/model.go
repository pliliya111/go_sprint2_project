package model

type Task struct {
	ID        string      `json:"id"`
	Arg1      interface{} `json:"arg1"`
	Arg2      interface{} `json:"arg2"`
	Operation string      `json:"operation"`
	Result    interface{} `json:"result"`
}

type Expression struct {
	ID         string      `json:"id"`
	Expression string      `json:"expression"`
	Status     string      `json:"status"` // pending, in_progress, completed
	Result     interface{} `json:"result"`
}
