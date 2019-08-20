package util

import (
	"fmt"
	"testing"
)

func TestCheckSimilarStruct(t *testing.T) {
	type FF struct {
		F string `json:"f"`
	}

	type TestStruct struct {
		A string `json:"a"`
		B int    `json:"b"`
		C struct {
			Aaa string `json:"aaa"`
		} `json:"c"`
		D interface{} `json:"d"`
		E bool        `json:"e"`
		F []FF        `json:"f"`
	}

	x := TestStruct{
		A: "ardo",
		C: struct {
			Aaa string `json:"aaa"`
		}{
			Aaa: "asuu",
		},
		D: struct {
			Aaa string `json:"aaa"`
		}{
			Aaa: "asuu",
		},
		E: false,
		F: []FF{
			FF{
				F: "test",
			},
		},
	}
	y := TestStruct{
		A: "ardo",
		C: struct {
			Aaa string `json:"aaa"`
		}{
			Aaa: "asuu",
		},
		D: struct {
			Aaa string `json:"aaa"`
		}{
			Aaa: "asuu",
		},
		E: false,
		F: []FF{
			FF{
				F: "test",
			},
		},
	}

	fmt.Println(CheckSimilarStruct(x, y))
}

func TestStringContaintsInSlice(t *testing.T)  {
	str := "indonesia"
	strArray := []string{"Indonesia", "Belgia"}

	ok := StringContaintsInSlice(str, strArray)

	if !ok {
		t.Errorf("Must True")
	}
}
