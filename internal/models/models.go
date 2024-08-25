package models

// определяем структурый тип Parcel ("посылка")
type Parcel struct {
	Number    int    // номер посылки, в БД это автоинкрементное поле
	Client    int    // идентификатор клиента
	Status    string // статус посылки
	Address   string // адрес посылки
	CreatedAt string // дата и время создания посылки
}
