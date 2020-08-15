package post

import (
	"fmt"
	"io/ioutil"
	"telegaBot/structers"
	"net/http"
	"crypto/tls"
	"errors"
	"bytes"
	"encoding/json"
)

const postURL = "http://redmine-api/?action=Issue"

func SendPost(data *Structers.NewRequest) error{
	var decodedResponse Structers.Response
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
	
	params, err := json.Marshal(data)
	if err != nil{
		return err
	}
	buf := bytes.NewBufferString(string(params))

	fmt.Printf("buf: %s\n", buf)

	response, err := client.Post(postURL, "application/x-www-form-urlencoded", buf)
	if err != nil{
		return err
	}
	//отложенное закрытие
	defer response.Body.Close()
	//Получаем ответ
	ba, err := ioutil.ReadAll(response.Body)
	if err != nil{
		return errors.New("ошибка получения ответа от сервера")
	}
	//декодируем json
	err = json.Unmarshal(ba, &decodedResponse)
	if err != nil{
		return err
	}
	if decodedResponse.Status == "error" {
		return errors.New(decodedResponse.Text)
	}
	
	if response.StatusCode != 200 {
		return errors.New("not 200 response")
	} 

	//Отработало, ошибок нет
	return nil
}