package form3_test

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type JsonMarshalMock struct {
	mock.Mock
}

func (m *JsonMarshalMock) Marshal(v any) ([]byte, error) {
	args := m.Called(v)

	return []byte{}, args.Error(1)
}

type JsonUnmmarshalMock struct {
	mock.Mock
}

func (m *JsonUnmmarshalMock) Unmarshal(data []byte, v any) error {
	args := m.Called(data, v)

	return args.Error(0)
}

type ReadAllMock struct {
	mock.Mock
}

func (m *ReadAllMock) ReadAll(r io.Reader) ([]byte, error) {
	args := m.Called(r)

	return []byte{}, args.Error(1)
}

type LogDebugMessageMock struct {
	mock.Mock
}

func (m *LogDebugMessageMock) LogDebugMessage(format string, v ...any) {
	m.Called(format, v)
}
