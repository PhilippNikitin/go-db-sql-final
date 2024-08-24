package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	serv "github.com/Yandex-Practicum/go-db-sql-final/internal/parcel/service"
	"github.com/Yandex-Practicum/go-db-sql-final/internal/parcel/store"
)

func main() {
	// подключаемся к БД
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		// если возникла ошибка при подключении к БД, выводим ее в консоль и завершаем программу
		fmt.Printf("Возникла ошибка при подключении к базе данных: %v", err)
		return
	}

	defer db.Close()

	// создаем объект ParcelStore функцией NewParcelStore
	store := store.NewParcelStore(db)
	service := serv.NewParcelService(store)

	// регистрация посылки
	client := 1
	address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
	p, err := service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// изменение адреса
	newAddress := "Саратов, д. Верхние Зори, ул. Козлова, д. 25"
	err = service.ChangeAddress(p.Number, newAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	// изменение статуса
	err = service.NextStatus(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// попытка удаления отправленной посылки
	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
	}

	// вывод посылок клиента
	// предыдущая посылка не должна удалиться, т.к. её статус НЕ «зарегистрирована»
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// регистрация новой посылки
	p, err = service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// удаление новой посылки
	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	// здесь не должно быть последней посылки, т.к. она должна была успешно удалиться
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}
}
