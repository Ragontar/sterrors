package main

import (
	"errors"
	"errors_pkg/pkg/somePkg"
	"errors_pkg/pkg/sterrors"
	"fmt"
)

func main() {
	errWithoutBasic := sterrors.NewWithdrawError()
	fmt.Println(errWithoutBasic.Error())

	errWithBasic := sterrors.NewWithdrawError(sterrors.BasicLabels{
		Context: "SOME BASIC CTX",
		Level:   "main",
	})

	fmt.Println(errWithBasic)

	fmt.Println(errorWithStackTraceIsHere().Error())

	fmt.Println(somePkg.FooWithErrorFromAnotherPkg().Error())

	wrappedErrors := WrappedErrorsHere()

	fmt.Println("------------------------- UNWRAPPING --------------------------------")
	for ; wrappedErrors != nil; wrappedErrors = errors.Unwrap(wrappedErrors) {
		fmt.Println("ERROR FOUND: ")
		fmt.Println(wrappedErrors)
	}
}

func errorWithStackTraceIsHere() error {
	return sterrors.NewNotFoundError(sterrors.BasicLabels{
		Context: "PIZDEC",
		Level:   "main",
	}).
		WithStackTrace().
		SetLabel("customLabel", "228").
		Wrap(somePkg.FooWithErrorFromAnotherPkg())
}

func WrappedErrorsHere() error {
	return sterrors.NewNotFoundError(sterrors.BasicLabels{
		Context: "wrapped errors of different types",
		Level:   "main",
	}).Wrap(errorWithStackTraceIsHere())
}
