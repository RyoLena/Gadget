package converter

import "time"

type Time2String struct {
	Pattern string
}

func (t Time2String) Convert(src time.Time) (string, error) {
	return src.Format(t.Pattern), nil
}
