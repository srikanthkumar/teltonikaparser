// Copyright 2019 Filip Kroča. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package teltonikaparser

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/filipkroca/b2n"
	"github.com/filipkroca/teltonikaparser/teltonikajson"
)

// HAvlData represent human readable set of a pointer to an AvlEncodeKey Decoding key and a pointer to IO element with RAW data
type HAvlData struct {
	AvlEncodeKey *AvlEncodeKey
	Element      *Element
}

// HumanDecoder is responsible for decoding
type HumanDecoder struct {
	elements map[string]map[uint16]AvlEncodeKey
}

// AvlEncodeKey represent parsed element values from JSON
type AvlEncodeKey struct {
	No              string `json:"No"`
	PropertyName    string `json:"PropertyName"`
	Bytes           string `json:"Bytes"`
	Type            string `json:"Type"`
	Min             string `json:"Min"`
	Max             string `json:"Max"`
	Multiplier      string `json:"Multiplier"`
	Units           string `json:"Units"`
	Description     string `json:"Description"`
	HWSupport       string `json:"HWSupport"`
	ParametrGroup   string `json:"Parametr Group"`
	FinalConversion string `json:"FinalConversion"`
}

// Human takes a pointer to Element, device type ["FMBXY", "FM64"] and return a pointer to decoding key
func (h *HumanDecoder) Human(el *Element, device string) (*HAvlData, error) {
	//init decoding key
	if len(h.elements) == 0 {
		h.loadElements()
	}

	// check if Element is valid
	if !((*el).Length > 0 && (*el).IOID > 0 && len((*el).Value) > 0) {
		//log.Fatal("Unable to decode empty element")
		return nil, errors.New("Unable to decode empty element")
	}

	// find decode key and pair it
	avl, ok := h.elements[device][(*el).IOID]
	if !ok {
		//log.Fatalf("Unknown element %v", (*el).IOID)
		return nil, fmt.Errorf("Unknown element %v", (*el).IOID)
	}

	// return pointer to merged struct with decode key AvlEncodeKey and data Element
	havl := HAvlData{
		AvlEncodeKey: &avl,
		Element:      el,
	}
	return &havl, nil
}

// loadElements parses ./decoding/.. into slice
func (h *HumanDecoder) loadElements() {
	// make map
	h.elements = make(map[string]map[uint16]AvlEncodeKey)

	// read our opened xmlFile as a byte array.
	byteValue := []byte(teltonikajson.FMBXY)
	fmbxy := make(map[uint16]AvlEncodeKey)
	//h.elements["FMBXY"] = make(map[uint16]AvlEncodeKey)
	err := json.Unmarshal(byteValue, &fmbxy)
	if err != nil {
		log.Panic(err)
	}
	h.elements["FMBXY"] = fmbxy

	// read our opened xmlFile as a byte array.
	byteValue = []byte(teltonikajson.FM64)
	fm64 := make(map[uint16]AvlEncodeKey)
	//h.elements["FMBXY"] = make(map[uint16]AvlEncodeKey)
	err = json.Unmarshal(byteValue, &fm64)
	if err != nil {
		log.Panic(err)
	}
	h.elements["FM64"] = fm64

}

// GetFinalValue return decimal value, if necesarry with float, return should be empty interface because there is many values to return
func (h *HAvlData) GetFinalValue() (interface{}, error) {

	if h.AvlEncodeKey.FinalConversion == "toUint8" {
		if h.AvlEncodeKey.Bytes != "1" || h.AvlEncodeKey.Type != "Unsigned" || len(h.Element.Value) != 1 {
			return nil, fmt.Errorf("Unable to convert %vBytes long parametr, %vBytes real long parametr to Uint8 %v", h.AvlEncodeKey.Bytes, len(h.Element.Value), h.AvlEncodeKey.PropertyName)
		}
		return b2n.ParseBs2Uint8(&h.Element.Value, 0)
	}

	if h.AvlEncodeKey.FinalConversion == "toUint16" {
		if h.AvlEncodeKey.Bytes != "2" || h.AvlEncodeKey.Type != "Unsigned" || len(h.Element.Value) != 2 {
			return nil, fmt.Errorf("Unable to convert %vBytes long parametr, %vBytes real long parametr to Uint16 %v", h.AvlEncodeKey.Bytes, len(h.Element.Value), h.AvlEncodeKey.PropertyName)
		}
		return b2n.ParseBs2Uint16(&h.Element.Value, 0)
	}

	if h.AvlEncodeKey.FinalConversion == "toUint32" {
		if h.AvlEncodeKey.Bytes != "4" || h.AvlEncodeKey.Type != "Unsigned" || len(h.Element.Value) != 4 {
			return nil, fmt.Errorf("Unable to convert %vBytes long parametr, %vBytes real long parametr to Uint32 %v", h.AvlEncodeKey.Bytes, len(h.Element.Value), h.AvlEncodeKey.PropertyName)
		}
		return b2n.ParseBs2Uint32(&h.Element.Value, 0)
	}

	if h.AvlEncodeKey.FinalConversion == "toUint64" {
		if h.AvlEncodeKey.Bytes != "8" || h.AvlEncodeKey.Type != "Unsigned" || len(h.Element.Value) != 8 {
			return nil, fmt.Errorf("Unable to convert %vBytes long parametr, %vBytes real long parametr to Uint64 %v", h.AvlEncodeKey.Bytes, len(h.Element.Value), h.AvlEncodeKey.PropertyName)
		}
		return b2n.ParseBs2Uint64(&h.Element.Value, 0)
	}

	if h.AvlEncodeKey.FinalConversion == "toInt8" {
		if h.AvlEncodeKey.Bytes != "1" || h.AvlEncodeKey.Type != "Signed" || len(h.Element.Value) != 1 {
			return nil, fmt.Errorf("Unable to convert %vBytes long parametr, %vBytes real long parametr to Int8 %v", h.AvlEncodeKey.Bytes, len(h.Element.Value), h.AvlEncodeKey.PropertyName)
		}
		return b2n.ParseBs2Int8TwoComplement(&h.Element.Value, 0)
	}

	if h.AvlEncodeKey.FinalConversion == "toInt16" {
		if h.AvlEncodeKey.Bytes != "2" || h.AvlEncodeKey.Type != "Signed" || len(h.Element.Value) != 2 {
			return nil, fmt.Errorf("Unable to convert %vBytes long parametr, %vBytes real long parametr to Int16 %v", h.AvlEncodeKey.Bytes, len(h.Element.Value), h.AvlEncodeKey.PropertyName)
		}
		return b2n.ParseBs2Int16TwoComplement(&h.Element.Value, 0)
	}

	if h.AvlEncodeKey.FinalConversion == "toInt32" {
		if h.AvlEncodeKey.Bytes != "4" || h.AvlEncodeKey.Type != "Signed" || len(h.Element.Value) != 4 {
			return nil, fmt.Errorf("Unable to convert %vBytes long parametr, %vBytes real long parametr to Int32 %v", h.AvlEncodeKey.Bytes, len(h.Element.Value), h.AvlEncodeKey.PropertyName)
		}
		return b2n.ParseBs2Int32TwoComplement(&h.Element.Value, 0)
	}

	if h.AvlEncodeKey.FinalConversion == "toInt64" {
		if h.AvlEncodeKey.Bytes != "8" || h.AvlEncodeKey.Type != "Signed" || len(h.Element.Value) != 8 {
			return nil, fmt.Errorf("Unable to convert %vBytes long parametr, %vBytes real long parametr, to Int64 %v", h.AvlEncodeKey.Bytes, len(h.Element.Value), h.AvlEncodeKey.PropertyName)
		}
		return b2n.ParseBs2Int64TwoComplement(&h.Element.Value, 0)
	}

	return string(h.Element.Value), nil
}
