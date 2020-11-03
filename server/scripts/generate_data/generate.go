//nolint
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/bxcodec/faker/v3"
)

type Register struct {
	Email     string `faker:"email"`
	Password  string `faker:"password"`
	FirstName string `faker:"first_name"`
	LastName  string `faker:"last_name"`
	Date      string `faker:"date"`
	Interests string `faker:"word"`
	Sex       string `faker:"oneof: M, F"`
	City      string `faker:"oneof: Moscow, New-York, Paris, London, Limassol"`
}

func main() {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	r := Register{}
	client := &http.Client{Transport: tr}

	wg := sync.WaitGroup{}

	for j := 0; j < 10; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 90_000; i++ {
				_ = faker.FakeData(&r)

				//fmt.Printf("%+v\n", v)
				reqBody, err := json.Marshal(map[string]string{
					"email":      r.Email,
					"password":   r.Password,
					"first_name": r.FirstName,
					"last_name":  r.LastName,
					"birth_date": fmt.Sprintf("%sT00:00:00Z", r.Date),
					"sex":        r.Sex,
					"interest":   r.Interests,
					"city":       r.City,
				})
				if err != nil {
					panic(err)
				}

				req, err := http.NewRequest("POST", "http://0.0.0.0:8091/api/register", bytes.NewBuffer(reqBody))
				if err != nil {
					panic(err)
				}
				req.Header.Set("Content-type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}

				err = printResp(resp)
				if err != nil {
					panic(err)
				}
			}
		}()
	}
	wg.Wait()
}

func printResp(resp *http.Response) error {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	return nil
}
