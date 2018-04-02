package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/tomogoma/go-typed-errors"
	"github.com/tomogoma/shoppingms/pkg/config"
	"github.com/tomogoma/shoppingms/pkg/logging"
	"strings"
	"io/ioutil"
	"github.com/tomogoma/shoppingms/pkg/shopping"
	"github.com/tomogoma/crdb"
	"strconv"
)

type contextKey string

type Guard interface {
	APIKeyValid(key []byte) (string, error)
}

type ShoppingManager interface {
	errors.ToHTTPResponser
	InsertShoppingList(JWT, name, mode string) (*shopping.ShoppingList, error)
	UpdateShoppingList(JWT, shoppingListID string, name, mode crdb.StringUpdate) (*shopping.ShoppingList, error)
	ShoppingLists(JWT string, offset, count int64) ([]shopping.ShoppingList, error)
	ShoppingListItems(JWT, shoppingListID string, offset, count int64) ([]shopping.ShoppingList, error)
}

type handler struct {
	errors.ErrToHTTP

	guard   Guard
	logger  logging.Logger
	manager ShoppingManager
}

type Config struct {
	Guard          Guard
	Logger         logging.Logger
	BaseURL        string
	AllowedOrigins []string
	Manager        ShoppingManager
}

const (
	keyAPIKey = "x-api-key"

	ctxKeyLog = contextKey("log")
)

func NewHandler(conf Config) (http.Handler, error) {
	if conf.Guard == nil {
		return nil, errors.New("Guard was nil")
	}
	if conf.Logger == nil {
		return nil, errors.New("Logger was nil")
	}
	if conf.Manager == nil {
		return nil, errors.New("ShoppingManager was nil")
	}

	r := mux.NewRouter().PathPrefix(conf.BaseURL).Subrouter()
	handler{guard: conf.Guard, logger: conf.Logger, manager: conf.Manager}.handleRoute(r)

	corsOpts := []handlers.CORSOption{
		handlers.AllowedHeaders([]string{
			"X-Requested-With", "Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Authorization", "X-api-key",
		}),
		handlers.AllowedOrigins(conf.AllowedOrigins),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
	}
	return handlers.CORS(corsOpts...)(r), nil
}

func (s handler) handleRoute(r *mux.Router) {
	s.handleStatus(r)
	s.handleDocs(r)

	s.handleNewShoppingList(r)
	s.handleUpdateShoppingList(r)
	s.handleGetShoppingLists(r)

	s.handleUpsertShoppingListItem(r)
	s.handleDeleteShoppingListItem(r)
	s.handleGetShoppingListItems(r)
	s.handleSearchShoppingItems(r)

	s.handleNotFound(r)
}

/**
 * @api {get} /status Status
 * @apiName Status
 * @apiVersion 0.1.0
 * @apiGroup Service
 *
 * @apiHeader x-api-key the api key
 *
 * @apiSuccess (200) {String} name Micro-service name.
 * @apiSuccess (200)  {String} version http://semver.org version.
 * @apiSuccess (200)  {String} description Short description of the micro-service.
 * @apiSuccess (200)  {String} canonicalName Canonical name of the micro-service.
 *
 */
func (s *handler) handleStatus(r *mux.Router) {
	r.Methods(http.MethodGet).
		PathPrefix("/status").
		HandlerFunc(
		s.apiGuardChain(func(w http.ResponseWriter, r *http.Request) {
			s.respondJsonOn(w, r, nil, struct {
				Name          string `json:"name"`
				Version       string `json:"version"`
				Description   string `json:"description"`
				CanonicalName string `json:"canonicalName"`
			}{
				Name:          config.Name,
				Version:       config.VersionFull,
				Description:   config.Description,
				CanonicalName: config.CanonicalWebName(),
			}, http.StatusOK, nil, s)
		}),
	)
}

/**
 * @api {get} /docs Docs
 * @apiName Docs
 * @apiVersion 0.1.0
 * @apiGroup Service
 *
 * @apiSuccess (200) {html} docs Docs page to be viewed on browser.
 *
 */
func (s *handler) handleDocs(r *mux.Router) {
	r.PathPrefix("/" + config.DocsPath).
		Handler(http.FileServer(http.Dir(config.DefaultDocsDir())))
}

/**
 * @api {put} /shoppinglists New Shopping List
 * @apiName InsertShoppingList
 * @apiVersion 0.1.0
 * @apiGroup Service
 * @apiDescription insert a shopping list by name if not exists.
 *
 * @apiHeader x-api-key the api key
 * @apiHeader Authorization Bearer token received from authentication
 * 		micro-service in the form "Bearer {token-value}".
 *
 * @apiParam (JSON Request Body) {String} name The name of the new shopping list.
 * @apiParam (JSON Request Body) {String="PREPARATION","SHOPPING"} [mode="PREPARATION"]
 * 		The current mode of the shopping list on the client apps.
 *
 * @apiUse ShoppingList200
 *
 */
func (s *handler) handleNewShoppingList(r *mux.Router) {
	r.Methods(http.MethodPut).
		PathPrefix("/shoppinglists").
		HandlerFunc(
		s.apiGuardChain(func(w http.ResponseWriter, r *http.Request) {

			req := struct {
				JWT  string
				Name string
				Mode string
			}{}

			if err := readJSONBody(r, &req); err != nil {
				handleError(w, r, req, err, s)
				return
			}

			var err error
			if req.JWT, err = readJWT(r); err != nil {
				handleError(w, r, req, err, s)
				return
			}

			sl, err := s.manager.InsertShoppingList(req.JWT, req.Name, req.Mode)
			s.respondJsonOn(w, r, req, NewShoppingList(sl), http.StatusOK, err, s.manager)
		}),
	)
}

/**
 * @api {put} /shoppinglists/{ID} Update Shopping List
 * @apiName UpdateShoppingList
 * @apiVersion 0.1.0
 * @apiGroup Service
 * @apiDescription update a shopping list with {ID}.
 *
 * @apiHeader x-api-key the api key
 * @apiHeader Authorization Bearer token received from authentication
 * 		micro-service in the form "Bearer {token-value}".
 *
 * @apiParam (URL Path Params) {String} id The ID of the shopping list.
 *
 * @apiParam (JSON Request Body) {String} name
 * 		Unique name of the shopping list.
 * @apiParam (JSON Request Body) {String="PREPARATION","SHOPPING"} mode
 * 		The current mode of the shopping list on the client apps.
 *
 * @apiUse ShoppingList200
 *
 */
func (s *handler) handleUpdateShoppingList(r *mux.Router) {
	r.Methods(http.MethodPut).
		PathPrefix("/shoppinglists/{ID}").
		HandlerFunc(
		s.apiGuardChain(func(w http.ResponseWriter, r *http.Request) {

			req := struct {
				JWT            string
				ShoppingListID string
				Name           JSONStringUpdate
				Mode           JSONStringUpdate
			}{}

			if err := readJSONBody(r, &req); err != nil {
				handleError(w, r, req, err, s)
				return
			}

			req.ShoppingListID = mux.Vars(r)["ID"]

			var err error
			if req.JWT, err = readJWT(r); err != nil {
				handleError(w, r, req, err, s)
				return
			}

			sl, err := s.manager.UpdateShoppingList(req.JWT, req.ShoppingListID,
				req.Name.StringUpdate, req.Mode.StringUpdate)

			s.respondJsonOn(w, r, req, NewShoppingList(sl), http.StatusOK, err, s.manager)
		}),
	)
}

/**
 * @api {get} /shoppinglists Get Shopping Lists
 * @apiName GetShoppingLists
 * @apiVersion 0.1.0
 * @apiGroup Service
 * @apiDescription Get shopping lists for a user.
 *
 * @apiHeader x-api-key the api key
 * @apiHeader Authorization Bearer token received from authentication
 * 		micro-service in the form "Bearer {token-value}".
 *
 * @apiParam (URL Query Params) {Long} [offset=0]
 * 		Offset index to fetch from.
 * @apiParam (URL Query Params) {Long} [count=10]
 * 		Number of shopping lists to fetch.
 *
 * @apiUse ShoppingLists200
 *
 */
func (s *handler) handleGetShoppingLists(r *mux.Router) {
	r.Methods(http.MethodGet).
		PathPrefix("/shoppinglists").
		HandlerFunc(
		s.apiGuardChain(func(w http.ResponseWriter, r *http.Request) {

			req := struct {
				JWT    string
				Offset int64
				Count  int64
			}{}

			var err error
			if req.JWT, err = readJWT(r); err != nil {
				handleError(w, r, req, err, s)
				return
			}

			if req.Offset, err = readOffset(r); err != nil {
				handleError(w, r, req, err, s)
				return
			}

			if req.Count, err = readCount(r); err != nil {
				handleError(w, r, req, err, s)
				return
			}

			sls, err := s.manager.ShoppingLists(req.JWT, req.Offset, req.Count)
			s.respondJsonOn(w, r, req, NewShoppingLists(sls), http.StatusOK, err, s.manager)
		}),
	)
}

/**
 * @api {put} /shoppinglists/{ID}/items Upsert Shopping List Item
 * @apiName UpsertShoppingListItem
 * @apiVersion 0.1.0
 * @apiGroup Service
 * @apiDescription Update/Insert a list item's values for a shopping list.
 *		Note that all details under the Price object are shared with
 * 		other users and will not be deleted during item deletion.
 *
 * @apiHeader x-api-key the api key
 * @apiHeader Authorization Bearer token received from authentication
 * 		micro-service in the form "Bearer {token-value}".
 *
 * @apiParam (URL Path Params) {String} id The ID of the shopping list.
 *
 * @apiParam (JSON Request Body) {String} itemName
 * 		Name of the ShoppingItem e.g. Toothpaste.
 * @apiParam (JSON Request Body) {Boolean} [inList]
 * 		True if item is in the shopping list, false otherwise.
 * @apiParam (JSON Request Body) {Boolean} [inCart]
 * 		True if item is in the shopping cart, false otherwise. Marking this as
 * 		true automatically sets inList to true.
 * @apiParam (JSON Request Body) {String} [brandName]
 * 		Name of the Brand of the itemName e.g. Colgate.
 * @apiParam (JSON Request Body) {Int} [quantity]
 * 		Number of items in the shopping List.
 * @apiParam (JSON Request Body) {String} [measurementUnit]
 * 		The measurement Unit to use e.g. 250ml Tub, KG, 5Kg bag, etc.
 * @apiParam (JSON Request Body) {Float} [unitPrice]
 * 		Price of one unit of measurement e.g. 200 if a 250ml Tub costs that.
 * @apiParam (JSON Request Body) {String} [currency=KES]
 *		Active ISO 4217 code denoting currency of the unitPrice.
 *
 * @apiUse ShoppingListItem200
 *
 */
func (s *handler) handleUpsertShoppingListItem(r *mux.Router) {
	r.Methods(http.MethodPut).
		PathPrefix("/shoppinglists/{ID}/items").
		HandlerFunc(
		s.apiGuardChain(func(w http.ResponseWriter, r *http.Request) {
			// TODO()
			handleError(w, r, r, errors.NewNotImplemented(), s)
		}),
	)
}

/**
 * @api {delete} /items/{ID} Delete Shopping List Item
 * @apiName DeleteShoppingListItem
 * @apiVersion 0.1.0
 * @apiGroup Service
 * @apiDescription Delete a shopping list Item. Note that this only deletes
 * 		the top level shopping list, the price details remain intact.
 *
 * @apiHeader x-api-key the api key
 * @apiHeader Authorization Bearer token received from authentication
 * 		micro-service in the form "Bearer {token-value}".
 *
 * @apiParam (URL Path Params) {String} id
 * 		The ID of the shopping list item to delete.
 *
 * @apiSuccess (200) emptyBody check status code for success.
 *
 */
func (s *handler) handleDeleteShoppingListItem(r *mux.Router) {
	r.Methods(http.MethodDelete).
		PathPrefix("/items/{ID}").
		HandlerFunc(
		s.apiGuardChain(func(w http.ResponseWriter, r *http.Request) {
			// TODO()
			handleError(w, r, r, errors.NewNotImplemented(), s)
		}),
	)
}

/**
 * @api {get} /shoppinglists/{ID}/items Get Shopping List Items
 * @apiName GetShoppingListItems
 * @apiVersion 0.1.0
 * @apiGroup Service
 * @apiDescription Get shopping items for a shopping list.
 *
 * @apiHeader x-api-key the api key
 * @apiHeader Authorization Bearer token received from authentication
 * 		micro-service in the form "Bearer {token-value}".
 *
 * @apiParam (URL Query Params) {Long} [offset=0]
 * 		Offset index to fetch from.
 * @apiParam (URL Query Params) {Long} [count=10]
 * 		Number of shopping lists to fetch.
 *
 * @apiSuccess (200 JSON Response Body) {Object[]} items
 *		List of ShoppingListItems. See "200 JSON Response Body" of
 *		<a href="#api-Service-UpsertShoppingListItem">Upsert Shopping List Item</a>
 *		for details on what each item looks like.
 *
 */
func (s *handler) handleGetShoppingListItems(r *mux.Router) {
	r.Methods(http.MethodGet).
		PathPrefix("/shoppinglists/{ID}/items").
		HandlerFunc(
		s.apiGuardChain(func(w http.ResponseWriter, r *http.Request) {

			req := struct {
				JWT            string
				ShoppingListID string
				Offset         int64
				Count          int64
			}{}

			req.ShoppingListID = mux.Vars(r)["ID"]

			var err error
			if req.JWT, err = readJWT(r); err != nil {
				handleError(w, r, req, err, s)
				return
			}

			if req.Offset, err = readOffset(r); err != nil {
				handleError(w, r, req, err, s)
				return
			}

			if req.Count, err = readCount(r); err != nil {
				handleError(w, r, req, err, s)
				return
			}

			sls, err := s.manager.ShoppingListItems(req.JWT, req.ShoppingListID, req.Offset, req.Count)
			s.respondJsonOn(w, r, req, NewShoppingLists(sls), http.StatusOK, err, s.manager)
		}),
	)
}

/**
 * @api {get} /items/search Search Shopping Items
 * @apiName SearchShoppingItems
 * @apiVersion 0.1.0
 * @apiGroup Service
 * @apiDescription Search Shopping Items not necessarily belonging to a
 *		specific ShoppingList.
 *
 * @apiHeader x-api-key the api key
 * @apiHeader Authorization Bearer token received from authentication
 * 		micro-service in the form "Bearer {token-value}".
 *
 * @apiParam (URL Query Params) {Long} [offset=0]
 * 		Offset index to fetch from.
 * @apiParam (URL Query Params) {Long} [count=10]
 * 		Number of shopping lists to fetch.
 * @apiParam (URL Query Params) {String} [brandName]
 * 		If provided, filter items where brandName contains provided text.
 * @apiParam (URL Query Params) {String} [itemName]
 * 		If provided, filter items where itemName contains provided text.
 * @apiParam (URL Query Params) {String} [brandPrice]
 * 		If provided, filter items where brandPrice contains provided text.
 * @apiParam (URL Query Params) {String} [measuringUnit]
 * 		If provided, filter items where measuringUnit contains provided text.
 *
 * @apiSuccess (200 JSON Response Body) {Object[]} items
 *		List of ShoppingListItems. See "200 JSON Response Body" of
 *		<a href="#api-Service-UpsertShoppingListItem">Upsert Shopping List Item</a>
 *		for details on what each item looks like.
 *
 */
func (s *handler) handleSearchShoppingItems(r *mux.Router) {
	r.Methods(http.MethodGet).
		PathPrefix("/items/search").
		HandlerFunc(
		s.apiGuardChain(func(w http.ResponseWriter, r *http.Request) {
			// TODO()
			handleError(w, r, r, errors.NewNotImplemented(), s)
		}),
	)
}

func (s handler) handleNotFound(r *mux.Router) {
	r.NotFoundHandler = http.HandlerFunc(
		s.prepLogger(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Nothing to see here", http.StatusNotFound)
		}),
	)
}

func (s *handler) apiGuardChain(next http.HandlerFunc) http.HandlerFunc {
	return s.prepLogger(s.guardRoute(next))
}

func (s handler) prepLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log := s.logger.WithHTTPRequest(r).
			WithField(logging.FieldTransID, uuid.New())

		log.WithFields(map[string]interface{}{
			logging.FieldURLPath:    r.URL.Path,
			logging.FieldHTTPMethod: r.Method,
		}).Info("new request")

		ctx := context.WithValue(r.Context(), ctxKeyLog, log)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (s *handler) guardRoute(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		APIKey := r.Header.Get(keyAPIKey)
		clUsrID, err := s.guard.APIKeyValid([]byte(APIKey))
		log := r.Context().Value(ctxKeyLog).(logging.Logger).
			WithField(logging.FieldClientAppUserID, clUsrID)
		ctx := context.WithValue(r.Context(), ctxKeyLog, log)
		if err != nil {
			handleError(w, r.WithContext(ctx), nil, err, s)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// respondJsonOn marshals respData to json and writes it and the code as the
// http header to w. If err is not nil, handleError is called instead of the
// documented write to w.
func (s *handler) respondJsonOn(w http.ResponseWriter, r *http.Request, reqData interface{},
	respData interface{}, code int, err error, errSrc errors.ToHTTPResponser) int {

	if err != nil {
		handleError(w, r, reqData, err, errSrc)
		return 0
	}

	respBytes, err := json.Marshal(respData)
	if err != nil {
		handleError(w, r, reqData, err, errSrc)
		return 0
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	i, err := w.Write(respBytes)
	if err != nil {
		log := r.Context().Value(ctxKeyLog).(logging.Logger)
		log.Errorf("unable write data to response stream: %v", err)
		return i
	}

	return i
}

// handleError writes an error to w using errSrc's logic and logs the error
// using the logger acquired by the prepLogger middleware on r. reqData is
// included in the log data.
func handleError(w http.ResponseWriter, r *http.Request, reqData interface{}, err error, errSrc errors.ToHTTPResponser) {
	reqDataB, _ := json.Marshal(reqData)
	log := r.Context().Value(ctxKeyLog).(logging.Logger).
		WithField(logging.FieldRequest, string(reqDataB))

	if code, ok := errSrc.ToHTTPResponse(err, w); ok {
		log.WithField(logging.FieldResponseCode, code).Warn(err)
		return
	}

	log.WithField(logging.FieldResponseCode, http.StatusInternalServerError).
		Error(err)
	http.Error(w, "Something wicked happened, please try again later",
		http.StatusInternalServerError)
}

func readJSONBody(r *http.Request, into interface{}) error {
	bodyB, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.NewClientf("read request body: %v", err)
	}
	defer r.Body.Close()
	if err = json.Unmarshal(bodyB, into); err != nil {
		return errors.NewClientf("invalid json body: %v", err)
	}
	return nil
}

func readJWT(r *http.Request) (string, error) {
	bearerPrfx := "bearer "
	bearerPrfxLen := len(bearerPrfx)
	for _, val := range r.Header["Authorization"] {
		if !strings.Contains(strings.ToLower(val), bearerPrfx) {
			continue
		}
		return val[bearerPrfxLen:], nil
	}
	return "", errors.NewUnauthorized("\"Authorization\" header with \"Bearer \" token not found")
}

func readOffset(r *http.Request) (int64, error) {
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		return 0, nil
	}
	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		return -1, errors.NewClientf("invalid offset: %v", err)
	}
	return offset, nil
}

func readCount(r *http.Request) (int64, error) {
	countStr := r.URL.Query().Get("count")
	if countStr == "" {
		return 10, nil
	}
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		return -1, errors.NewClientf("invalid count: %v", err)
	}
	return count, nil
}
