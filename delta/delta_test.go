package delta

import (
	"bufio"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/deathtrix/rdf/model"
)

func TestDiffOverlap100(t *testing.T) {
	blockSize := 3
	textOld := "123456789"
	textNew := "123456789"
	diffComp := &model.Delta{
		Blocks: []model.Block{
			{Op: "O", FromStart: 0, FromEnd: 3, ToStart: 0, ToEnd: 3},
			{Op: "O", FromStart: 3, FromEnd: 6, ToStart: 3, ToEnd: 6},
			{Op: "O", FromStart: 6, FromEnd: 9, ToStart: 6, ToEnd: 9},
		},
	}

	readerOld := bufio.NewReaderSize(strings.NewReader(textOld), 1)
	sig, err := BuildSignature(readerOld, blockSize)
	if err != nil {
		t.Fatal(err)
	}

	readerNew := bufio.NewReaderSize(strings.NewReader(textNew), 1)
	diff, err := BuildDelta(readerNew, sig)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(diff, diffComp) {
		t.Fatalf("expected: %v , got %v", diffComp, diff)
	}
}

func TestDiffOverlap0(t *testing.T) {
	blockSize := 3
	textOld := "123456789"
	textNew := "987654321"
	diffComp := &model.Delta{
		Blocks: []model.Block{
			{Op: "C", FromStart: 0, FromEnd: 0, ToStart: 0, ToEnd: 3},
			{Op: "C", FromStart: 0, FromEnd: 0, ToStart: 3, ToEnd: 6},
			{Op: "C", FromStart: 0, FromEnd: 0, ToStart: 6, ToEnd: 9},
		},
	}

	readerOld := bufio.NewReaderSize(strings.NewReader(textOld), 1)
	sig, err := BuildSignature(readerOld, blockSize)
	if err != nil {
		t.Fatal(err)
	}

	readerNew := bufio.NewReaderSize(strings.NewReader(textNew), 1)
	diff, err := BuildDelta(readerNew, sig)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(diff, diffComp) {
		t.Fatalf("expected: %v , got %v", diffComp, diff)
	}
}

func TestDiffChange(t *testing.T) {
	blockSize := 3
	textOld := "123456789"
	textNew := "123356789"
	diffComp := &model.Delta{
		Blocks: []model.Block{
			{Op: "O", FromStart: 0, FromEnd: 3, ToStart: 0, ToEnd: 3},
			{Op: "O", FromStart: 6, FromEnd: 9, ToStart: 6, ToEnd: 9},
			{Op: "C", FromStart: 0, FromEnd: 0, ToStart: 3, ToEnd: 6},
		},
	}

	readerOld := bufio.NewReaderSize(strings.NewReader(textOld), 1)
	sig, err := BuildSignature(readerOld, blockSize)
	if err != nil {
		t.Fatal(err)
	}

	readerNew := bufio.NewReaderSize(strings.NewReader(textNew), 1)
	diff, err := BuildDelta(readerNew, sig)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(diff, diffComp) {
		t.Fatalf("expected: %v , got %v", diffComp, diff)
	}
}

func TestDiffInsertBlock(t *testing.T) {
	blockSize := 3
	textOld := "123456789"
	textNew := "123abc456789"
	diffComp := &model.Delta{
		Blocks: []model.Block{
			{Op: "O", FromStart: 0, FromEnd: 3, ToStart: 0, ToEnd: 3},
			{Op: "M", FromStart: 3, FromEnd: 6, ToStart: 6, ToEnd: 9},
			{Op: "M", FromStart: 6, FromEnd: 9, ToStart: 9, ToEnd: 12},
			{Op: "C", FromStart: 0, FromEnd: 0, ToStart: 3, ToEnd: 6},
		},
	}

	readerOld := bufio.NewReaderSize(strings.NewReader(textOld), 1)
	sig, err := BuildSignature(readerOld, blockSize)
	if err != nil {
		t.Fatal(err)
	}

	readerNew := bufio.NewReaderSize(strings.NewReader(textNew), 1)
	diff, err := BuildDelta(readerNew, sig)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(diff, diffComp) {
		t.Fatalf("expected: %v , got %v", diffComp, diff)
	}
}

func TestDiffInsertLessThanBlock(t *testing.T) {
	blockSize := 3
	textOld := "123456789"
	textNew := "123a456789"
	diffComp := &model.Delta{
		Blocks: []model.Block{
			{Op: "O", FromStart: 0, FromEnd: 3, ToStart: 0, ToEnd: 3},
			{Op: "M", FromStart: 3, FromEnd: 6, ToStart: 4, ToEnd: 7},
			{Op: "M", FromStart: 6, FromEnd: 9, ToStart: 7, ToEnd: 10},
			{Op: "C", FromStart: 0, FromEnd: 0, ToStart: 3, ToEnd: 4},
		},
	}

	readerOld := bufio.NewReaderSize(strings.NewReader(textOld), 1)
	sig, err := BuildSignature(readerOld, blockSize)
	if err != nil {
		t.Fatal(err)
	}

	readerNew := bufio.NewReaderSize(strings.NewReader(textNew), 1)
	diff, err := BuildDelta(readerNew, sig)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(diff, diffComp) {
		t.Fatalf("expected: %v , got %v", diffComp, diff)
	}
}

func TestDiffDeleteBlock(t *testing.T) {
	blockSize := 3
	textOld := "123abc456789"
	textNew := "123456789"
	diffComp := &model.Delta{
		Blocks: []model.Block{
			{Op: "O", FromStart: 0, FromEnd: 3, ToStart: 0, ToEnd: 3},
			{Op: "M", FromStart: 6, FromEnd: 9, ToStart: 3, ToEnd: 6},
			{Op: "M", FromStart: 9, FromEnd: 12, ToStart: 6, ToEnd: 9},
		},
	}

	readerOld := bufio.NewReaderSize(strings.NewReader(textOld), 1)
	sig, err := BuildSignature(readerOld, blockSize)
	if err != nil {
		t.Fatal(err)
	}

	readerNew := bufio.NewReaderSize(strings.NewReader(textNew), 1)
	diff, err := BuildDelta(readerNew, sig)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(diff, diffComp) {
		t.Fatalf("expected: %v , got %v", diffComp, diff)
	}
}

func TestDiffDeleteLessThanBlock(t *testing.T) {
	blockSize := 3
	textOld := "123b456789"
	textNew := "123456789"
	diffComp := &model.Delta{
		Blocks: []model.Block{
			{Op: "O", FromStart: 0, FromEnd: 3, ToStart: 0, ToEnd: 3},
			{Op: "M", FromStart: 6, FromEnd: 9, ToStart: 5, ToEnd: 8},
			{Op: "M", FromStart: 9, FromEnd: 12, ToStart: 8, ToEnd: 9},
			{Op: "C", FromStart: 0, FromEnd: 0, ToStart: 3, ToEnd: 5},
		},
	}

	readerOld := bufio.NewReaderSize(strings.NewReader(textOld), 1)
	sig, err := BuildSignature(readerOld, blockSize)
	if err != nil {
		t.Fatal(err)
	}

	readerNew := bufio.NewReaderSize(strings.NewReader(textNew), 1)
	diff, err := BuildDelta(readerNew, sig)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(diff, diffComp) {
		t.Fatalf("expected: %v , got %v", diffComp, diff)
	}
}

func TestDiffSubstitute(t *testing.T) {
	blockSize := 3
	textOld := "123456789"
	textNew := "789456123"
	diffComp := &model.Delta{
		Blocks: []model.Block{
			{Op: "M", FromStart: 0, FromEnd: 3, ToStart: 6, ToEnd: 9},
			{Op: "O", FromStart: 3, FromEnd: 6, ToStart: 3, ToEnd: 6},
			{Op: "M", FromStart: 6, FromEnd: 9, ToStart: 0, ToEnd: 3},
		},
	}

	readerOld := bufio.NewReaderSize(strings.NewReader(textOld), 1)
	sig, err := BuildSignature(readerOld, blockSize)
	if err != nil {
		t.Fatal(err)
	}

	readerNew := bufio.NewReaderSize(strings.NewReader(textNew), 1)
	diff, err := BuildDelta(readerNew, sig)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(diff, diffComp) {
		t.Fatalf("expected: %v , got %v", diffComp, diff)
	}
}

func TestDiffRepeatBlock(t *testing.T) {
	blockSize := 3
	textOld := "123456789"
	textNew := "123123123"
	diffComp := &model.Delta{
		Blocks: []model.Block{
			{Op: "O", FromStart: 0, FromEnd: 3, ToStart: 0, ToEnd: 3},
			{Op: "M", FromStart: 0, FromEnd: 3, ToStart: 3, ToEnd: 6},
			{Op: "M", FromStart: 0, FromEnd: 3, ToStart: 6, ToEnd: 9},
		},
	}

	readerOld := bufio.NewReaderSize(strings.NewReader(textOld), 1)
	sig, err := BuildSignature(readerOld, blockSize)
	if err != nil {
		t.Fatal(err)
	}

	readerNew := bufio.NewReaderSize(strings.NewReader(textNew), 1)
	diff, err := BuildDelta(readerNew, sig)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(diff, diffComp) {
		t.Fatalf("expected: %v , got %v", diffComp, diff)
	}
}

func TestDiffRepeatChar(t *testing.T) {
	blockSize := 3
	textOld := "aaaaaaaaa"
	textNew := "aaaaaaaaa"
	diffComp := &model.Delta{
		Blocks: []model.Block{
			{Op: "O", FromStart: 0, FromEnd: 3, ToStart: 0, ToEnd: 3},
			{Op: "M", FromStart: 0, FromEnd: 3, ToStart: 3, ToEnd: 6},
			{Op: "M", FromStart: 0, FromEnd: 3, ToStart: 6, ToEnd: 9},
		},
	}

	readerOld := bufio.NewReaderSize(strings.NewReader(textOld), 1)
	sig, err := BuildSignature(readerOld, blockSize)
	if err != nil {
		t.Fatal(err)
	}

	readerNew := bufio.NewReaderSize(strings.NewReader(textNew), 1)
	diff, err := BuildDelta(readerNew, sig)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(diff, diffComp) {
		t.Fatalf("expected: %v , got %v", diffComp, diff)
	}
}

func BenchmarkDiffNonZero(b *testing.B) {
	blockSize := 2
	bytes1 := []byte{}
	bytes2 := []byte{}
	textOld := string(bytes1)
	textNew := string(bytes2)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 1_000_000_000; i++ {
		bytes1 = append(bytes1, byte(r.Uint32()))
		bytes2 = append(bytes2, byte(r.Uint32()))
	}

	readerOld := bufio.NewReaderSize(strings.NewReader(textOld), 1)
	readerNew := bufio.NewReaderSize(strings.NewReader(textNew), 1)

	b.ResetTimer()
	sig, err := BuildSignature(readerOld, blockSize)
	if err != nil {
		b.Fatal(err)
	}

	_, err = BuildDelta(readerNew, sig)
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkDiffOnlyZeros(b *testing.B) {
	blockSize := 2
	bytes := make([]byte, 100_000_000)
	textOld := string(bytes)
	textNew := string(bytes)

	readerOld := bufio.NewReaderSize(strings.NewReader(textOld), 1)
	readerNew := bufio.NewReaderSize(strings.NewReader(textNew), 1)

	b.ResetTimer()
	sig, err := BuildSignature(readerOld, blockSize)
	if err != nil {
		b.Fatal(err)
	}

	_, err = BuildDelta(readerNew, sig)
	if err != nil {
		b.Fatal(err)
	}
}
