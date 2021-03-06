define({ "api": [
  {
    "type": "delete",
    "url": "/items/{ID}",
    "title": "Delete Shopping List Item",
    "name": "DeleteShoppingListItem",
    "version": "0.1.0",
    "group": "Service",
    "description": "<p>Delete a shopping list Item. Note that this only deletes the top level shopping list, the price details remain intact.</p>",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          },
          {
            "group": "Header",
            "optional": false,
            "field": "Authorization",
            "description": "<p>Bearer token received from authentication micro-service in the form &quot;Bearer {token-value}&quot;.</p>"
          }
        ]
      }
    },
    "parameter": {
      "fields": {
        "URL Path Params": [
          {
            "group": "URL Path Params",
            "type": "String",
            "optional": false,
            "field": "id",
            "description": "<p>The ID of the shopping list item to delete.</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "200": [
          {
            "group": "200",
            "optional": false,
            "field": "emptyBody",
            "description": "<p>check status code for success.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service"
  },
  {
    "type": "get",
    "url": "/docs",
    "title": "Docs",
    "name": "Docs",
    "version": "0.1.0",
    "group": "Service",
    "success": {
      "fields": {
        "200": [
          {
            "group": "200",
            "type": "html",
            "optional": false,
            "field": "docs",
            "description": "<p>Docs page to be viewed on browser.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service"
  },
  {
    "type": "get",
    "url": "/shoppinglists/{ID}/items",
    "title": "Get Shopping List Items",
    "name": "GetShoppingListItems",
    "version": "0.1.0",
    "group": "Service",
    "description": "<p>Get shopping items for a shopping list.</p>",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          },
          {
            "group": "Header",
            "optional": false,
            "field": "Authorization",
            "description": "<p>Bearer token received from authentication micro-service in the form &quot;Bearer {token-value}&quot;.</p>"
          }
        ]
      }
    },
    "parameter": {
      "fields": {
        "URL Query Params": [
          {
            "group": "URL Query Params",
            "type": "Long",
            "optional": true,
            "field": "offset",
            "defaultValue": "0",
            "description": "<p>Offset index to fetch from.</p>"
          },
          {
            "group": "URL Query Params",
            "type": "Long",
            "optional": true,
            "field": "count",
            "defaultValue": "10",
            "description": "<p>Number of shopping lists to fetch.</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "200 JSON Response Body": [
          {
            "group": "200 JSON Response Body",
            "type": "Object[]",
            "optional": false,
            "field": "items",
            "description": "<p>List of ShoppingListItems. See &quot;200 JSON Response Body&quot; of <a href=\"#api-Service-UpsertShoppingListItem\">Upsert Shopping List Item</a> for details on what each item looks like.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service"
  },
  {
    "type": "get",
    "url": "/shoppinglists",
    "title": "Get Shopping Lists",
    "name": "GetShoppingLists",
    "version": "0.1.0",
    "group": "Service",
    "description": "<p>Get shopping lists for a user.</p>",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          },
          {
            "group": "Header",
            "optional": false,
            "field": "Authorization",
            "description": "<p>Bearer token received from authentication micro-service in the form &quot;Bearer {token-value}&quot;.</p>"
          }
        ]
      }
    },
    "parameter": {
      "fields": {
        "URL Query Params": [
          {
            "group": "URL Query Params",
            "type": "Long",
            "optional": true,
            "field": "offset",
            "defaultValue": "0",
            "description": "<p>Offset index to fetch from.</p>"
          },
          {
            "group": "URL Query Params",
            "type": "Long",
            "optional": true,
            "field": "count",
            "defaultValue": "10",
            "description": "<p>Number of shopping lists to fetch.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service",
    "success": {
      "fields": {
        "200 JSON Response Body": [
          {
            "group": "200 JSON Response Body",
            "type": "Object[]",
            "optional": false,
            "field": "shoppingLists",
            "description": "<p>List of shoppingLists.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "shoppingLists.ID",
            "description": "<p>Unique ID of the shopping list.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "shoppingLists.userID",
            "description": "<p>ID of the user who owns the shopping list.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "shoppingLists.name",
            "description": "<p>Unique name of the shopping list.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "allowedValues": [
              "\"PREPARATION\"",
              "\"SHOPPING\""
            ],
            "optional": false,
            "field": "shoppingLists.mode",
            "description": "<p>The current mode of the shopping list on the client apps.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "shoppingLists.created",
            "description": "<p>ISO8601 date of shopping list creation.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "shoppingLists.lastUpdated",
            "description": "<p>ISO8601 date denoting last time the list was updated.</p>"
          }
        ]
      }
    }
  },
  {
    "type": "put",
    "url": "/shoppinglists",
    "title": "New Shopping List",
    "name": "InsertShoppingList",
    "version": "0.1.0",
    "group": "Service",
    "description": "<p>insert a shopping list by name if not exists</p>",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          },
          {
            "group": "Header",
            "optional": false,
            "field": "Authorization",
            "description": "<p>Bearer token received from authentication micro-service in the form &quot;Bearer {token-value}&quot;.</p>"
          }
        ]
      }
    },
    "parameter": {
      "fields": {
        "JSON Request Body": [
          {
            "group": "JSON Request Body",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>The name of the new shopping list.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service",
    "success": {
      "fields": {
        "200 existed JSON Response Body": [
          {
            "group": "200 existed JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "ID",
            "description": "<p>Unique ID of the shopping list.</p>"
          },
          {
            "group": "200 existed JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "userID",
            "description": "<p>ID of the user who owns the shopping list.</p>"
          },
          {
            "group": "200 existed JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>Unique name of the shopping list.</p>"
          },
          {
            "group": "200 existed JSON Response Body",
            "type": "String",
            "allowedValues": [
              "\"PREPARATION\"",
              "\"SHOPPING\""
            ],
            "optional": false,
            "field": "mode",
            "description": "<p>The current mode of the shopping list on the client apps.</p>"
          },
          {
            "group": "200 existed JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "created",
            "description": "<p>ISO8601 date of shopping list creation.</p>"
          },
          {
            "group": "200 existed JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "lastUpdated",
            "description": "<p>ISO8601 date denoting last time the list was updated.</p>"
          }
        ],
        "201 created JSON Response Body": [
          {
            "group": "201 created JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "ID",
            "description": "<p>Unique ID of the shopping list.</p>"
          },
          {
            "group": "201 created JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "userID",
            "description": "<p>ID of the user who owns the shopping list.</p>"
          },
          {
            "group": "201 created JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>Unique name of the shopping list.</p>"
          },
          {
            "group": "201 created JSON Response Body",
            "type": "String",
            "allowedValues": [
              "\"PREPARATION\"",
              "\"SHOPPING\""
            ],
            "optional": false,
            "field": "mode",
            "description": "<p>The current mode of the shopping list on the client apps.</p>"
          },
          {
            "group": "201 created JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "created",
            "description": "<p>ISO8601 date of shopping list creation.</p>"
          },
          {
            "group": "201 created JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "lastUpdated",
            "description": "<p>ISO8601 date denoting last time the list was updated.</p>"
          }
        ]
      }
    }
  },
  {
    "type": "get",
    "url": "/items/search",
    "title": "Search Shopping Items",
    "name": "SearchShoppingItems",
    "version": "0.1.0",
    "group": "Service",
    "description": "<p>Search Shopping Items not necessarily belonging to a specific ShoppingList.</p>",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          },
          {
            "group": "Header",
            "optional": false,
            "field": "Authorization",
            "description": "<p>Bearer token received from authentication micro-service in the form &quot;Bearer {token-value}&quot;.</p>"
          }
        ]
      }
    },
    "parameter": {
      "fields": {
        "URL Query Params": [
          {
            "group": "URL Query Params",
            "type": "Long",
            "optional": true,
            "field": "offset",
            "defaultValue": "0",
            "description": "<p>Offset index to fetch from.</p>"
          },
          {
            "group": "URL Query Params",
            "type": "Long",
            "optional": true,
            "field": "count",
            "defaultValue": "10",
            "description": "<p>Number of shopping lists to fetch.</p>"
          },
          {
            "group": "URL Query Params",
            "type": "String",
            "optional": true,
            "field": "brandName",
            "description": "<p>If provided, filter items where brandName contains provided text.</p>"
          },
          {
            "group": "URL Query Params",
            "type": "String",
            "optional": true,
            "field": "itemName",
            "description": "<p>If provided, filter items where itemName contains provided text.</p>"
          },
          {
            "group": "URL Query Params",
            "type": "String",
            "optional": true,
            "field": "brandPrice",
            "description": "<p>If provided, filter items where brandPrice contains provided text.</p>"
          },
          {
            "group": "URL Query Params",
            "type": "String",
            "optional": true,
            "field": "measuringUnit",
            "description": "<p>If provided, filter items where measuringUnit contains provided text.</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "200 JSON Response Body": [
          {
            "group": "200 JSON Response Body",
            "type": "Object[]",
            "optional": false,
            "field": "items",
            "description": "<p>List of ShoppingListItems. See &quot;200 JSON Response Body&quot; of <a href=\"#api-Service-UpsertShoppingListItem\">Upsert Shopping List Item</a> for details on what each item looks like.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service"
  },
  {
    "type": "get",
    "url": "/status",
    "title": "Status",
    "name": "Status",
    "version": "0.1.0",
    "group": "Service",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "200": [
          {
            "group": "200",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>Micro-service name.</p>"
          },
          {
            "group": "200",
            "type": "String",
            "optional": false,
            "field": "version",
            "description": "<p>http://semver.org version.</p>"
          },
          {
            "group": "200",
            "type": "String",
            "optional": false,
            "field": "description",
            "description": "<p>Short description of the micro-service.</p>"
          },
          {
            "group": "200",
            "type": "String",
            "optional": false,
            "field": "canonicalName",
            "description": "<p>Canonical name of the micro-service.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service"
  },
  {
    "type": "put",
    "url": "/shoppinglists/{ID}",
    "title": "Update Shopping List",
    "name": "UpdateShoppingList",
    "version": "0.1.0",
    "group": "Service",
    "description": "<p>update a shopping list with {ID}.</p>",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          },
          {
            "group": "Header",
            "optional": false,
            "field": "Authorization",
            "description": "<p>Bearer token received from authentication micro-service in the form &quot;Bearer {token-value}&quot;.</p>"
          }
        ]
      }
    },
    "parameter": {
      "fields": {
        "URL Path Params": [
          {
            "group": "URL Path Params",
            "type": "String",
            "optional": false,
            "field": "id",
            "description": "<p>The ID of the shopping list.</p>"
          }
        ],
        "JSON Request Body": [
          {
            "group": "JSON Request Body",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>Unique name of the shopping list.</p>"
          },
          {
            "group": "JSON Request Body",
            "type": "String",
            "allowedValues": [
              "\"PREPARATION\"",
              "\"SHOPPING\""
            ],
            "optional": false,
            "field": "mode",
            "description": "<p>The current mode of the shopping list on the client apps.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service",
    "success": {
      "fields": {
        "200 existed JSON Response Body": [
          {
            "group": "200 existed JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "ID",
            "description": "<p>Unique ID of the shopping list.</p>"
          },
          {
            "group": "200 existed JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "userID",
            "description": "<p>ID of the user who owns the shopping list.</p>"
          },
          {
            "group": "200 existed JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>Unique name of the shopping list.</p>"
          },
          {
            "group": "200 existed JSON Response Body",
            "type": "String",
            "allowedValues": [
              "\"PREPARATION\"",
              "\"SHOPPING\""
            ],
            "optional": false,
            "field": "mode",
            "description": "<p>The current mode of the shopping list on the client apps.</p>"
          },
          {
            "group": "200 existed JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "created",
            "description": "<p>ISO8601 date of shopping list creation.</p>"
          },
          {
            "group": "200 existed JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "lastUpdated",
            "description": "<p>ISO8601 date denoting last time the list was updated.</p>"
          }
        ]
      }
    }
  },
  {
    "type": "put",
    "url": "/shoppinglists/{ID}/items",
    "title": "Upsert Shopping List Item",
    "name": "UpsertShoppingListItem",
    "version": "0.1.0",
    "group": "Service",
    "description": "<p>Update/Insert a list item's values for a shopping list. Note that all details under the Price object are shared with other users and will not be deleted during item deletion.</p>",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          },
          {
            "group": "Header",
            "optional": false,
            "field": "Authorization",
            "description": "<p>Bearer token received from authentication micro-service in the form &quot;Bearer {token-value}&quot;.</p>"
          }
        ]
      }
    },
    "parameter": {
      "fields": {
        "URL Path Params": [
          {
            "group": "URL Path Params",
            "type": "String",
            "optional": false,
            "field": "id",
            "description": "<p>The ID of the shopping list.</p>"
          }
        ],
        "JSON Request Body": [
          {
            "group": "JSON Request Body",
            "type": "String",
            "optional": false,
            "field": "itemName",
            "description": "<p>Name of the ShoppingItem e.g. Toothpaste.</p>"
          },
          {
            "group": "JSON Request Body",
            "type": "Boolean",
            "optional": true,
            "field": "inList",
            "description": "<p>True if item is in the shopping list, false otherwise.</p>"
          },
          {
            "group": "JSON Request Body",
            "type": "Boolean",
            "optional": true,
            "field": "inCart",
            "description": "<p>True if item is in the shopping cart, false otherwise. Marking this as true automatically sets inList to true.</p>"
          },
          {
            "group": "JSON Request Body",
            "type": "String",
            "optional": true,
            "field": "brandName",
            "description": "<p>Name of the Brand of the itemName e.g. Colgate.</p>"
          },
          {
            "group": "JSON Request Body",
            "type": "Int",
            "optional": true,
            "field": "quantity",
            "description": "<p>Number of items in the shopping List.</p>"
          },
          {
            "group": "JSON Request Body",
            "type": "String",
            "optional": true,
            "field": "measurementUnit",
            "description": "<p>The measurement Unit to use e.g. 250ml Tub, KG, 5Kg bag, etc.</p>"
          },
          {
            "group": "JSON Request Body",
            "type": "Float",
            "optional": true,
            "field": "unitPrice",
            "description": "<p>Price of one unit of measurement e.g. 200 if a 250ml Tub costs that.</p>"
          },
          {
            "group": "JSON Request Body",
            "type": "String",
            "optional": true,
            "field": "currency",
            "defaultValue": "KES",
            "description": "<p>Active ISO 4217 code denoting currency of the unitPrice.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service",
    "success": {
      "fields": {
        "200 JSON Response Body": [
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "ID",
            "description": "<p>Unique ID of the ShoppingListItem</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "Int",
            "optional": false,
            "field": "quantity",
            "description": "<p>Number of items to get.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "Boolean",
            "optional": false,
            "field": "inList",
            "description": "<p>True if item is in list, false otherwise.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "Boolean",
            "optional": false,
            "field": "inCart",
            "description": "<p>True if item is in cart, false otherwise.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "Object",
            "optional": false,
            "field": "shoppingList",
            "description": "<p>The shopping list to which price point belongs, see &quot;201 created JSON Response Body&quot; of <a href=\"#api-Service-InsertShoppingList\">New Shopping List</a> for details.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "Object",
            "optional": false,
            "field": "price",
            "description": "<p>Price details of the shopping list item.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "price.ID",
            "description": "<p>Unique ID of the Price.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "Float",
            "optional": false,
            "field": "price.value",
            "description": "<p>The price point of the item e.g. 200.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "price.currency",
            "description": "<p>Active ISO 4217 code denoting currency of value field e.g. KES.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "Object",
            "optional": false,
            "field": "price.brand",
            "description": "<p>The brand for which provided price point applies e.g. brand.name=Colgate.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "price.brand.ID",
            "description": "<p>Unique ID of the Brand.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "price.brand.name",
            "description": "<p>Name of the Brand.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "Object",
            "optional": false,
            "field": "price.brand.measuringUnit",
            "description": "<p>The unit measurement to which the brand can be priced.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "price.brand.measuringUnit.ID",
            "description": "<p>Unique ID of the measuring unit.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "price.brand.measuringUnit.name",
            "description": "<p>Unique name of the measuring unit.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "Object",
            "optional": false,
            "field": "price.brand.item",
            "description": "<p>The item to which the brand is derived e.g. item.name=Toothpaste.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "price.item.ID",
            "description": "<p>Unique ID of the item.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "price.brand.item.name",
            "description": "<p>Unique name of the item.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "Object",
            "optional": true,
            "field": "price.atStoreBranch",
            "description": "<p>The store branch for which provided price point applies. This is only available if the price has ever been checked out at a store.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "price.atStoreBranch.ID",
            "description": "<p>Unique ID of the Store Branch.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "price.atStoreBranch.name",
            "description": "<p>Unique Name of the Store Branch.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "Object",
            "optional": false,
            "field": "price.atStoreBranch.store",
            "description": "<p>The store to which the store branch belongs.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "price.atStoreBranch.store.ID",
            "description": "<p>Unique ID of the Store.</p>"
          },
          {
            "group": "200 JSON Response Body",
            "type": "String",
            "optional": false,
            "field": "price.atStoreBranch.store.name",
            "description": "<p>Unique Name of the Store.</p>"
          }
        ]
      }
    }
  }
] });
