package def

type StateMachine struct {
	Name    string            `json:"Name"`
	StartAt string            `json:"StartAt"`
	States  map[string]string `json:"States"`
}

type Task struct {
	Type     string `json:"Type"`
	Resource string `json:"Resource"`
	Next     string `json:"Next"`
	End      bool   `json:"End"`
}

type Options struct {
	Variable     string `json:"Variable"`
	StringEquals string `json:"StringEquals"`
	Next         string `json:"Next"`
}

type Choice struct {
	Type    string    `json:"Type"`
	Choices []Options `json:"Choices"`
}
