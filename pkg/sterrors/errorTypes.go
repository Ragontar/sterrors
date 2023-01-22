package sterrors

/*

В этом файле должны быть имплементированы все "архетипы" используемых ошибок. Это необходимо для возможности проверки
ошибки на соответствие какому-то конкретному типу, чтобы была возможность отвечать клиенту (например) статус-кодом,
отличным от 500. (i.e. если errors.As(err, &NotFoundError{}) == true, то вернуть 404)

Контекст ошибок позволит как улучшить информативность ответов, так и даст возможность отслеживать ошибки каких-то
конкретных типов для мониторинга/метрик/пр. Простор для идей широкий.

PS для создания ошибки конкретного типа !!обязательно!! использовать конструкторы, через них инициализируются стектрейсы.

*/

type BasicLabels struct {
	Context string
	Level   string
}

type WithdrawError struct {
	GenericError
}

func NewWithdrawError(basic ...BasicLabels) GenericErrorInterface {
	return &WithdrawError{newGenericError("withdrawErr", basic)}
}

type NotFoundError struct {
	GenericError
}

func NewNotFoundError(basic ...BasicLabels) GenericErrorInterface {
	return &NotFoundError{newGenericError("notFoundErr", basic)}
}

type RepositoryError struct {
	GenericError
}

func NewRepositoryError(basic ...BasicLabels) GenericErrorInterface {
	return &RepositoryError{newGenericError("repoErr", basic)}
}
