package shopping

type ShoppingList struct {
	ID          string
	UserID      string
	Name        string
	Mode        string
	Created     string
	LastUpdated string
}

type MeasuringUnit struct {
	ID   string
	Name string
}

type Item struct {
	ID   string
	Name string
}

type Brand struct {
	ID            string
	Name          string
	MeasuringUnit MeasuringUnit
	Item          Item
}

type Store struct {
	ID   string
	Name string
}

type StoreBranch struct {
	ID    string
	Name  string
	Store Store
}

type Price struct {
	ID            string
	Value         float32
	Currency      string
	Brand         Brand
	AtStoreBranch StoreBranch
}

type ShoppingListItem struct {
	ID           string
	Quantity     int
	InList       bool
	InCart       bool
	ShoppingList ShoppingList
	Price        Price
}
