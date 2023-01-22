package main

import (
	"errors"
	"errors_pkg/pkg/sterrors"
	"fmt"
)

func main() {
	fmt.Println("------------------------- UNWRAPPING --------------------------------")
	wrappedErrors := usecaseFoo()
	for ; wrappedErrors != nil; wrappedErrors = errors.Unwrap(wrappedErrors) {
		fmt.Println("---NEXT ERROR: ")
		/*
			В еррор хендлере, во время распаковки ошибок, будут заполняться лейблы и прочее для отправки в локи.
			(через type switch)
		*/
		fmt.Println(wrappedErrors)
	}
}

func sqlErrorFoo() error {
	return fmt.Errorf("ordinary sql error from outer package")
}

func repoFoo() error {
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
