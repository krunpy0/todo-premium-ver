package shop

import (
	"errors"

	"github.com/krunpy0/todo-premium-ver/db"
)

var (
	ErrItemNotFound      = errors.New("shop item not found")
	ErrInsufficientCoins = errors.New("insufficient coins")
)

var shopItems = []ShopItem{
	{
		ItemKey:     "placeholder_item",
		Name:        "Placeholder Item",
		Price:       5,
		Description: "Placeholder description",
	},
}

func QueryShopItems() ([]ShopItem, error) {
	var items = []ShopItem{}
	rows, err := db.DB.Query(`SELECT * FROM shop_items`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i = ShopItem{}
		err := rows.Scan(
			&i.ItemKey,
			&i.Name,
			&i.Price,
			&i.Description,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, err

}

func QueryShopItem(itemKey string) (ShopItem, error) {
	var item = ShopItem{}
	err := db.DB.QueryRow(`SELECT * FROM shop_items WHERE item_key=$1`, itemKey).Scan(&item.ItemKey, &item.Name, &item.Price, &item.Description)
	if err != nil {
		return ShopItem{}, err
	}
	return item, err
}

func PurchaseItem(userID string, itemKey string) (InventoryItem, int, error) {
	item, err := QueryShopItem(itemKey)
	if err != nil {
		return InventoryItem{}, 0, err
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return InventoryItem{}, 0, err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS user_inventory (
			user_id TEXT NOT NULL,
			item_key TEXT NOT NULL,
			name TEXT NOT NULL,
			price INT NOT NULL,
			description TEXT NOT NULL,
			quantity INT NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (user_id, item_key)
		);`); err != nil {
		return InventoryItem{}, 0, err
	}

	var coins int
	if err = tx.QueryRow(`SELECT coins FROM users WHERE id = $1 FOR UPDATE;`, userID).Scan(&coins); err != nil {
		return InventoryItem{}, 0, err
	}
	if coins < item.Price {
		return InventoryItem{}, 0, ErrInsufficientCoins
	}

	var updatedCoins int
	if err = tx.QueryRow(`
		UPDATE users
		SET coins = coins - $1, updated_at = NOW()
		WHERE id = $2
		RETURNING coins;`, item.Price, userID).Scan(&updatedCoins); err != nil {
		return InventoryItem{}, 0, err
	}

	var quantity int
	if err = tx.QueryRow(`
		INSERT INTO user_inventory (user_id, item_key, name, price, description, quantity)
		VALUES ($1, $2, $3, $4, $5, 1)
		ON CONFLICT (user_id, item_key)
		DO UPDATE SET quantity = user_inventory.quantity + 1, updated_at = NOW()
		RETURNING quantity;`,
		userID, item.ItemKey, item.Name, item.Price, item.Description,
	).Scan(&quantity); err != nil {
		return InventoryItem{}, 0, err
	}

	if err = tx.Commit(); err != nil {
		return InventoryItem{}, 0, err
	}

	return InventoryItem{
		UserID:      userID,
		ItemKey:     item.ItemKey,
		Name:        item.Name,
		Price:       item.Price,
		Description: item.Description,
		Quantity:    quantity,
	}, updatedCoins, nil
}
