package parcel_service

import (
	"fmt"
	"time"

	"github.com/Yandex-Practicum/go-db-sql-final/internal/constants"
	"github.com/Yandex-Practicum/go-db-sql-final/internal/models"
	"github.com/Yandex-Practicum/go-db-sql-final/internal/parcel_store"
)

// создаем структурный тип ParcelService
type ParcelService struct {
	store parcel_store.ParcelStore // поле store содержит структуру типа ParcelStore
}

// Функция NewParcelService возвращает новый экземпляр типа ParcelService
// Параметры
// store - экземпляр типа ParcelStore
func NewParcelService(store parcel_store.ParcelStore) ParcelService {
	return ParcelService{store: store}
}

// Метод Register типа ParcelService
// возвращает экземпляр типа Parcel и ошибку,
// а также выводит в консоль сообщение о создании новой посылки
// Параметры
// client - идентификатор клиента, целое число
// address - адрес посылки, строка
func (s ParcelService) Register(client int, address string) (models.Parcel, error) {
	// создаем новый экземпляр типа Parcel
	parcel := models.Parcel{
		Client:    client,                                // значение поля Client устанавливаем равным параметру client
		Status:    constants.ParcelStatusRegistered,      // для всех новых посылок устанавливаем статус "посылка зарегистрирована"
		Address:   address,                               // значение поля Address устанавливаем равным параметру address
		CreatedAt: time.Now().UTC().Format(time.RFC3339), // для заполнения поля CreatedAt получаем актуальное время
	}

	// получаем id новой посылки после добавления ее в базу данных
	id, err := s.store.Add(parcel)
	if err != nil {
		return parcel, err // в случае, если ошибка не равна nil, возвращаем экземпляр посылки и ошибку
	}

	//  заполняем поле Number у посылки parcel значением переменной id
	parcel.Number = id

	fmt.Printf("Новая посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s\n",
		parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt)

	return parcel, nil
}

// Метод PrintClientParcels типа ParcelService
// выводит в консоль все посылки интересующего клиента
// возвращает ошибку (по умолчанию - nil)
// Параметры
// client - идентификатор интересующего клиента (целое число)
func (s ParcelService) PrintClientParcels(client int) error {
	// получаем все посылки интересующего клиента
	parcels, err := s.store.GetByClient(client)
	if err != nil {
		return err
	}

	// выводим посылки интересующего клиента в консоль
	fmt.Printf("Посылки клиента %d:\n", client)
	for _, parcel := range parcels {
		fmt.Printf("Посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s, статус %s\n",
			parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt, parcel.Status)
	}
	fmt.Println()

	return nil
}

// Метод NextStatus типа ParcelService
// устанавливает для посылки с заданным номером следующий по порядку статус,
// выводит в консоль сообщение об обновлении статуса посылки,
// возвращает ошибку (по умолчанию nil)
// Параметры
// number - номер интересующей посылки
func (s ParcelService) NextStatus(number int) error {
	// получаем посылку из БД
	parcel, err := s.store.Get(number)
	if err != nil {
		return err
	}

	// создаем локальную переменную nextStatus, куда мы будем сохранять следующий статус посылки
	var nextStatus string

	switch parcel.Status {
	case constants.ParcelStatusRegistered:
		nextStatus = constants.ParcelStatusSent
	case constants.ParcelStatusSent:
		nextStatus = constants.ParcelStatusDelivered
	case constants.ParcelStatusDelivered: // если у посылки уже статус "Доставлена" - завершаем выполнение функции и возвращаем nil
		return nil
	}

	// выводим сообщение об обновлении статуса посылки
	fmt.Printf("У посылки № %d новый статус: %s\n", number, nextStatus)

	// возвращаем результат вызова метода s.store.SetStatus, при помощи которого обновляем статус заказа
	return s.store.SetStatus(number, nextStatus)
}

// Метод ChangeAddress типа ParcelService
// изменяет адрес доставки посылки,
// возвращает ошибку
// Параметры
// number - номер посылки, у которой необходимо изменить адрес
// address - новый адрес
func (s ParcelService) ChangeAddress(number int, address string) error {
	return s.store.SetAddress(number, address) // вызываем метод s.store.SetAddress для установки нового адреса
}

// Метод Delete типа ParcelService
// удаляет посылку с заданным номером
// возвращает ошибку
// Параметры
// number - номер посылки, которую необходимо удалить
func (s ParcelService) Delete(number int) error {
	return s.store.Delete(number)
}
