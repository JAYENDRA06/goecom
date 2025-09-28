package cart

import (
	"fmt"

	"github.com/JAYENDRA06/apiproject/types"
)

func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	productIDs := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductID)
		}
		productIDs[i] = item.ProductID
	}

	return productIDs, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}
	// check if all products are actually in stock
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}
	// calculate the total price
	totalPrice, err := getTotalPrice(items, productMap)
	if err != nil {
		return 0, 0, err
	}
	// reduce quantity of products in our db
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}
	// create the order
	order := types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Address: "random",
		Status:  "pending",
	}
	orderID, err := h.store.CreateOrder(order)
	if err != nil {
		return 0, 0, err
	}
	// create order items
	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}
	return orderID, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []types.CartItem, productMap map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	var itemsOutOfStock []string
	for _, item := range cartItems {
		if product, ok := productMap[item.ProductID]; ok {
			if item.Quantity > product.Quantity {
				itemsOutOfStock = append(itemsOutOfStock, productMap[item.ProductID].Name)
			}
		}
	}
	if len(itemsOutOfStock) > 0 {
		itemNames := "following items are out of stock: "
		for idx, it := range itemsOutOfStock {
			itemNames += it
			if idx < len(itemsOutOfStock)-1 {
				itemNames += ", "
			}
		}
		return fmt.Errorf("%s", itemNames)
	}
	return nil
}

func getTotalPrice(cartItems []types.CartItem, productMap map[int]types.Product) (float64, error) {
	var total float64
	for _, item := range cartItems {
		if product, ok := productMap[item.ProductID]; ok {
			total += product.Price
		} else {
			return 0, fmt.Errorf("some items in cart are removed")
		}
	}
	return total, nil
}
