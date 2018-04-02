package http

import (
	"github.com/tomogoma/shoppingms/pkg/shopping"
	"github.com/tomogoma/crdb"
	"encoding/json"
)

type JSONStringUpdate struct {
	crdb.StringUpdate
}

func (jsu *JSONStringUpdate) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set in the JSON string.
	if string(data) == "null" { // Ignore the null literal value.
		return nil
	}

	jsu.Updating = true
	if err := json.Unmarshal(data, &jsu.NewVal); err != nil {
		return err
	}
	return nil
}

/**
 * @apiDefine ShoppingLists200
 * @apiSuccess (200 JSON Response Body) {Object[]} shoppingLists
 * 		List of shoppingLists.
 * @apiSuccess (200 JSON Response Body) {String} shoppingLists.ID
 * 		Unique ID of the shopping list.
 * @apiSuccess (200 JSON Response Body) {String} shoppingLists.userID
 * 		ID of the user who owns the shopping list.
 * @apiSuccess (200 JSON Response Body) {String} shoppingLists.name
 * 		Unique name of the shopping list.
 * @apiSuccess (200 JSON Response Body) {String="PREPARATION","SHOPPING"} shoppingLists.mode
 * 		The current mode of the shopping list on the client apps.
 * @apiSuccess (200 JSON Response Body) {String} shoppingLists.created
 * 		ISO8601 date of shopping list creation.
 * @apiSuccess (200 JSON Response Body) {String} shoppingLists.lastUpdated
 * 		ISO8601 date denoting last time the list was updated.
 */
/**
 * @apiDefine ShoppingList200
 * @apiSuccess (200 existed JSON Response Body) {String} ID
 *		Unique ID of the shopping list.
 * @apiSuccess (200 existed JSON Response Body) {String} userID
 *		ID of the user who owns the shopping list.
 * @apiSuccess (200 existed JSON Response Body) {String} name
 *	 	Unique name of the shopping list.
 * @apiSuccess (200 existed JSON Response Body) {String="PREPARATION","SHOPPING"} mode
 * 		The current mode of the shopping list on the client apps.
 * @apiSuccess (200 existed JSON Response Body) {String} created
 *		ISO8601 date of shopping list creation.
 * @apiSuccess (200 existed JSON Response Body) {String} lastUpdated
 * 		ISO8601 date denoting last time the list was updated.
 */
type ShoppingList struct {
	ID          string `json:"ID,omitempty"`
	UserID      string `json:"userID,omitempty"`
	Name        string `json:"name,omitempty"`
	Mode        string `json:"mode,omitempty"`
	Created     string `json:"created,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}

func NewShoppingList(list *shopping.ShoppingList) *ShoppingList {
	if list == nil {
		return nil
	}
	return &ShoppingList{
		ID:          list.ID,
		UserID:      list.UserID,
		Name:        list.Name,
		Mode:        list.Mode,
		Created:     list.Created,
		LastUpdated: list.LastUpdated,
	}
}

func NewShoppingLists(lists []shopping.ShoppingList) []ShoppingList {
	if len(lists) == 0 {
		return nil
	}
	var ress []ShoppingList
	for _, list := range lists {
		res := NewShoppingList(&list)
		ress = append(ress, *res)
	}
	return ress
}

type MeasuringUnit struct {
	ID   string `json:"ID,omitempty"`
	Name string `json:"name,omitempty"`
}

type Item struct {
	ID   string `json:"ID,omitempty"`
	Name string `json:"name,omitempty"`
}

type Brand struct {
	ID            string         `json:"ID,omitempty"`
	Name          string         `json:"name,omitempty"`
	MeasuringUnit *MeasuringUnit `json:"measuringUnit,omitempty"`
	Item          *Item          `json:"item,omitempty"`
}

type Store struct {
	ID   string `json:"ID,omitempty"`
	Name string `json:"name,omitempty"`
}

type StoreBranch struct {
	ID    string `json:"ID,omitempty"`
	Name  string `json:"name,omitempty"`
	Store *Store `json:"store,omitempty"`
}

/**
 * @apiDefine PriceObject200
 */
type Price struct {
	ID            string       `json:"ID,omitempty"`
	Value         float32      `json:"value,omitempty"`
	Currency      string       `json:"currency,omitempty"`
	Brand         *Brand       `json:"brand,omitempty"`
	AtStoreBranch *StoreBranch `json:"atStoreBranch,omitempty"`
}

/**
 * @apiDefine ShoppingListItem200
 * @apiSuccess (200 JSON Response Body) {String} ID
 *		Unique ID of the ShoppingListItem
 * @apiSuccess (200 JSON Response Body) {Int} quantity
 *		Number of items to get.
 * @apiSuccess (200 JSON Response Body) {Boolean} inList
 *		True if item is in list, false otherwise.
 * @apiSuccess (200 JSON Response Body) {Boolean} inCart
 *		True if item is in cart, false otherwise.
 * @apiSuccess (200 JSON Response Body) {Object} shoppingList
 *		The shopping list to which price point belongs, see "201 created JSON
 *		Response Body" of
 *		<a href="#api-Service-InsertShoppingList">New Shopping List</a>
 *		for details.
 * @apiSuccess (200 JSON Response Body) {Object} price
 *		Price details of the shopping list item.
 * @apiSuccess (200 JSON Response Body) {String} price.ID
 *		Unique ID of the Price.
 * @apiSuccess (200 JSON Response Body) {Float} price.value
 *		The price point of the item e.g. 200.
 * @apiSuccess (200 JSON Response Body) {String} price.currency
 * 		Active ISO 4217 code denoting currency of value field e.g. KES.
 * @apiSuccess (200 JSON Response Body) {Object} price.brand
 *		The brand for which provided price point applies e.g. brand.name=Colgate.
 * @apiSuccess (200 JSON Response Body) {String} price.brand.ID
 *		Unique ID of the Brand.
 * @apiSuccess (200 JSON Response Body) {String} price.brand.name
 *		Name of the Brand.
 * @apiSuccess (200 JSON Response Body) {Object} price.brand.measuringUnit
 *		The unit measurement to which the brand can be priced.
 * @apiSuccess (200 JSON Response Body) {String} price.brand.measuringUnit.ID
 *		Unique ID of the measuring unit.
 * @apiSuccess (200 JSON Response Body) {String} price.brand.measuringUnit.name
 *		Unique name of the measuring unit.
 * @apiSuccess (200 JSON Response Body) {Object} price.brand.item
 *		The item to which the brand is derived e.g. item.name=Toothpaste.
 * @apiSuccess (200 JSON Response Body) {String} price.item.ID
 *		Unique ID of the item.
 * @apiSuccess (200 JSON Response Body) {String} price.brand.item.name
 *		Unique name of the item.
 * @apiSuccess (200 JSON Response Body) {Object} [price.atStoreBranch]
 *		The store branch for which provided price point applies. This is only
 * 		available if the price has ever been checked out at a store.
 * @apiSuccess (200 JSON Response Body) {String} price.atStoreBranch.ID
 *		Unique ID of the Store Branch.
 * @apiSuccess (200 JSON Response Body) {String} price.atStoreBranch.name
 *		Unique Name of the Store Branch.
 * @apiSuccess (200 JSON Response Body) {Object} price.atStoreBranch.store
 *		The store to which the store branch belongs.
 * @apiSuccess (200 JSON Response Body) {String} price.atStoreBranch.store.ID
 *		Unique ID of the Store.
 * @apiSuccess (200 JSON Response Body) {String} price.atStoreBranch.store.name
 *		Unique Name of the Store.
 */
type ShoppingListItem struct {
	ID           string        `json:"ID,omitempty"`
	Quantity     int           `json:"quantity,omitempty"`
	InList       bool          `json:"inList"`
	InCart       bool          `json:"inCart"`
	ShoppingList *ShoppingList `json:"shoppingList,omitempty"`
	Price        *Price        `json:"price,omitempty"`
}
