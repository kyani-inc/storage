package elasticsearch_test

import (
	"os"
	"testing"

	"encoding/json"
	"github.com/kyani-inc/storage/providers/elasticsearch"
	"github.com/subosito/gotenv"
	"reflect"
)

var (
	host, index, scheme, namespace, awsKey, awsSecret string
)

type product struct {
	Name        string      `json:"name"`
	Price       int         `json:"price"`
	Description description `json:"description"`
}

type description struct {
	Description string `json:"nickname"`
}

func init() {
	gotenv.Load(".env")
	host = os.Getenv("ES_HOST")
	index = os.Getenv("ES_INDEX")
	namespace = os.Getenv("ES_NAMESPACE")
	scheme = os.Getenv("ES_SCHEME")
	awsKey = os.Getenv("AWS_KEY")
	awsSecret = os.Getenv("AWS_SECRET")
}

func emptyVars() bool {
	return host == "" || index == "" || scheme == "" || awsKey == "" || awsSecret == ""

}

func TestES(t *testing.T) {
	if emptyVars() {
		t.Skip("Must have ENV set")
	}

	es, err := elasticsearch.New(host, scheme, index, namespace, awsKey, awsSecret)

	if err != nil {
		t.Errorf("Building new ES client error %+v", err)
	}

	td := product{
		Name:        "Kyäni Sunrise",
		Price:       100,
		Description: description{Description: "Complete Nutrition & Antioxidant Powerhouse"},
	}

	data, err := json.Marshal(td)

	err = es.Put("Sunrise", data)

	if err != nil {
		t.Errorf("Error on Put error:%s", err)
	}

	crit := es.Get("Sunrise")

	if len(crit) < 1 {
		t.Errorf("Error on Get - data retrieved: %s", crit)
	}

	var container product

	err = json.Unmarshal(crit, &container)

	if err != nil {
		t.Errorf("Unmarshal error: %s", err)
	}

	if !(reflect.DeepEqual(td, container)) {
		t.Errorf("Incorrect data returned original key: %v, recieved: %v", td, container)
	}

	es.Delete("Sunrise")

	crit = es.Get("Sunrise")

	if len(crit) > 0 {
		t.Errorf("No successful deletion retrieved for:%v ", crit)
	}
}

// need to change to doctype being a part of the new func
// probably should call it category
