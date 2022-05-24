package parsers

import (
	"bytes"
	"context"
	"encoding/json"
)

// Item элемент списка приборов учета/записей архива показаний
type Item struct {
	// V результат парсинга элемента списка
	V interface{}

	// E ошибка парсинга элемента списка
	E error
}

func (i Item) Error() bool {
	return i.E != nil
}

func of(x interface{}) Item {
	return Item{V: x}
}

func e(err error) Item {
	return Item{E: err}
}

// ParseGaugesList разбирает ответ метода /api/cascade/counter-house
func ParseGaugesList(ctx context.Context, b []byte) (<-chan Item, error) {
	decoder := json.NewDecoder(bytes.NewReader(b))

	_, err := decoder.Token()

	if err != nil {
		return nil, err
	}

	out := make(chan Item)

	go func(decoder *json.Decoder, b []byte) {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return

			default:
				if !decoder.More() {
					return
				}

				var counterHouse Gauge

				if err := decoder.Decode(&counterHouse); err != nil {
					out <- e(err)
				} else {
					out <- of(&counterHouse)
				}
			}
		}
	}(decoder, b)

	return out, nil
}

// ParseReadings разбирает ответ метода /api/cascade/counter-house/readings
func ParseReadings(ctx context.Context, b []byte) (<-chan Item, error) {
	decoder := json.NewDecoder(bytes.NewReader(b))

	_, err := decoder.Token()

	if err != nil {
		return nil, err
	}

	out := make(chan Item)

	go func(decoder *json.Decoder, b []byte) {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return

			default:
				if !decoder.More() {
					return
				}

				var reading Readings

				if err := decoder.Decode(&reading); err != nil {
					out <- e(err)
				} else {
					out <- of(&reading)
				}
			}
		}
	}(decoder, b)

	return out, nil
}
