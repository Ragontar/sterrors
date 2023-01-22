# Errors PKG

## Описание
Пакет для стандартизации работы с ошибками. Он содержит интерфейс и типы ошибок, 
а также методы и конструкторы для удобства работы с ними. Основные фичи:
- Stacktrace
- Метод, инкапсулирующий оборачивание ошибок
- Стандартизация "архетипов" ошибок
- Создание лейблов для Loki
- Формирование полноценного отчета при правильном использовании
- Бойлерплейт

## Как использовать

Каждая возникающая ошибка должна создаваться с помощью конструктора NewXXXError.

Конструктор может принимать структуру BasicLabels. (опционально, но рекомендовано ее заполнять) 
Описание стандартных лейблов, часть из которых заплняется в этой структуре, ниже.

Если ошибка возвращается из какого-либо метода (например, в юзкейс вернулась ошибка
из репозитория), то должна быть создана новая ошибка через коструктор NewXXXError,
а вернувшаяся - обернута методом Wrap. (см. пример) Wrap - должен быть последним методом
в цепочке, т.к. возвращает интерфейс error, что не позволит чейнить методы дальше.

Если необходимо записать стактрейс, можно использовать WithStackTrace(). 
Рекомендуется использовать его на как можно более "низком" уровне.

## Лейблы и все-все-все

Каждая ошибка строится на основе GenericError. Она содержит 4 основных лейбла:
- Type - тип ошибки, заполняется автоматически в конструкторе.
- Context - контекст ошибки, основная задача - ее описание для просматривающего лог.
    Например, можно записать id транзакции, пользователя, etc.
- Location - метод, где произошла ошибка. Заполняется автоматически.
- Level - уровень, на котором произошла ошибка. (пример: usecase, repository, delivery)

На основе этих четырех лейблов формируется текст ошибки. Пример со стактрейсом:

```
TYPE: repoErr
CONTEXT: some context about sql error
LOCATION: main.repoFoo
LEVEL: repository
/Users/Kurisutina/Documents/work/errors_pkg/cmd/main.go:34
        main.repoFoo
/Users/Kurisutina/Documents/work/errors_pkg/cmd/main.go:42
        main.usecaseFoo
/Users/Kurisutina/Documents/work/errors_pkg/cmd/main.go:11
        main.main
/usr/local/go/src/runtime/proc.go:250
        runtime.main
/usr/local/go/src/runtime/asm_amd64.s:1594
        runtime.goexit
```

Если какие-то специфические ошибки необходимо мониторить отдельно, то можно
создать дополнительный лейбл методом SetLabel(key, val). Эти лейблы не фигурируют
в тексте ошибок и являются метаданными для Loki.

## Создание своей ошибки

Каждый тип ошибок должен быть максимально общим и, при этом, репрезентативным. 
Если появилась необходимость добавить новый тип, необходимо создать структуру 
в errorTypes.go и написать конструктор по аналогии. В конструкторе задается лейбл
TYPE.

## Правила использования пакета

1. Выбирать типы ошибок и заполнять их вдумчиво.
2. Всегда оборачивать ошибки на каждом уровне. Простой return err не канает, иначе
    итоговый отчет в локи выйдет неполным и малоинформативным.
3. Для заполнения лейблов (кроме контекста) использовать энумераторы из lokiLabels.go.
    Если необходимо - дополнить. Это критично для мониторинга. (особенно касается key лейбла)
4. Не плодить лишнего без необходимости, особенно типы. 

## Пример

Так выглядит обработка и создание ошибок на примере цепочки (usecase -> repo -> sql-driver)

```go


func main() {
	fmt.Println("------------------------- UNWRAPPING --------------------------------")
	wrappedErrors := usecaseFoo()
	for ; wrappedErrors != nil; wrappedErrors = errors.Unwrap(wrappedErrors) {
		fmt.Println("---NEXT ERROR: ")
		/*
			В еррор хендлере, во время распаковки ошибок, будут заполняться 
		    лейблы и прочее для отправки в локи. (через type switch)
		*/
		fmt.Println(wrappedErrors)
	}
}

func sqlErrorFoo() error {
	return fmt.Errorf("ordinary sql error from outer package")
}

func repoFoo() error {
	/* ... */
	// ex: error from select
	err := sqlErrorFoo()
	if err != nil {
		return sterrors.NewRepositoryError(sterrors.BasicLabels{
			Context: "some context about sql error",
			Level:   "repository",
		}).
			WithStackTrace().
			Wrap(err)
	}

	return nil
}

func usecaseFoo() error {
	/* ... */
	err := repoFoo()
	if err != nil {
		return sterrors.NewWithdrawError(sterrors.BasicLabels{
			Context: "some context about withdraw (txId, client, etc)",
			Level:   "usecase",
		}).
			SetLabel("customLabelForTrackingInLoki", "tracking value").
			Wrap(err)
	}

	return nil
}

```

```text
------------------------- UNWRAPPING --------------------------------
---NEXT ERROR: 
TYPE: withdrawErr
CONTEXT: some context about withdraw (txId, client, etc)
LOCATION: main.usecaseFoo
LEVEL: usecase

---NEXT ERROR: 
TYPE: repoErr
CONTEXT: some context about sql error
LOCATION: main.repoFoo
LEVEL: repository
/Users/Kurisutina/Documents/work/errors_pkg/cmd/main.go:34
        main.repoFoo
/Users/Kurisutina/Documents/work/errors_pkg/cmd/main.go:42
        main.usecaseFoo
/Users/Kurisutina/Documents/work/errors_pkg/cmd/main.go:11
        main.main
/usr/local/go/src/runtime/proc.go:250
        runtime.main
/usr/local/go/src/runtime/asm_amd64.s:1594
        runtime.goexit

---NEXT ERROR: 
ordinary sql error from outer package
```