package parcel_errors

import "errors"

// статически создаем ошибку ErrUnsuitableParcelStatus,
// которая будет возникать в ситуациях, когда происходит попытка изменить адрес доставки посылки
// или удалить посылку, но статус последней не равен `зарегистрирована`
var ErrUnsuitableParcelStatus = errors.New("статус посылки не равен `зарегистрирована`")