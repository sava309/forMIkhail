package main

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

type Product struct {
	gorm.Model
	Title  string
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// перенос схемы
	db.AutoMigrate(&Product{})

	// вставка содержимого
	db.Create(&Product{Title: "New mobile phone", Code: "D42", Price: 1000})
	db.Create(&Product{Title: "New Computer", Code: "D43", Price: 3500})

	// чтение содержимого
	var product Product
	db.First(&product, 1) // найти продукт с первичным ключём
	db.First(&product, "code = ?", "D42") // найти товар с кодом D42

	// обновить одно поле
	db.Model(&product).Update("Price", 2000)

	// обновление нескольких полей
	db.Model(&product).Updates(Product{Price: 2000, Code: "F42"}) // нулевые поля
	db.Model(&product).Updates(map[string]interface{}{"Price": 2000, "Code": "F42"})

	// удаление продукта:
	db.Delete(&product, 1)
}