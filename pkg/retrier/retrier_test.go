package retrier

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithRetry_NilReturn(t *testing.T) {
	f := func() error {
		return nil
	}

	err := WithRetry(10, f)

	require.NoError(t, err)
}

func TestWithRetry_Unretriable(t *testing.T) {
	f := func() error {
		return fmt.Errorf("Some error: %w", ErrUnretriable)
	}

	err := WithRetry(10, f)

	require.True(t, errors.Is(err, ErrUnretriable))
}

func TestWithRetry_ThirdAttempt(t *testing.T) {
	attempts := 0
	f := func() error {
		attempts++
		if attempts < 3 {
			return errors.New("Some error")
		}
		return nil
	}

	err := WithRetry(3, f)

	require.Equal(t, 3, attempts)
	require.NoError(t, err)
}
