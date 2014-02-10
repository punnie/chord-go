package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

func readMessage(reader *bufio.Reader) ([]byte, error) {
	payloadSize, err := readSize(reader)

	if err != nil {
		return nil, err
	}

	payload, err := readPayload(reader, payloadSize)

	if err != nil {
		return nil, err
	}

	return payload, nil
}

func readSize(reader *bufio.Reader) (uint32, error) {
	size := new(uint32)
	sizeBuffer := make([]byte, 4)

	for i := 0; i < 4; {
		n, err := reader.Read(sizeBuffer) // TODO: this has bugs

		if err != nil {
			return 0, err
		}

		i += n
	}

	intReader := bytes.NewReader(sizeBuffer)
	err := binary.Read(intReader, binary.LittleEndian, size)

	if err != nil {
		return 0, err
	}

	return *size, nil
}

func readPayload(reader *bufio.Reader, size uint32) ([]byte, error) {
	payload := new(bytes.Buffer)
	payloadBuffer := make([]byte, size)

	for i := 0; i < int(size); {
		n, err := reader.Read(payloadBuffer)

		if err != nil {
			return nil, err
		}

		payload.Write(payloadBuffer[i : i+n])

		i += n
	}

	return payload.Bytes(), nil
}
