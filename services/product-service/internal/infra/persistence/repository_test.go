package persistence

import (
	"context"
	"log"
	"testing"

	"github.com/boretsotets/e-commerce-platform/product-service/internal/domain/models"
	"github.com/boretsotets/e-commerce-platform/product-service/internal/infra/db"
)

func TestGetByID(t *testing.T) {
	ctx := context.Background()

	db, err := db.NewTestPostgres()
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}

	repo := NewProductRepo(db)

	// 1. тест создания продукта
	pencil, err := repo.InsertNewProduct(ctx, "Карандаш", 10.5, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pencil.Name == "Карандаш" && pencil.Price == 10.5 && pencil.Stock == 100 {
	} else {
		t.Errorf("inserted product does not match input, Nae is %s, price is %v, stock is %v, id is %v", pencil.Name, pencil.Price, pencil.Stock, pencil.ID)
	}

	// 2. тест получения продукта по ID
	pencil, err = repo.GetById(ctx, pencil.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3. тест обновления остатков
	newStock := int32(95)
	updatedStock, err := repo.UpdateStock(ctx, pencil.ID, newStock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if newStock != updatedStock {
		t.Errorf("stocks doesn't match")
	}

	// 4. тест проверки существования продукта по имени
	name := "Карандаш"
	exists, err := repo.CheckProductExsistence(ctx, name)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Errorf("existing product was not found")
	}

	// 5. Проверка группового изменения остатков
	pen, _ := repo.InsertNewProduct(ctx, "Ручка", 50.5, 70)
	notebook, _ := repo.InsertNewProduct(ctx, "Блокнот", 150, 25)
	eraser, _ := repo.InsertNewProduct(ctx, "Ластик", 5, 300)

	input := []*models.StockChangeItem{
		{
			ProductID: pencil.ID,
			Delta:     10,
		},
		{
			ProductID: pen.ID,
			Delta:     -15,
		},
		{
			ProductID: notebook.ID,
			Delta:     20,
		},
		{
			ProductID: eraser.ID,
			Delta:     -50,
		},
	}

	err = repo.BatchChangeStock(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pencil, _ = repo.GetById(ctx, pencil.ID)
	if pencil.Stock != 105 {
		t.Errorf("stock for %v didn't update correctly, expected: %d fact: %d", pencil.Name, 105, pencil.Stock)
	}
	pen, _ = repo.GetById(ctx, pen.ID)
	if pen.Stock != 55 {
		t.Errorf("stock for %v didn't update correctly, expected: %d fact: %d", pen.Name, 55, pen.Stock)
	}
	notebook, _ = repo.GetById(ctx, notebook.ID)
	if notebook.Stock != 45 {
		t.Errorf("stock for %v didn't update correctly, expected: %d fact: %d", notebook.Name, 45, notebook.Stock)
	}
	eraser, _ = repo.GetById(ctx, eraser.ID)
	if eraser.Stock != 250 {
		t.Errorf("stock for %v didn't update correctly, expected: %d fact: %d", eraser.Name, 250, eraser.Stock)
	}

	// 5. Проверка получения списка продуктов
	list, err := repo.GetList(ctx, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}

	list0 := models.Product{
		ID:    pen.ID,
		Name:  "Ручка",
		Price: 50.5,
		Stock: 55,
	}
	list1 := models.Product{
		ID:    notebook.ID,
		Name:  "Блокнот",
		Price: 150,
		Stock: 45,
	}
	idmatch := list0.ID == list[0].ID && list1.ID == list[1].ID
	namematch := list0.Name == list[0].Name && list1.Name == list[1].Name
	pricematch := int64(list0.Price) == int64(list[0].Price) && int64(list1.Price) == int64(list[1].Price)
	stockmatch := int64(list0.Stock) == int64(list[0].Stock) && int64(list1.Stock) == int64(list[1].Stock)

	if !(idmatch && namematch && pricematch && stockmatch) {
		t.Errorf("got incorrect list values: id %d and id %d", list[0].ID, list[1].ID)
		t.Errorf("got incorrect list values: name %s and name %s", list[0].Name, list[1].Name)
		t.Errorf("got incorrect list values: price %f and price %f", list[0].Price, list[1].Price)
		t.Errorf("got incorrect list values: stock %d and stock %d", list[0].Stock, list[1].Stock)
	}

	// 6. проверка удаления продукта
	err = repo.DeleteProduct(ctx, pencil.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	exists, _ = repo.CheckProductExsistence(ctx, "Карандаш")
	if exists {
		t.Errorf("still exists by name")
	}
	_, err = repo.GetById(ctx, pencil.ID)
	if err == nil {
		t.Errorf("no error on get by id somewhy")
	}

	// очистка тестовой таблицы
	var p models.Product
	repo.Repo.Where("1=1").Delete(&p)
}
