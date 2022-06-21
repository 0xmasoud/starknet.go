package types

import (
	"bytes"
	"encoding/json"
	"math/big"
	"strconv"
	"testing"
)

var (
	rawTest string = `{
		"max_felt": 3618502788666131213697322783095070105623107215331596699973092056135872020481,
		"felts": [
			{"value": "0x0", "expected": 0},
			{"value": "0x1277312773", "expected": 79309121395},
			{"value": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB", "expected": "1071767867375995473349368877325274214414350531515"},
			{"value": "2458502865976494910213617956670505342647705497324144349552978333078363662855", "expected": "2458502865976494910213617956670505342647705497324144349552978333078363662855"}
		]
	}`
)

var feltTest FeltTest

type FeltTest struct {
	MaxFelt *big.Int    `json:"max_felt"`
	Felts   []FeltValue `json:"felts"`
}

type FeltValue struct {
	Value    string `json:"value"`
	Expected *Felt  `json:"expected"`
}

func TestJSONUnmarshal(t *testing.T) {
	json.Unmarshal([]byte(rawTest), &feltTest)

	fetchedMax := &Felt{feltTest.MaxFelt}
	if fetchedMax.String() != MaxFelt.String() {
		t.Errorf("Incorrect unmarshal and for max felt: %v %v\n", MaxFelt, feltTest.MaxFelt)
	}

	for _, felt := range feltTest.Felts {
		f := StrToFelt(felt.Value)

		if f.Int.Cmp(felt.Expected.Int) != 0 {
			t.Errorf("Incorrect unmarshal and felt comparison: %v %v\n", f.Int, felt.Expected.Int)
		}

		if f.String() != felt.Expected.String() {
			t.Errorf("Incorrect unmarshal and hex comparison: %v %v\n", f.Int, felt.Expected.Int)
		}
	}
}

func TestJSONMarshal(t *testing.T) {
	var newFelts FeltTest
	var newBigs []*big.Int
	for i, felt := range feltTest.Felts {
		nb := new(big.Int).Add(big.NewInt(int64(i)+7), felt.Expected.Int)
		newBigs = append(newBigs, nb)

		felt.Expected.Int = nb
		newFelts.Felts = append(newFelts.Felts, felt)
	}

	raw, err := json.Marshal(newFelts)
	if err != nil {
		t.Errorf("Could not marshal felt: %v\n", err)
	}

	var newTest FeltTest
	json.Unmarshal(raw, &newTest)

	for _, nb := range newBigs {
		innerBytes := []byte(nb.String())
		result := bytes.Index(raw, innerBytes)

		if result <= 0 {
			t.Errorf("Could not marshal felt: %v\n", result)
		}
		if string(raw[result:result+len(innerBytes)]) != string(innerBytes) {
			t.Errorf("Could not marshal felt: %v\n", result)
		}
	}
}

func TestGQLUnmarshal(t *testing.T) {
	for _, felt := range feltTest.Felts {
		if err := new(Felt).UnmarshalGQL(felt.Value); err != nil {
			t.Errorf("Could not unmarshal GQL for felt: %v\n", err)
		}
	}
	if err := new(Felt).UnmarshalGQL(1000000); err != nil {
		t.Errorf("Could not unmarshal GQL for felt: %v\n", err)
	}
}

func TestGQLMarshal(t *testing.T) {
	for i, felt := range feltTest.Felts {
		buf := bytes.NewBuffer(nil)
		felt.Expected.MarshalGQL(buf)

		cmp := &Felt{Int: new(big.Int).Add(big.NewInt(int64(i)+7), ToFelt(felt.Value).Int)}

		if buf.String() != strconv.Quote(cmp.String()) {
			t.Errorf("Could not marshal GQL for felt: %v %v\n", buf.String(), strconv.Quote(cmp.String()))
		}
	}
}
