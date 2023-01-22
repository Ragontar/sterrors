package sterrors

type Labeler interface {
	Labels() []Label
	LabelsMap() map[string]string
	addLabel(label Label)
}

type Label struct {
	Key   string
	Value string
}

type labels struct {
	ls map[string]string
}

// Labels возвращает все метки в виде массива
func (l *labels) Labels() []Label {
	var res []Label
	for key, val := range l.ls {
		res = append(res, Label{
			Key:   key,
			Value: val,
		})
	}
	return res
}

// LabelsMap возвращает метки в виде карты
func (l *labels) LabelsMap() map[string]string {
	return l.ls
}

func (l *labels) addLabel(label Label) {
	if l.ls == nil {
		l.ls = make(map[string]string)
	}
	l.ls[label.Key] = label.Value
}
