package shop

type ShopItem struct {
	ItemKey     string `json:"item_key"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

type InventoryItem struct {
	UserID      string `json:"user_id"`
	ItemKey     string `json:"item_key"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}
