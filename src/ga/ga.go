package ga

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SendEvent(gaTrackingID, gaSecret, userID, command, parameter string) {
	// async send event
	go sendEvent(gaTrackingID, gaSecret, userID, command, parameter)
}

func sendEvent(gaTrackingID, gaSecret, userID, command, parameter string) {

	// Данные события
	event := map[string]interface{}{
		"name": "telegram-bot",
		"params": map[string]interface{}{
			"command":   command,
			"parameter": parameter,
		},
	}

	// Сформируйте запрос
	data := map[string]interface{}{
		"client_id": userID,
		"events":    []map[string]interface{}{event},
	}

	// Кодирование данных в JSON
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Println("Ошибка кодирования данных:", err)
		return
	}

	// Создание запроса
	url := fmt.Sprintf("https://www.google-analytics.com/mp/collect?api_secret=%s&measurement_id=%s", gaSecret, gaTrackingID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataJSON))
	if err != nil {
		log.Println("Ошибка создания запроса:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка выполнения запроса:", err)
		return
	}
	defer resp.Body.Close()

}
