package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func BenchmarkRouter(b *testing.B) {

	url := "http://localhost:5656/pan/file/getlist"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		client := http.Client{}
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Set("web-token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODc3MDMwMDYsInVzZXJJZCI6ImNlc2hpIn0.ISMRI49Ha1q26SsHHCsCLZH93PYKorlSHTaurtoOi_4")

		response, e := client.Do(request)
		if e != nil {
			fmt.Println(e.Error())
		}
		body, _ := ioutil.ReadAll(response.Body)
		_ = string(body)
		//fmt.Println(bodystr)
	}

}
