package constants

// в пакете хранятся константы - возможные значения, которые может принимать поле Status структуры Parcel

const (
	// объявляем константы с возможными статусами посылок
	ParcelStatusRegistered = "registered" // посылка зарегистрирована
	ParcelStatusSent       = "sent"       // посылка отправлена
	ParcelStatusDelivered  = "delivered"  // посылка доставлена

	// объявляем константу с адресом локальной БД
	PathToLocalDB = "/home/nikitin/Dev/Yandex Practicum/sprint8/go-db-sql-final/tracker.db"
)
