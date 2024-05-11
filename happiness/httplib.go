package happiness

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

const (
	// BindJSON : if you want to specify specific json bind for your easy gin handler
	BindJSON = "json_bind"
	// BindQuery : if you want to specify specific query bind for your easy gin handler
	BindQuery = "query_bind"
	// BindURI : if you want to specify specific uri bind for your easy gin handler
	BindURI = "uri_bind"

	bindContext = "context_bind"
	ctxSep      = ":::"
)

type PaginationResponse struct {
	Page    int64 `json:"page"`
	PerPage int64 `json:"perPage"`
	Skip    int64 `json:"skip"`
	Count   int   `json:"count"`
}

type GeneralResponse struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
type GeneralGetResponse struct {
	Status     int                    `json:"status"`
	Pagination PaginationResponse     `json:"pagination"`
	FilterBy   map[string]interface{} `json:"filter_by"`
	SortBy     bson.D                 `json:"sort_by"`
	Data       interface{}            `json:"data,omitempty"`
	Message    string                 `json:"message,omitempty"`
}

type C struct {
	W http.ResponseWriter
	R *http.Request
}

func (c *C) BindJSON(data interface{}) error {
	// Read the request body and decode JSON into the provided interface
	if err := json.NewDecoder(c.R.Body).Decode(data); err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}
	return nil
}

func responseJSON(res http.ResponseWriter, status int, object interface{}) {
	res.Header().Set("Content-Resource", "application/json")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	err := json.NewEncoder(res).Encode(object)

	if err != nil {
		return
	}
}

func (c *C) Query(key string) string {
	query := c.R.URL.Query()
	return query.Get(key)
}

func (c *C) Params(key string) string {
	return mux.Vars(c.R)[key]
}

func (c *C) GetResponse(status int, filterBy map[string]interface{}, sort bson.D, pagination PaginationResponse, data interface{}, message string) {
	responseSuccess := GeneralGetResponse{
		FilterBy:   filterBy,
		Pagination: pagination,
		SortBy:     sort,
		Status:     status,
		Message:    message,
		Data:       data,
	}
	responseJSON(c.W, status, responseSuccess)
}

func (c *C) Response(status int, data interface{}, message string) {
	responseSuccess := GeneralResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
	responseJSON(c.W, status, responseSuccess)
}

type HTTPResponse struct {
	StatusCode int
	Body       []byte
}

type RequestOptions struct {
	Header map[string]string
	Params map[string]string
}

func (r *HTTPResponse) JSONBody(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

func MakePostRequest(url string, payload interface{}, headers ...map[string]string) (*HTTPResponse, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}

	// Add headers if provided
	for _, headerMap := range headers {
		for key, value := range headerMap {
			fmt.Println(key, value)
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("here cool ", resp.StatusCode)
	fmt.Println(resp.Body)

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       responseBody,
	}, nil
}

func MakeGetRequest(url string, opts ...*RequestOptions) (*HTTPResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if len(opts) > 0 {
		opt := opts[0]
		for key, value := range opt.Params {
			query.Add(key, value)
		}
		req.URL.RawQuery = query.Encode()

		for key, value := range opt.Header {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       responseBody,
	}, nil
}

func buildQuery(filter map[string]interface{}) bson.M {
	query := bson.M{}
	fmt.Println(query)
	for key, value := range filter {
		fmt.Println("????????????", key, value)
		switch key {
		case "$or", "$and":
			subqueries := make([]bson.M, 0)
			subfilter, ok := value.([]map[string]interface{})
			if !ok {
				continue
			}

			for _, sub := range subfilter {
				subqueries = append(subqueries, buildQuery(sub))
			}

			if len(subqueries) > 0 {
				query[key] = subqueries
			}
		default:
			valueType := reflect.TypeOf(value)

			if valueType.Kind() == reflect.Map {
				subquery := buildQuery(value.(map[string]interface{}))
				query[key] = subquery
			} else {
				query[key] = value
			}
		}
	}
	fmt.Println("query", query)

	return query
}

func buildSort(m map[string]interface{}) bson.D {
	var bsonD bson.D
	fmt.Println(bson.D{{Key: "created_at", Value: -1}})
	for key, value := range m {
		sortOrder := 1
		if s, ok := value.(string); ok && strings.ToLower(s) == "desc" {
			sortOrder = -1 // Set to descending
		}
		bsonD = append(bsonD, bson.E{Key: key, Value: sortOrder})
	}
	return bsonD
}

func SerializerFunc(result interface{}, obj interface{}) (interface{}, error) {

	resultType := reflect.TypeOf(result)
	objType := reflect.TypeOf(obj)
	schemaInstance := reflect.New(reflect.TypeOf(obj)).Interface()

	hydrate := reflect.ValueOf(schemaInstance).MethodByName("Hydrate")

	if hydrate.IsValid() {
		newObj := map[string]interface{}{}
		if resultType.Kind() == reflect.Slice {
			var resp []interface{}
			resultValue := reflect.ValueOf(result)
			for i := 0; i < resultValue.Len(); i++ {
				element := resultValue.Index(i)
				elementJSON, err := json.Marshal(element.Interface())
				if err != nil {
					return nil, err
				}
				var newItem map[string]interface{}
				err = json.Unmarshal(elementJSON, &newItem)
				if err != nil {
					return nil, err
				}
				hydrateResult := hydrate.Call([]reflect.Value{reflect.ValueOf(newItem)})

				if len(hydrateResult) < 1 {
					return nil, errors.New("hydrate should return an interface")
				}
				resp = append(resp, hydrateResult[0].Interface())
			}
			result = resp
		} else {
			elementJSON, err := json.Marshal(result)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(elementJSON, &newObj)
			if err != nil {
				return nil, err
			}
			hydrateResult := hydrate.Call([]reflect.Value{reflect.ValueOf(newObj)})

			if len(hydrateResult) < 1 {
				return nil, errors.New("hydrate should return an interface")
			}
			result = hydrateResult[0].Interface()
		}

	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	if resultType.Kind() == reflect.Slice {
		sliceType := reflect.SliceOf(objType)
		newObj := reflect.New(sliceType).Interface()
		if err := json.Unmarshal(resultJSON, newObj); err != nil {
			return nil, err
		}
		return newObj, nil
	} else {
		newObj := reflect.New(objType).Interface()
		if err := json.Unmarshal(resultJSON, newObj); err != nil {
			return nil, err
		}
		return newObj, nil
	}
}

func extractContextBinds(bindFroms []string) (ctxKeys []string) {
	for _, v := range bindFroms {
		if strings.Contains(v, bindContext) {
			ctxKeys = append(ctxKeys, strings.Split(v, ctxSep)[1])
		}
	}
	return ctxKeys
}

func BindContext(ctxKey string) string {
	return bindContext + ctxSep + ctxKey
}

func withoutContexts(bindFroms []string) (nonCtx []string) {
	for _, v := range bindFroms {
		if !strings.Contains(v, bindContext) {
			nonCtx = append(nonCtx, v)
		}
	}
	return nonCtx
}
