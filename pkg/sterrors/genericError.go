package sterrors

import (
	"fmt"
	"runtime"
)

type GenericErrorInterface interface {
	Error() string
	Wrap(err error) GenericErrorInterface
	Unwrap() error
	WithStackTrace() GenericErrorInterface

	SetLabel(key, val string) GenericErrorInterface
	SetContext(context string) GenericErrorInterface
	SetLevel(level string) GenericErrorInterface

	Labeler
	StackTracer
}

type GenericError struct {
	wrapped error
	labels
	stackTrace
}

func (e GenericError) Error() string {
	var errorText string
	lm := e.labels.LabelsMap()
	if t, ok := lm[KEY_ERROR_TYPE]; ok {
		errorText += fmt.Sprintf("TYPE: %s\n", t)
	}
	if ctx, ok := lm[KEY_ERROR_CONTEXT]; ok {
		errorText += fmt.Sprintf("CONTEXT: %s\n", ctx)
	}
	if loc, ok := lm[KEY_ERROR_LOCATION]; ok {
		errorText += fmt.Sprintf("LOCATION: %s\n", loc)
	}
	if lvl, ok := lm[KEY_ERROR_LEVEL]; ok {
		errorText += fmt.Sprintf("LEVEL: %s\n", lvl)
	}
	if e.Trace() != "" {
		errorText += e.Trace()
	}
	return errorText
}

// WithStackTrace добавляет стактрейс к тексту ошибки
func (e GenericError) WithStackTrace() GenericErrorInterface {
	e.stackTrace = newStackTrace()
	return &e
}

// SetContext устанавливает контекст ошибики, можно использовать для удобства
func (e GenericError) SetContext(context string) GenericErrorInterface {
	e.addLabel(Label{
		Key:   KEY_ERROR_CONTEXT,
		Value: context,
	})

	return &e
}

// SetLevel устанавливает уровень ошибики, можно использовать для удобства
func (e GenericError) SetLevel(level string) GenericErrorInterface {
	e.addLabel(Label{
		Key:   KEY_ERROR_LEVEL,
		Value: level,
	})

	return &e
}

func newGenericError(errType string, basic []BasicLabels) GenericError {
	var labels []Label
	if len(basic) == 0 {
		labels = []Label{
			{
				Key:   KEY_ERROR_TYPE,
				Value: errType,
			},
			{
				Key:   KEY_ERROR_CONTEXT,
				Value: "unknown",
			},
			{
				Key:   KEY_ERROR_LEVEL,
				Value: "unknown",
			},
			{
				Key:   KEY_ERROR_LOCATION,
				Value: getLocation(),
			},
		}
	} else {
		labels = []Label{
			{
				Key:   KEY_ERROR_TYPE,
				Value: errType,
			},
			{
				Key:   KEY_ERROR_CONTEXT,
				Value: basic[0].Context,
			},
			{
				Key:   KEY_ERROR_LEVEL,
				Value: basic[0].Level,
			},
			{
				Key:   KEY_ERROR_LOCATION,
				Value: getLocation(),
			},
		}
	}

	ge := GenericError{}
	for _, l := range labels {
		ge.addLabel(l)
	}

	return ge
}

// SetLabel добавляет метку ошибки для использования в Loki. На текст ошибки не влияет.
func (e GenericError) SetLabel(key, val string) GenericErrorInterface {
	e.addLabel(Label{
		Key:   key,
		Value: val,
	})

	return &e
}

func getLocation() string {
	pc, _, _, ok := runtime.Caller(3)
	if !ok {
		return "getLocation error"
	}
	fn := runtime.FuncForPC(pc)

	return fn.Name()
}

// Wrap оборачивает ошибку. Использовать всегда, кроме случаев, когда ошибка создается текущем уровне и оборачивать нечего.
func (e GenericError) Wrap(err error) GenericErrorInterface {
	e.wrapped = err
	return &e
}

func (e GenericError) Unwrap() error {
	return e.wrapped
}
