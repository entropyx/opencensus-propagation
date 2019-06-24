package propagation

import "net/http"

type HTTPHeader http.Header

var _ TextMap = HTTPHeader{}

func (h HTTPHeader) ForEach(handler func(key, val string) error) error {
	for k, v := range h {
		for _, v2 := range v {
			if err := handler(k, v2); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h HTTPHeader) Set(key, value string) {
	header := http.Header(h)
	header.Set(key, value)
}
