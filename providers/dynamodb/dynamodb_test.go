package dynamodb_test

import (
	"os"
	"testing"

	"github.com/kyani-inc/storage/providers/dynamodb"
	"github.com/subosito/gotenv"
)

var (
	ddb       dynamodb.DynamoDB
	envLoaded = false
	region    = ""
	dbtable   = ""
	endpoint  = ""
)

func checkEnv(t *testing.T) {
	if envLoaded == false {
		gotenv.Load(".env")
		region = os.Getenv("AWS_REGION") //"us-west-2"
		dbtable = os.Getenv("DYNAMO_DB_TABLE")
		endpoint = os.Getenv("DYNAMO_DB_ENDPOINT")
		envLoaded = true
	}

	if region == "" || dbtable == "" || endpoint == "" {
		t.Skip("env vars not defined")
	}
}

func TestDDB(t *testing.T) {
	checkEnv(t)

	var err error
	ddb, err = dynamodb.New(region, endpoint, dbtable)

	if err != nil {
		t.Fatal("Failed to establish connection with DynamoDB!")
	} else {
		t.Log("Connected to local DynamoDB server")
	}

	k, v := "test1", "hello, world!!"
	err = ddb.Put(k, []byte(v))

	if err != nil {
		t.Error("Error putting value", err.Error())
	}

	b := ddb.Get(k)

	if v != string(b) {
		t.Error("item `test1` does not contain expected values")
	}

	ddb.Delete(k)

	b = ddb.Get(k)

	if v == string(b) {
		t.Error("key test2 was not deleted!")
	}

	k, v = "nodata", ""
	err = ddb.Put(k, []byte(v))

	if err == nil {
		t.Error("An error was expected but passed for some reason..")
	}

	b = ddb.Get("nodata")

	if v != string(b) {
		t.Error("item `nodata` should not contain data")
	}

	k, v = "test2", "data..."
	err = ddb.Put(k, []byte(v))

	if err != nil {
		t.Error("Error putting value", err.Error())
	}

	ddb.Flush()

	b = ddb.Get(k)

	if v == string(b) {
		t.Error("Failed to flush the table..")
	}
}
