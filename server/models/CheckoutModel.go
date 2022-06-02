package models

import (
	"database/sql"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// GET ALL PRODUCTS IN A CUSTOMER'S CART
// get customer_id from Customer_Cart
// 	if Customer_Cart not found (0), then create one
// 	if present, continue
// then get all data from cart_item table with the customer_id
func Checkout(customer_id string, logistic string) (result *sql.Rows, err error) {
	db, err := sql.Open("mysql", os.Getenv("DB_WRITER_INSTANCE"))
	if err != nil {
		panic(err.Error())
	}

	// random tracking number
	rand.Seed(time.Now().UnixNano())
	min, max := 1000000000, 9999999999
	tracking_number := strconv.Itoa(rand.Intn(max-min+1) + min)

	// query weight - temporary
	weight := strconv.Itoa(1200)

	query := "INSERT INTO Orders (Customer_Id, Total, Weight_Gram, Order_Date, Logistic, Tracking_Number) SELECT Customers.Customer_Id, Customer_Cart.Total, " + weight + " AS Weight_Gram, NOW() AS Order_Date, '" + logistic + "' AS Logistic, '" + tracking_number + "' AS Tracking_Number FROM Customers, Customer_Cart WHERE Customers.Customer_Id=" + customer_id + " AND Customer_Cart.Customer_Id=" + customer_id + ";"

	result, err = db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	// then insert into Order_Item FROM Cart_Item
	query = "INSERT INTO Order_Item (Order_Id, Product_Id, Quantity) SELECT (SELECT Order_Id FROM Orders WHERE Order_Date=(SELECT MAX(Order_Date) FROM Orders)), Cart_Item.Product_Id, Cart_Item.Quantity FROM Cart_Item;"
	result, err = db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	// then delete Cart_Item WHERE Customer_Id equals to...
	query = "DELETE FROM Cart_Item WHERE Customer_Id=" + customer_id
	result, err = db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	// then delete Customer_Cart WHERE Customer_Id equals to...
	query = "DELETE FROM Customer_Cart WHERE Customer_Id=" + customer_id
	result, err = db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	return result, err
}
