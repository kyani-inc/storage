package elasticsearch_test

import (
	"net/http"
	"os"
	"testing"

	"encoding/json"
	"github.com/kyani-inc/storage/providers/elasticsearch"
	"github.com/subosito/gotenv"
	"reflect"
)

var (
	host, index, scheme, namespace string
)

type person struct {
	Name      string   `json:"name"`
	Age       int      `json:"age"`
	Nicknames nickname `json:"nicknames"`
}

type nickname struct {
	Name string `json:"nickname"`
}

func init() {
	gotenv.Load(".env")
	host = os.Getenv("ES_HOST")
	index = os.Getenv("ES_INDEX")
	namespace = os.Getenv("ES_NAMESPACE")
	scheme = os.Getenv("ES_SCHEME")
}

func emptyVars() bool {
	return host == "" || index == "" || scheme == ""

}

func TestES(t *testing.T) {
	if emptyVars() {
		t.Error("Must have ENV set")
	}

	es, err := elasticsearch.New(host, scheme, index, namespace, http.Client{})

	if err != nil {
		t.Errorf("Building new ES client error %+v", err)
	}

	td := person{
		Name:      "Crit",
		Age:       100,
		Nicknames: nickname{Name: "Waffle"},
	}

	data, err := json.Marshal(td)

	err = es.Put("Crit", data)

	if err != nil {
		t.Errorf("Error on Put error:%s", err)
	}

	crit := es.Get("Crit")

	if len(crit) < 1 {
		t.Errorf("Error on Get - data retrieved: %s", crit)
	}

	var container person

	err = json.Unmarshal(crit, &container)

	if err != nil {
		t.Errorf("Unmarshal error: %s", err)
	}

	if !(reflect.DeepEqual(td, container)) {
		t.Errorf("Incorrect data returned original key: %v, recieved: %v", td, container)
	}

	es.Delete("Crit")

	crit = es.Get("Crit")

	if len(crit) > 0 {
		t.Errorf("No successful deletion retrieved for:%v ", crit)
	}
}

// need to change to doctype being a part of the new func
// probably should call it category
