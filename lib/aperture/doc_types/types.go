package doc_types

import (
	"github.com/goaperture/goaperture/lib/aperture"
	"github.com/goaperture/goaperture/lib/aperture/types"
)

type TestItem[P types.Payload] struct {
	Path string
	Test func(client aperture.Client[P]) TestData
}

type TestData struct {
	Inputs  []any
	Outputs []any
}
