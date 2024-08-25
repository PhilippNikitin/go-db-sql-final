package store

import (
	// импортируем пакеты standard library
	"database/sql"
	"math/rand"
	"testing"
	"time"

	// импортируем пакеты third-party
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	// импортируем локальные пакеты проекта
	"github.com/Yandex-Practicum/go-db-sql-final/internal/constants"
	"github.com/Yandex-Practicum/go-db-sql-final/internal/models"
	"github.com/Yandex-Practicum/go-db-sql-final/internal/parcel/errors"
)

var (
	// randSource источник псевдо случайных чисел.
	// в качестве seed используется текущее время
	// в unix формате
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() models.Parcel {
	return models.Parcel{
		Client:    1000,
		Status:    constants.ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// подключаемся к БД
	db, err := sql.Open("sqlite", "../../../tracker.db")
	require.NoError(t, err)

	defer db.Close()

	// получаем экземпляр ParcelStore
	store := NewParcelStore(db)
	// получаем тестовый экземпляр посылки
	parcel := getTestParcel()

	// add
	// добавляем новую посылку в БД, проверяем отсутствие ошибки и наличие идентификатора
	num, err := store.Add(parcel)
	require.NoError(t, err)  // убеждаемся в отсутствии ошибки
	require.NotEmpty(t, num) // убеждаемся в наличии идентификатора

	// get
	// получаем только что добавленную посылку, проверяем отсутствие ошибки
	storedParcel, err := store.Get(num)
	require.NoError(t, err) // убеждаемся в отсутствии ошибки

	// проверяем, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	// структуры равны, если значения всех их полей равны (и типы полей допускают сравнения)

	// устанавливаем значение поле Number у тестового экземпляра равное переменной num
	parcel.Number = num
	// сравниваем две структуры
	assert.Equal(t, parcel, storedParcel)

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	err = store.Delete(num)
	require.NoError(t, err) // убеждаемся в отсутствии ошибки

	// проверьте, что посылку больше нельзя получить из БД
	_, err = store.Get(num)
	require.Error(t, err)                 // если мы запрашиваем из БД несуществующую посылку, должна вернуться ошибка
	assert.ErrorIs(t, err, sql.ErrNoRows) // проверяем, что по крайней мере одна ошибка из соответствующей цепи ошибок err равна sql.ErrNoRows

}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// подключаемся к БД
	db, err := sql.Open("sqlite", "../../../tracker.db")
	require.NoError(t, err)

	defer db.Close()

	// получаем экземпляр ParcelStore
	store := NewParcelStore(db)
	// получаем тестовый экземпляр посылки
	parcel := getTestParcel()

	// add
	// добавляем новую посылку в БД, проверяем отсутствие ошибки и наличие идентификатора
	num, err := store.Add(parcel)
	require.NoError(t, err)  // убеждаемся в отсутствии ошибки
	require.NotEmpty(t, num) // убеждаемся в наличии идентификатора

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(num, newAddress)
	require.NoError(t, err) // убеждаемся в отсутствии ошибки

	// check
	// получаем добавленную посылку, проверяем, что адрес обновился
	storedParcel, err := store.Get(num)
	require.NoError(t, err)                           // проверяем, что при получении посылки не возникло ошибки
	assert.Equal(t, newAddress, storedParcel.Address) // проверяем, что адрес посылки изменился на предполагаемый
}

// TestSetStatus проверяет обновление статуса, а также
// то, что мы не можем изменить адрес посылки или удалить ее,
// если статус посылки не равен `зарегистрирована`
func TestSetStatus(t *testing.T) {
	// подключаемся к БД
	db, err := sql.Open("sqlite", "../../../tracker.db")
	require.NoError(t, err)

	defer db.Close()

	// получаем экземпляр ParcelStore
	store := NewParcelStore(db)
	// получаем тестовый экземпляр посылки
	parcel := getTestParcel()

	// add
	// добавляем новую посылку в БД, проверяем отсутствие ошибки и наличие идентификатора
	num, err := store.Add(parcel)
	require.NoError(t, err)  // убеждаемся в отсутствии ошибки
	require.NotEmpty(t, num) // убеждаемся в наличии идентификатора

	// set status
	// обновляем статус, проверяем отсутствие ошибки
	err = store.SetStatus(num, constants.ParcelStatusSent)
	require.NoError(t, err) // убеждаемся в отсутствии ошибки

	// check
	// получаем добавленную посылку и убеждаемся, что статус обновился
	storedParcel, err := store.Get(num)
	require.NoError(t, err)                                           // убеждаемся в отсутствии ошибки
	require.Equal(t, constants.ParcelStatusSent, storedParcel.Status) // проверяем, что статус обновился

	// проверяем, что нельзя изменить адрес, если статус посылки не равен `зарегистрирована`
	newAddress := "new test address"
	oldAddress := "test"
	err = store.SetAddress(num, newAddress)
	// убеждаемся, что вернулась ошибка
	// и она равна ErrUnsuccessful
	assert.ErrorIs(t, err, errors.ErrUnsuccessful)
	// проверяем, что адрес посылки не изменился
	storedParcel, err = store.Get(num)
	require.NoError(t, err)                            // убеждаемся в отсутствии ошибки
	require.Equal(t, oldAddress, storedParcel.Address) // убеждаемся, что адрес не изменился

	// проверяем, что мы не можем удалить посылку, если ее статус не равен `зарегистрирована`
	err = store.Delete(num)
	// убеждаемся, что вернулась ошибка
	// и она равна ErrUnsuccessful
	assert.ErrorIs(t, err, errors.ErrUnsuccessful)
	// проверяем, что мы по-прежнему можем получить посылку из БД
	testParcel := parcel
	// устанавливаем значение поля Number посылки testParcel
	// равное значению переменной num
	testParcel.Number = num
	// устанавиливаем статус testParcel равным ParcelStatusSent
	testParcel.Status = constants.ParcelStatusSent

	storedParcel, err = store.Get(num)
	require.NoError(t, err)                   // убеждаемся в отсутствии ошибки
	assert.Equal(t, testParcel, storedParcel) // проверяем, что поля посылки не изменились

}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// подключаемся к БД
	db, err := sql.Open("sqlite", "../../../tracker.db")
	require.NoError(t, err)

	defer db.Close()

	// получаем экземпляр ParcelStore
	store := NewParcelStore(db)

	parcels := []models.Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)

	// add
	for i := 0; i < len(parcels); i++ {
		// задаём всем посылкам один и тот же идентификатор клиента
		parcels[i].Client = client

		num, err := store.Add(parcels[i]) // добавляем новую посылку в БД, убеждаемся в отсутствии ошибки и наличии идентификатора
		require.NoError(t, err)           // убеждаемся в отсутствии ошибки
		require.NotEmpty(t, num)          // убеждаемся в наличии идентификатора

		// обновляем идентификатор у добавленной посылки
		parcels[i].Number = num
	}

	// get by client
	storedParcels, err := store.GetByClient(client) // получаем список посылок по идентификатору клиента, сохранённому в переменной client
	require.NoError(t, err)                         // проверяем отсутствие ошибки

	// проверяем, что количество элементов в слайсах parcels и storedParcels равно, и каждому элементу из одного слайса
	// есть соответсвующий равный элемент из другого слайса, то есть все посылки из storedParcels есть в parcels
	// и значения полей полученных посылок заполнены верно
	assert.ElementsMatch(t, parcels, storedParcels)
}
