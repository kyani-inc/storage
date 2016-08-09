package elasticsearch

import (
	"encoding/json"
	"gopkg.in/olivere/elastic.v3"
	"net/http"
)

type ElasticSearch struct {
	client    *elastic.Client
	index     string
	namespace string
}

// New creates a new instance of an ES client
// We use http.Client as a param so that we can
// be provider agnostic - eg. AWS ElasticSearch
func New(host, scheme, index, namespace string, h http.Client) (ElasticSearch, error) {
	var e ElasticSearch

	client, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetMaxRetries(10),
		elastic.SetScheme(scheme),
		elastic.SetHttpClient(&h),
	)

	if err != nil {
		return e, err
	}

	// the namespace becomes the doctype for ease of use
	e.client = client
	e.index = index
	e.namespace = namespace

	return e, err
}

// Put overwrites or creates a new ES entry based on key, index, and doctype
func (e ElasticSearch) Put(key string, data []byte) error {
	// check for the index, create if does not exist
	exists, err := e.client.IndexExists(e.index).Do()

	if err != nil {
		return err
	}

	if !exists {
		createIndex, err := e.client.CreateIndex(e.index).Do()

		if err != nil {
			return err
		}

		// it is possible that no error is thrown the index creation failed
		if !createIndex.Acknowledged {
			return err
		}
	}

	var container interface{}

	err = json.Unmarshal(data, &container)

	if err != nil {
		return err
	}

	_, err = e.client.Index().
		Index(e.index).
		Type(e.namespace).
		Id(key).
		BodyJson(container).
		Do()

	if err != nil {
		return err
	}

	return nil
}

// Get attempts to grab ES entries by key, index and doctype
func (e ElasticSearch) Get(key string) []byte {
	data, err := e.client.Get().
		Index(e.index).
		Type(e.namespace).
		Id(key).
		Do()

	if err != nil || data.Id != key {
		return []byte{}
	}

	res, err := json.Marshal(data.Source)

	if err != nil {
		return []byte{}
	}

	return res
}

// Delete removes an ES entry by key
func (e ElasticSearch) Delete(key string) {
	e.client.Delete().
		Index(e.index).
		Type(e.namespace).
		Id(key).
		Do()
}

// Flush tells ES to free memory from the index and flush data to disk.
func (e ElasticSearch) Flush() {
	e.client.Flush().
		Index(e.index).
		AllowNoIndices(true).
		IgnoreUnavailable(true).
		Do()
}
