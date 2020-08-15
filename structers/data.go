package Structers


type NewRequest struct {
	State 			int //Текущее состояние опроса пользователя
	Project_id 		string `json:"project_id"`//1
	Tracker_id 		string `json:"tracker_id"`//2
	Status_id 		string `json:"status_id"`//1
	Priority_id		string `json:"priority_id"`//2
	Subject			string `json:"subject"`//Тест
	Description		string `json:"description"`//Описание
	Custom_fields	[]*Customf `json:"custom_fields"`//структура
}
//bson{"id":"3", "name":"Подразделение-инициатор", "value":"Внешний канал"}
type Customf struct{
	Id 		string `json:"id"`//3
	Name	string `json:"name"`//Подразделение-инициатор
	Value	string `json:"value"`//Внешний канал
}

type Response struct {
	Status		string `json:"status"`
	Text		string `json:"text"`
}

type UpdateRequestLetter struct {
	Request_id		string	`json:"project_id"`
	Subject			string	`json:"subject"`
	Status_id		string	`json:"status_id"`
	UserExternal_id	string	`json:"userId"`
	Text			string  `json:"note"`
}
