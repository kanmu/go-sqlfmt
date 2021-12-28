package scanner

import (
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer/internal/reader"
	"github.com/stretchr/testify/require"
)

func TestScanner(t *testing.T) {
	t.Parallel()

	const (
		src          = "abcdefgàéèç"
		maxLookAhead = 3
	)
	scanner := NewRuneScanner(src, WithReaderOptions(reader.WithLookAhead(maxLookAhead)))
	runes := make([]rune, 0, len(src))

	t.Run("Read eof error", func(t *testing.T) {
		for _, expectedRune := range src {
			runes = append(runes, expectedRune)

			ch, err := scanner.Read()
			require.NoError(t, err)
			require.Equalf(t, expectedRune, ch, "expected %c got %c", expectedRune, ch)
		}

		ch, err := scanner.Read()
		require.NoError(t, err)
		require.Equal(t, EOFRune, ch)
	})

	t.Run("Unread backward error", func(t *testing.T) {
		require.NotPanics(t, scanner.Unread)
		lowerBound := len(runes) - maxLookAhead

		for i := len(runes) - 1; i >= lowerBound; i-- {
			ch, err := scanner.Read()
			require.NoError(t, err)

			expectedRune := runes[i]
			require.Equalf(t, expectedRune, ch, "expected rune [%d] to  be %c, got %c", i, expectedRune, ch)

			require.NotPanics(t, scanner.Unread)

			if i > lowerBound {
				require.NotPanics(t, func() {
					scanner.Rewind(1)
				})
			}
		}

		require.Panics(t, scanner.Unread)
	})
}

func TestScannerInvalidRune(t *testing.T) {
	t.Parallel()

	const (
		src = "\xF0\xA4\xAD"
	)
	scanner := NewRuneScanner(src)
	_, err := scanner.Read()
	require.ErrorIs(t, err, ErrInvalidRune)
}
