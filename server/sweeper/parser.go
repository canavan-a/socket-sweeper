package sweeper

import "encoding/json"

type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Parser struct{}

func (p *Parser) parseCoordinates(inputData []byte) (x, y int, err error) {

	var coords Coordinates

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(inputData, &coords)
	if err != nil {
		return
	}

	return coords.X, coords.Y, nil
}
