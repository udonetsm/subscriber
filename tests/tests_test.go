package tests

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"subscriber/controllers"
	"subscriber/db"
	"testing"
	"time"
)

func Test(t *testing.T) {
	var reqs []*http.Request
	url1 := "http://localhost:8585?id=notexistsuri?"
	url2 := "http://localhost:8585/get?id=b563feb7b2b84b6test"
	url3 := "http://localhost:8585/get?id=notexistsorder"
	req1 := httptest.NewRequest(http.MethodGet, url1, nil) //unknown uri
	req3 := httptest.NewRequest(http.MethodGet, url2, nil) //valid request
	req4 := httptest.NewRequest(http.MethodGet, url3, nil) //unknown order_id
	reqs = append(reqs, req1, req3, req4)
	w := httptest.NewRecorder()
	resps := DoReqs(reqs, w)
	for _, res := range resps {
		log.Println(res)
	}
}

func DoReqs(reqs []*http.Request, w *httptest.ResponseRecorder) (answers []string) {
	for _, req := range reqs {
		controllers.GetById(w, req)
		body, err := ioutil.ReadAll(w.Body)
		if err != nil {
			log.Println(err)
			return
		}
		answers = append(answers, string(body))
	}
	return
}

func TestSetDB(t *testing.T) {
	var d [][]byte
	data2 := []byte("testing") //invalid data(not json)
	// valid data(order json)
	data3 := []byte(`{"order_uid":"b563feb7b2b84b6test","track_number":"WBILMTESTTRACK","entry":"WBIL","delivery":{"phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},"payment":{"currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,"goods_total":317},"items":[{"chrt_id":9934930,"track_number":"WBILMTESTTRACK","price":453,"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,"brand":"Vivienne Sabo","status":202}],"locale":"en","customer_id":"test","delivery_service":"meest","shardKey":"9","sm_id":99,"date_created":"2021-11-26T06:22:19Z","oof_shard":"1"}`)
	data4 := []byte(`"testkey":"testvalue"`) //invalid json(not order)
	d = append(d, data2, data3, data4)
	ids, errs := TrySet(d)
	for i, res := range ids {
		log.Println(res, errs[i])
	}
}

func TrySet(bytesdata [][]byte) (id []string, errs []error) {
	for _, m := range bytesdata {
		res, err := db.Set(m, time.Now().UnixNano())
		errs = append(errs, err)
		id = append(id, res)
	}
	return
}
