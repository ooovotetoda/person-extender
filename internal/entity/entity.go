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

type Filters struct {
	Name       *string `json:"name,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	Patronymic *string `json:"patronymic,omitempty"`
	Age        *int64  `json:"age,omitempty"`
	Gender     *string `json:"gender,omitempty"`
	Country    *string `json:"country,omitempty"`
}
