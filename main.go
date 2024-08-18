package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Sale struct {
	Product int
	Volume  int
	Date    string
}

// String реализует метод интерфейса fmt.Stringer для Sale, возвращает строковое представление объекта Sale.
// Теперь, если передать объект Sale в fmt.Println(), то выведется строка, которую вернёт эта функция.
func (s Sale) String() string {
	return fmt.Sprintf("Product: %d Volume: %d Date:%s", s.Product, s.Volume, s.Date)
}

func selectSales(client int) ([]Sale, error) {
	// подключаемся к БД
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		return nil, err
	}
	defer db.Close() // выполняем отложенное закрытие соединения с БД

	// получаем количество покупок, т.е. количество строк, у клиента с интересующим id
	var rowsCount int
	err = db.QueryRow("SELECT COUNT(*) FROM sales WHERE client = :client", sql.Named("client", client)).Scan(&rowsCount)
	if err != nil {
		return nil, err
	}
	// зная количество строк, мы можем создать слайс с известной capacity (равной rowsCount) для оптимизации выделения памяти
	var actualSales = make([]Sale, 0, rowsCount)

	// выполняем SELECT-запрос по получению информации об интересующем клиенте
	rows, err := db.Query("SELECT product, volume, date FROM sales WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// заполненяем массив в переменной actualSales объектами Sale, в которых будут данные из таблицы
	for rows.Next() {
		sale := Sale{}

		err := rows.Scan(&sale.Product, &sale.Volume, &sale.Date)
		if err != nil {
			return nil, err
		}
		actualSales = append(actualSales, sale)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actualSales, nil
}

func main() {
	client := 208

	sales, err := selectSales(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, sale := range sales {
		fmt.Println(sale)
	}
}
