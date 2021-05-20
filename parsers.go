package cascade

import (
	"bytes"
	"context"
	"encoding/json"
)

// ParseCounterHouseDto разбирает ответ метода /api/cascade/counter-house
func ParseCounterHouseDto(ctx context.Context, b []byte) <-chan struct {
	*CounterHouseDto
	error
} {
	out := make(chan struct {
		*CounterHouseDto
		error
	})

	go func(b []byte) {
		defer close(out)

		decoder := json.NewDecoder(bytes.NewReader(b))

		_, err := decoder.Token()

		if err != nil {
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !decoder.More() {
					return
				}

				var counterHouse CounterHouseDto

				if err := decoder.Decode(&counterHouse); err != nil {
					out <- struct {
						*CounterHouseDto
						error
					}{nil, err}
				} else {
					out <- struct {
						*CounterHouseDto
						error
					}{&counterHouse, nil}
				}
			}
		}
	}(b)

	return out
}

// ParseCounterHouseReadingDto разбирает ответ метода /api/cascade/counter-house/readings
func ParseCounterHouseReadingDto(ctx context.Context, b []byte) <-chan struct {
	*CounterHouseReadingDto
	error
} {
	out := make(chan struct {
		*CounterHouseReadingDto
		error
	})

	go func(b []byte) {
		defer close(out)

		decoder := json.NewDecoder(bytes.NewReader(b))

		_, err := decoder.Token()

		if err != nil {
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !decoder.More() {
					return
				}

				var reading CounterHouseReadingDto

				if err := decoder.Decode(&reading); err != nil {
					out <- struct {
						*CounterHouseReadingDto
						error
					}{nil, err}
				} else {
					out <- struct {
						*CounterHouseReadingDto
						error
					}{&reading, nil}
				}
			}
		}
	}(b)

	return out
}
