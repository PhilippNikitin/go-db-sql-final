package store

import (
	"database/sql"
	"fmt"

	"github.com/Yandex-Practicum/go-db-sql-final/internal/constants"
	"github.com/Yandex-Practicum/go-db-sql-final/internal/models"
	"github.com/Yandex-Practicum/go-db-sql-final/internal/parcel/errors"
)

// определяем структурный тип ParcelStore для работы с БД
type ParcelStore struct {
	db *sql.DB // единственное поле db - указатель на БД
}

// функция NewParcelStore для создания нового экземпляра ParcelStore
// Параметры
// db - указатель на БД
// возвращает новый экземпляр ParcelStore
func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Метод Add типа ParcelStore добавляет
// в таблицу parcel в БД запись для новой посылки
// Параметры
// p - экземпляр типа Parcel
// поля данной переменной будут использоваться
// для заполнения соответствующих атрибутов в таблице parcel
// возвращает идентификатор последней добавленной записи
func (s ParcelStore) Add(p models.Parcel) (int, error) {
	res, err := s.db.Exec(`INSERT INTO parcel (client, status, address, created_at)
						 VALUES (:client, :status, :address, :created_at)`,
		sql.Named("client", p.Client), sql.Named("status", p.Status),
		sql.Named("address", p.Address), sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	// получаем id последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// возвращаем id последней добавленной записи
	return int(id), nil
}

// Метод Get типа ParcelStore
// получает данные о посылке из БД по идентификатору посылки
// Параметры
// number - идентификатор посылки
// возвращает экземпляр типа Parcel
// и ошибку, если она возникла в ходе выполнения функции
func (s ParcelStore) Get(number int) (models.Parcel, error) {
	// из таблицы возвращается только одна строка
	row := s.db.QueryRow(`SELECT number, client, status, address, created_at
						  FROM parcel
						  WHERE number = :number`,
		sql.Named("number", number))

	// заполняем объект Parcel полученными данными
	p := models.Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

// Метод GetByClient типа ParcelStore
// применяется для получения всех посылок интересующего клиента из БД
// Параметры
// client - идентификатор клиента, посылки которого мы хотим получить
// возвращает слайс из структур типа Parcel
// и ошибку, если она возникла в ходе выполнения функции
func (s ParcelStore) GetByClient(client int) ([]models.Parcel, error) {
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query(`SELECT number, client, status, address, created_at
							 FROM parcel
							 WHERE client = :client`, sql.Named("client", client))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// cоздаем слайс для сохранения всех посылок клиента
	var res = make([]models.Parcel, 0)

	for rows.Next() {
		p := models.Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, err
		}
		res = append(res, p)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}

// метод SetStatus типа ParcelStore
// позволяет изменить статус у заданной посылки
// Параметры
// number - номер посылки
// status - новый статус посылки
// возвращает ошибку
func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

// метод SetAddress типа ParcelStore
// изменяет адрес у посылки с заданным идентификатором
// Изменение адреса возможно, только если
// статус посылки равен `зарегистрирована`
// Параметры
// number - идентификатор посылки
// address - новый адрес
// возвращает ошибку
func (s ParcelStore) SetAddress(number int, address string) error {
	res, err := s.db.Exec(`UPDATE parcel
						SET address = :address
						WHERE number = :number AND
							  status = :registered`,
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("registered", constants.ParcelStatusRegistered))

	if err != nil {
		return err
	}

	// проверяем количество измененных строк при помощи переменной res
	// равенство 0 возможно в двух случаях:
	// в функцию был передан несуществующий номер посылки
	// или ее статус не равен `зарегистрирована`
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// информируем о неудачной операции
	if rowsAffected == 0 {
		return errors.ErrUnsuccessful
	}

	return nil
}

// Метод Delete типа ParcelStore
// удаляет посылку из БД (таблицы parcel)
// удалить посылку можно, только если ее статус
// равен `зарегистрирована`
// Параметры
// number - номер посылки, которую требуется удалить
func (s ParcelStore) Delete(number int) error {
	res, err := s.db.Exec(`DELETE FROM parcel
						WHERE number = :number AND
						status = :registered`,
		sql.Named("number", number),
		sql.Named("registered", constants.ParcelStatusRegistered))
	if err != nil {
		fmt.Println(err)
		return err
	}

	// проверяем, что посылка была удалена, при помощи res.RowsAffected()
	// неуспешная операция возможна в двух случаях:
	// в функцию был передан несуществующий номер посылки
	// или ее статус не равен `зарегистрирована`
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// информируем о неудачной операции
	if rowsAffected == 0 {
		return errors.ErrUnsuccessful
	}

	return nil
}
