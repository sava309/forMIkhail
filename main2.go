package main

import (
 "database/sql"
 "fmt"
 "log"

 _ "github.com/mattn/go-sqlite3"
)

// Заказ
type Order struct {
 ID         int
 Status     string
 OrderItems []OrderItem
}

// Позиция заказа
type OrderItem struct {
 ID        int
 OrderID   int
 ItemID    int
 Quantity  int
 Item      Item
}

// Продукты
type Item struct {
 ID       int
 ItemName string
 Amount   float32
}

var (
 items = []Item{
  {ItemName: "Go Mug", Amount: 12.49},
  {ItemName: "Go Keychain", Amount: 6.95},
  {ItemName: "Go Tshirt", Amount: 17.99},
 }
)

func main() {
 db, err := sql.Open("sqlite3", "mydatabase.db")
 if err != nil {
  log.Fatal(err)
 }
 defer db.Close()

 // Создание таблицы "Заказ"
 _, err = db.Exec(`
  CREATE TABLE IF NOT EXISTS Заказ (
   ID INTEGER PRIMARY KEY,
   Status TEXT
  )
 `)
 if err != nil {
  log.Fatal(err)
 }

 // Создание таблицы "Позиция заказа"
 _, err = db.Exec(`
  CREATE TABLE IF NOT EXISTS Позиция_заказа (
   ID INTEGER PRIMARY KEY,
   OrderID INTEGER,
   ItemID INTEGER,
   Quantity INTEGER,
   FOREIGN KEY (OrderID) REFERENCES Заказ(ID)
  )
 `)
 if err != nil {
  log.Fatal(err)
 }

 // Создание таблицы "Продукты"
 _, err = db.Exec(`
  CREATE TABLE IF NOT EXISTS Продукты (
   ID INTEGER PRIMARY KEY,
   ItemName TEXT,
   Amount REAL
  )
 `)
 if err != nil {
  log.Fatal(err)
 }

 // Вставка данных в таблицу "Продукты"
 for _, item := range items {
  _, err = db.Exec("INSERT INTO Продукты (ItemName, Amount) VALUES (?, ?)",
   item.ItemName, item.Amount)
  if err != nil {
   log.Fatal(err)
  }
 }

 // Вставка данных в таблицу "Заказ"
 order := Order{Status: "pending"}
 _, err = db.Exec("INSERT INTO Заказ (Status) VALUES (?)", order.Status)
 if err != nil {
  log.Fatal(err)
 }

 // Вставка данных в таблицу "Позиция заказа"
 item1 := OrderItem{OrderID: 1, ItemID: 1, Quantity: 1}
 item2 := OrderItem{OrderID: 1, ItemID: 2, Quantity: 4}
 _, err = db.Exec("INSERT INTO Позиция_заказа (OrderID, ItemID, Quantity) VALUES (?, ?, ?)",
  item1.OrderID, item1.ItemID, item1.Quantity)
 if err != nil {
  log.Fatal(err)
 }
 _, err = db.Exec("INSERT INTO Позиция_заказа (OrderID, ItemID, Quantity) VALUES (?, ?, ?)",
  item2.OrderID, item2.ItemID, item2.Quantity)
 if err != nil {
  log.Fatal(err)
 }

 // Запрос данных о заказе с использованием JOIN
 rows, err := db.Query(`
  SELECT Заказ.ID, Заказ.Status, Позиция_заказа.ID, Позиция_заказа.OrderID, Позиция_заказа.ItemID,
   Позиция_заказа.Quantity, Продукты.ItemName, Продукты.Amount
  FROM Заказ
  JOIN Позиция_заказа ON Позиция_заказа.OrderID = Заказ.ID
  JOIN Продукты ON Продукты.ID = Позиция_заказа.ItemID
  WHERE Заказ.ID = ? AND Заказ.Status = ?
 `, order.ID, order.Status)
 if err != nil {
  log.Fatal(err)
 }
 defer rows.Close()

 // Загрузка данных в структуру Order
 newOrder := &Order{}
 newOrder.OrderItems = make([]OrderItem, 0)

 for rows.Next() {
  orderItem := OrderItem{}
  item := Item{}
  err = rows.Scan(&newOrder.ID, &newOrder.Status, &orderItem.ID, &orderItem.OrderID,
   &orderItem.ItemID, &orderItem.Quantity, &item.ItemName, &item.Amount)
  if err != nil {
   log.Fatal(err)
  }
  orderItem.Item = item
  newOrder.OrderItems = append(newOrder.OrderItems, orderItem)
 }

 fmt.Printf("%#v\n", newOrder)
}
