package cascade

import (
	"bytes"
	"encoding/json"
)

// ParseCounterHouseDto разбирает ответ метода /api/cascade/counter-house
func ParseCounterHouseDto(b []byte) <-chan struct {
	*CounterHouseDto
	error
} {
	out := make(chan struct {
		*CounterHouseDto
		error
	})

	go func() {
		defer close(out)

		decoder := json.NewDecoder(bytes.NewReader(b))

		_, err := decoder.Token()

		if err != nil {
			out <- struct {
				*CounterHouseDto
				error
			}{nil, err}

			return
		}

		for decoder.More() {
			var counterHouse CounterHouseDto

			if err := decoder.Decode(&counterHouse); err != nil {
				out <- struct {
					*CounterHouseDto
					error
				}{nil, err}
				break
			} else {
				out <- struct {
					*CounterHouseDto
					error
				}{&counterHouse, nil}
			}
		}
	}()

	return out
}

// ParseCounterHouseReadingDto разбирает ответ метода /api/cascade/counter-house/readings
func ParseCounterHouseReadingDto(b []byte) <-chan struct {
	*CounterHouseReadingDto
	error
} {
	out := make(chan struct {
		*CounterHouseReadingDto
		error
	})

	go func() {
		defer close(out)

		decoder := json.NewDecoder(bytes.NewReader(b))

		_, err := decoder.Token()

		if err != nil {
			out <- struct {
				*CounterHouseReadingDto
				error
			}{nil, err}

			return
		}

		for decoder.More() {
			var reading CounterHouseReadingDto

			if err := decoder.Decode(&reading); err != nil {
				out <- struct {
					*CounterHouseReadingDto
					error
				}{nil, err}
				break
			} else {
				out <- struct {
					*CounterHouseReadingDto
					error
				}{&reading, nil}
			}
		}
	}()

	return out
}
