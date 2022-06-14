package server

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestWriteBack(t *testing.T) {
	reader := bytes.NewReader([]byte("hello world yet again"))
	buffer := bytes.NewBuffer([]byte{})
	err := (&Server{opt: &Option{BufferSize: 10}}).createStream(reader, buffer, true)
	if err == io.EOF {
		fmt.Println(buffer.String())
	} else if err != nil {
		t.Fatal(err)
	}
}
