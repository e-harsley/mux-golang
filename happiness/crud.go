package happiness

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type EndpointOption struct {
	allowDelete bool
	allowList   bool
	allowFetch  bool
	allowCreate bool
	allowUpdate bool
	bindOption  []string
}

type SerializerHandler struct {
	createDto      interface{}
	updateDto      IRequest
	response       interface{}
	updateResponse interface{}
	fetchResponse  interface{}
}

func defaultSerializerHandler() *SerializerHandler {
	return &SerializerHandler{}
}

type SerializerHandlerFunc func(*SerializerHandler)

func CreateDto[cdh IRequest](schema cdh) SerializerHandlerFunc {
	return func(opts *SerializerHandler) {
		opts.createDto = schema
	}
}

func UpdateDto[udh IRequest](schema udh) SerializerHandlerFunc {
	return func(opts *SerializerHandler) {
		opts.updateDto = schema
	}
}

func ResponseSchema(schema interface{}) SerializerHandlerFunc {
	return func(opts *SerializerHandler) {
		opts.response = schema
	}
}

func UpdateResponseSchema(schema interface{}) SerializerHandlerFunc {
	return func(opts *SerializerHandler) {
		opts.updateResponse = schema
	}
}

func FetchResponseSchema(schema interface{}) SerializerHandlerFunc {
	return func(opts *SerializerHandler) {
		opts.fetchResponse = schema
	}
}

func defaultEndpointOpts() *EndpointOption {
	return &EndpointOption{
		allowList:   true,
		allowFetch:  true,
		allowCreate: true,
		allowDelete: false,
		allowUpdate: true,
	}
}

func WithoutList(opt *EndpointOption) {
	opt.allowList = false
}

func endpointContextBind(bindContext string) EndpointOptFunc {
	return func(opts *EndpointOption) {
		opts.bindOption = append(opts.bindOption, bindContext)
	}
}

func WithoutFetch(opt *EndpointOption) {
	opt.allowFetch = false
}

func WithoutCreate(opt *EndpointOption) {
	opt.allowCreate = false
}

func WithoutUpdate(opt *EndpointOption) {
	opt.allowUpdate = false
}

func WithDelete(opt *EndpointOption) {
	opt.allowDelete = true
}

type EndpointOptFunc func(*EndpointOption)

// CrudOption is an option to construct the router group.
type CrudOption func(group *mux.Router) *mux.Router

// Action defines a custom action for the CRUD.
type Action struct {
	Name        string
	Handler     http.HandlerFunc
	Middlewares []mux.MiddlewareFunc
	Method      EndpointActionTypes
}

// CustomActions is a slice of custom actions.
type CustomActions []Action

// AddAction adds a custom action to the slice of custom actions.
func (ca *CustomActions) AddAction(action Action) {
	*ca = append(*ca, action)
}

// CrudOptions is a slice of CrudOption functions.
type CrudOptions []CrudOption

// AddOptions adds CrudOption functions to the slice of options.
func (co *CrudOptions) AddOptions(options ...CrudOption) {
	*co = append(*co, options...)
}

func crudOptions[M DocumentModel](repository BaseRepository[M], option EndpointOption, serializer SerializerHandler) CrudOptions {
	var crud CrudOptions
	fmt.Println(option)
	if option.allowList {
		crud = append(crud, getList(repository, serializer.response, option.bindOption...))
	}
	if option.allowFetch {
		crud = append(crud, getOne(repository, serializer.response))
	}
	if option.allowCreate {
		crud = append(crud, create(repository, serializer.createDto, serializer.response))
	}
	if option.allowDelete {
		crud = append(crud, destroy(repository))
	}
	if option.allowUpdate {
		crud = append(crud, update(repository, serializer.createDto, serializer.response))
	}

	return crud
}

// Crud sets up CRUD routes for the given router.
func Crud[M DocumentModel](router *mux.Router, repository BaseRepository[M], relativePath string, options ...interface{}) *mux.Router {
	subRouter := router.PathPrefix(relativePath).Subrouter()

	fmt.Printf("Adding CRUD routes for path: %s\n", relativePath)

	defaultEndpointOptions := defaultEndpointOpts()
	defaultSerializer := defaultSerializerHandler()

	var customActions CustomActions
	var crudOpts CrudOptions

	for _, option := range options {
		switch opt := option.(type) {
		case []mux.MiddlewareFunc:
			subRouter.Use(opt...)
		case func(*EndpointOption):
			opt(defaultEndpointOptions)
		case Action:
			customActions.AddAction(opt)
		case CustomActions:
			for _, customAction := range opt {
				customActions.AddAction(customAction)
			}
		case string:
			if strings.HasPrefix(opt, bindPrefix) {
				customBindFunc := endpointContextBind(opt)
				customBindFunc(defaultEndpointOptions) // Apply the custom binding function to modify options
			}
		case SerializerHandlerFunc:
			opt(defaultSerializer)
		}
	}

	for _, action := range customActions {
		subRouter.Use(action.Middlewares...)
		if action.Method == POST {
			subRouter.HandleFunc(action.Name, action.Handler).Methods("POST")
		}
		if action.Method == PUT {
			subRouter.HandleFunc(action.Name, action.Handler).Methods("PUT")
		}
		if action.Method == DELETE {
			subRouter.HandleFunc(action.Name, action.Handler).Methods("DELETE")
		}
		if action.Method == GET {
			subRouter.HandleFunc(action.Name, action.Handler).Methods("GET")
		}
	}

	crudOpts = append(crudOpts, crudOptions(repository, *defaultEndpointOptions, *defaultSerializer)...)
	for _, opt := range crudOpts {
		subRouter = opt(subRouter)
	}

	return subRouter
}

func getList[M DocumentModel](repository BaseRepository[M], response interface{}, bindFrom ...string) CrudOption {
	return func(router *mux.Router) *mux.Router {

		wrapper := func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("BINDFROM", bindFrom)
			c := C{W: w, R: r}
			paramMap := map[string]interface{}{}
			contextKeysToBind := extractContextBinds(bindFrom)
			for _, cKey := range contextKeysToBind {
				val := c.R.Context().Value(cKey)
				if val == nil {
					fmt.Println("warning!! context key == " + cKey + " does not exist")
					continue
				}

				jsonString, _ := json.Marshal(val)

				err := json.Unmarshal([]byte(jsonString), &paramMap)
				if err != nil {
					http.Error(w, fmt.Sprintf("context key %s made a bad cast: %v", cKey, err), http.StatusInternalServerError)
					return
				}
				fmt.Println("dtoCastFrom", paramMap)
			}

			param := c.Query("filter_by")
			pageSize, pageNumber := int64(10), int64(1)
			sortParam := c.Query("sort_by")
			pageSizeStr, pageNumberStr := c.Query("page_size"), c.Query("page")
			if pageSizeStr != "" {

				lim, _ := strconv.Atoi(pageSizeStr)

				pageSize = int64(lim)

				log.Println(pageSize)

			}
			if pageNumberStr != "" {

				page, _ := strconv.Atoi(pageNumberStr)

				pageNumber = int64(page)

				log.Println(pageNumber)
			}
			skip := (pageNumber * pageSize) - pageSize

			var sortMap map[string]interface{}

			if sortParam == "" {
				sortMap = map[string]interface{}{"created_at": "desc"}
			} else {
				err := json.Unmarshal([]byte(sortParam), &sortMap)
				if err != nil {
					c.Response(400, err.Error(), "Failed to fetch")
					return
				}
			}

			if param != "" {
				err := json.Unmarshal([]byte(param), &paramMap)
				if err != nil {
					c.Response(400, err.Error(), "Failed to fetch")
					return
				}
			}
			params := buildQuery(paramMap)
			sorts := buildSort(sortMap)
			//params = cls.limit(params, c)
			findOptions := options.FindOptions{Limit: &pageSize, Skip: &skip, Sort: sorts}
			models, err := repository.Find(params, &findOptions)
			if err != nil {
				fmt.Println(">>>>>>.")
				c.Response(400, err.Error(), "Failed to fetch")
				return
			}
			count := len(models)

			pagination := PaginationResponse{
				Page:    pageNumber,
				PerPage: pageSize,
				Skip:    skip,
				Count:   count,
			}

			if response == nil {
				c.GetResponse(http.StatusOK, params, sorts, pagination, &models, "FETCHED SUCCESSFULLY")
				return
			}
			resp, err := SerializerFunc(models, &response)
			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "FAILED TO FETCH")
				return
			}

			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "FAILED TO FETCH")
				return
			}
			c.GetResponse(http.StatusOK, params, sorts, pagination, resp, "FETCHED SUCCESSFULLY")

		}

		router.HandleFunc("", wrapper).Methods("GET")
		return router
	}
}

func getOne[M DocumentModel](repository BaseRepository[M], response interface{}) CrudOption {
	return func(router *mux.Router) *mux.Router {
		wrapper := func(w http.ResponseWriter, r *http.Request) {
			c := C{W: w, R: r}

			id := c.Params("id")
			idHex, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "Invalid object id")
				return
			}
			res, err := repository.FindOne(bson.M{"_id": idHex})
			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "VALIDATION ERROR")
				return
			}
			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "Permission ERROR")
				return
			}
			if response == nil {
				c.Response(http.StatusOK, res, "FETCHED SUCCESSFULLY")
				return
			}
			resp, err := SerializerFunc(res, response)

			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "FAILED TO FETCH")
				return
			}
			c.Response(http.StatusOK, resp, "FETCHED SUCCESSFULLY")

		}
		router.HandleFunc("/{id}", wrapper).Methods("GET")
		return router
	}
}

func create[M DocumentModel](repository BaseRepository[M], schema interface{}, response interface{}, bindFrom ...string) CrudOption {
	return func(router *mux.Router) *mux.Router {
		wrapper := func(w http.ResponseWriter, r *http.Request) {
			var (
				err error

				varsToBind = mux.Vars(r)
				canBindAll = len(withoutContexts(bindFrom)) == 0

				c = C{W: w, R: r}
			)

			//fmt.Println("reflect.ValueOf(schema)", reflect.ValueOf(schema))
			if schema == nil {
				c.Response(http.StatusBadRequest, "request coreschema is required", "VALIDATION ERROR")
				return
			}
			schemaType := reflect.TypeOf(schema)

			sche := reflect.New(schemaType).Interface()
			if err := c.BindJSON(sche); err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "VALIDATION ERROR")
				return
			}
			if canBindAll || SliceContains(bindFrom, BindQuery) {
				queryParams := r.URL.Query()
				for key, value := range queryParams {
					varsToBind[key] = value[0]
				}
			}
			for _, cKey := range varsToBind {
				val, exists := varsToBind[cKey]
				if !exists {
					fmt.Println("warning!! context key == " + cKey + " does not exist")
					continue
				}
				err := copier.Copy(&sche, val)
				if err != nil {
					http.Error(w, fmt.Sprintf("context key %s made a bad cast: %v", cKey, err), http.StatusInternalServerError)
					return
				}
			}

			if err := sche.(IRequest).Validate(); err != nil {
				if validationErr, ok := err.(*ValidationError); ok {
					c.Response(http.StatusBadRequest, validationErr.Errors, "VALIDATION ERROR")
					return
				}
				c.Response(http.StatusBadRequest, err.Error(), "VALIDATION ERROR")
				return
			}
			mo, err := repository.BindDataOperationStruct(&sche)
			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "FAILED TO SAVE")
				return
			}
			fmt.Println(mo)
			mode, err := repository.InsertOne(mo)

			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "FAILED TO SAVE")
				return
			}

			fmt.Println(mode)

			if response == nil {

				c.Response(http.StatusCreated, mode, "CREATED SUCCESSFULLY")
				return
			}
			resp, err := SerializerFunc(mode, &response)
			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "FAILED TO SAVE")
				return
			}

			c.Response(http.StatusCreated, resp, "CREATED SUCCESSFULLY")

		}
		router.HandleFunc("", wrapper).Methods("POST")
		return router
	}
}

func update[M DocumentModel](repository BaseRepository[M], schema interface{}, response interface{}, bindFrom ...string) CrudOption {
	return func(router *mux.Router) *mux.Router {
		wrapper := func(w http.ResponseWriter, r *http.Request) {
			var (
				err error

				varsToBind = mux.Vars(r)
				canBindAll = len(withoutContexts(bindFrom)) == 0

				c = C{W: w, R: r}
			)

			id := c.Params("id")
			idHex, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "Invalid object id")
				return
			}
			if schema == nil {
				c.Response(http.StatusBadRequest, "request coreschema is required", "VALIDATION ERROR")
				return
			}
			schemaType := reflect.TypeOf(schema)

			sche := reflect.New(schemaType).Interface()
			if err := c.BindJSON(sche); err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "VALIDATION ERROR")
				return
			}
			if canBindAll || SliceContains(bindFrom, BindQuery) {
				queryParams := r.URL.Query()
				for key, value := range queryParams {
					varsToBind[key] = value[0]
				}
			}

			for _, cKey := range varsToBind {
				val, exists := varsToBind[cKey]
				if !exists {
					fmt.Println("warning!! context key == " + cKey + " does not exist")
					continue
				}
				err := copier.Copy(&sche, val)
				if err != nil {
					http.Error(w, fmt.Sprintf("context key %s made a bad cast: %v", cKey, err), http.StatusInternalServerError)
					return
				}
			}
			if err := sche.(IRequest).Validate(); err != nil {
				if validationErr, ok := err.(*ValidationError); ok {
					c.Response(http.StatusBadRequest, validationErr.Errors, "VALIDATION ERROR")
					return
				}
				c.Response(http.StatusBadRequest, err.Error(), "VALIDATION ERROR")
				return
			}

			mo, err := repository.BindDataOperationStruct(&sche)
			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "FAILED TO SAVE")
				return
			}
			mode, err := repository.FindOneAndUpdate(bson.M{"_id": idHex}, mo)
			if response == nil {

				c.Response(http.StatusCreated, mode, "CREATED SUCCESSFULLY")
				return
			}
			resp, err := SerializerFunc(mode, &response)
			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "FAILED TO SAVE")
				return
			}

			c.Response(http.StatusCreated, resp, "CREATED SUCCESSFULLY")

		}
		router.HandleFunc("/{id}", wrapper).Methods("PUT")
		return router
	}
}

func destroy[M DocumentModel](repository BaseRepository[M]) CrudOption {
	return func(router *mux.Router) *mux.Router {
		wrapper := func(w http.ResponseWriter, r *http.Request) {
			c := C{W: w, R: r}
			objID := c.Params("objID")
			id, err := primitive.ObjectIDFromHex(objID)
			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "Invalid object id")
				return
			}
			_, err = repository.DeleteOne(bson.M{"_id": id})
			if err != nil {
				c.Response(http.StatusBadRequest, err.Error(), "FAILED TO SAVE")
				return
			}
			c.Response(http.StatusOK, map[string]interface{}{
				"message": fmt.Sprintf("%s deleted successfully", objID),
			}, "DELETED SUCCESSFULLY")
			return

		}
		router.HandleFunc("/{id}", wrapper).Methods("DELETE")
		return router
	}
}

func Depend[T IRequest](handler func(req T, c C) *Response, bindFrom ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			dtoCastFrom       T
			err               error
			contextKeysToBind = extractContextBinds(bindFrom)
			varsToBind        = mux.Vars(r)
			canBindAll        = len(withoutContexts(bindFrom)) == 0

			c = C{W: w, R: r}
		)

		if canBindAll || SliceContains(bindFrom, BindJSON) {
			err = json.NewDecoder(r.Body).Decode(&dtoCastFrom)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		if canBindAll || SliceContains(bindFrom, BindQuery) {
			queryParams := r.URL.Query()
			for key, value := range queryParams {
				varsToBind[key] = value[0]
			}
		}

		for _, cKey := range contextKeysToBind {
			val := c.R.Context().Value(cKey)
			if val == nil {
				fmt.Println("warning!! context key == " + cKey + " does not exist")
				continue
			}
			fmt.Println(val)
			err := copier.Copy(&dtoCastFrom, val)
			if err != nil {
				http.Error(w, fmt.Sprintf("context key %s made a bad cast: %v", cKey, err), http.StatusInternalServerError)
				return
			}
			fmt.Println("dtoCastFrom", dtoCastFrom)
		}

		err = dtoCastFrom.Validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res := handler(dtoCastFrom, c)
		if res == nil {
			http.Error(w, "No response in body", http.StatusInternalServerError)
			return
		}
		responseCode := res.HTTPStatusCode
		if responseCode == 0 {
			responseCode = MethodToStatusCode[r.Method]
			res.HTTPStatusCode = responseCode
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseCode)
		json.NewEncoder(w).Encode(res.JSONDATA())

	}
}
