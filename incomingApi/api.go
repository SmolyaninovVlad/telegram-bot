package api

import (
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
	"telegaBot/structers"
	"github.com/gorilla/mux"
	"telegaBot/bot"
	// "strings"
)
// SetHeaders ...
func SetHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "mode, access-control-allow-origin, Authorization, Origin, Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func MainAPI(){
	router := mux.NewRouter()
	router.HandleFunc("/updateRequestToUser", updateRequestToUser).Methods("POST", "OPTIONS")
	router.Use(SetHeaders)
	fmt.Printf("API initialized on :2178\n")
	http.ListenAndServe(":2178", router)
}





func updateRequestToUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var (
		updateLetter 	Structers.UpdateRequestLetter
		decodedResponse Structers.Response
	)

	json.NewDecoder(request.Body).Decode(&updateLetter)
	chat_id, err := strconv.ParseInt(updateLetter.UserExternal_id, 10, 64)
	if err != nil {
		decodedResponse.Status = "failure"
		decodedResponse.Text = "Ошибка конвертации userExternal_id"
		fmt.Println("Ошибка конвертации", err)
		json.NewEncoder(response).Encode(decodedResponse)
		return
	}
	//отправка сообщения юзверю в телегу
	bot.SendMessageToUser(chat_id, getMessage(updateLetter))



	decodedResponse.Status = "success"
	decodedResponse.Text = "ok"
	json.NewEncoder(response).Encode(decodedResponse)
}

func getMessage(updateLetter Structers.UpdateRequestLetter) string {
	description := ""
	if len(updateLetter.Text)>0 {
		description = "\nСообщение: " + updateLetter.Text
	}
	message := "Обновление запроса №" + updateLetter.Request_id + "\n\nТема запроса: " + updateLetter.Subject +
	"\nСтатус запроса: " + getStatusRequest(strconv.Atoi(updateLetter.Status_id)) + description
	return message
}
func getStatusRequest( status_id int, err error) string {
	if err != nil {
		return "Не определено"
	}
	switch status_id {
		case 1: {
			return "Новая"
		}
		case 2: {
			return "В работе"
		}
		case 3: {
			return "Решена"
		}
		case 4: {
			return "Закрыта"
		}
	}
	return "Не определено"
}