package helper

import (
	"fmt"
	"io"
	"net/http"
)

func GetURL(u string) (body []byte, err error) {

	res, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("helper url [%s] could not be fetched: %s", u, err))
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf(fmt.Sprintf("helper url [%s] could not be fetched, http code: %d", u, res.StatusCode))
	}
	body, err = io.ReadAll(res.Body)
	//log.Println("helper get url:", u, " responce: ", len(body)/1000, "kB")

	return
}
