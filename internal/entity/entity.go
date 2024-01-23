package entity

type Person struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Age        int64  `json:"age"`
	Gender     string `json:"gender"`
	Country    string `json:"country"`
}
