package reader

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuneReader(t *testing.T) {
	t.Parallel()

	const src = "abcdefgàéèç"

	reader := NewRewindRuneReader(src, WithLookAhead(len(src)))
	runes := make([]rune, 0, len(src))

	t.Run("ReadRune forward", func(t *testing.T) {
		for _, expectedRune := range src {
			runes = append(runes, expectedRune)

			ch, _, err := reader.ReadRune()
			require.NoError(t, err)
			require.Equalf(t, expectedRune, ch, "expected %c got %c", expectedRune, ch)
		}
	})

	t.Run("ReadRune backwards", func(t *testing.T) {
		require.NoError(t, reader.UnreadRune())

		for i := len(runes) - 1; i >= 0; i-- {
			ch, _, err := reader.ReadRune()
			require.NoError(t, err)

			expectedRune := runes[i]
			require.Equalf(t, expectedRune, ch, "expected %c got %c", expectedRune, ch)

			err = reader.UnreadRune()
			require.NoError(t, err)

			if i > 0 {
				require.NoError(t, reader.UnreadRune())
			}
		}

		require.Error(t, reader.UnreadRune())
	})
}

func TestRuneReaderErrors(t *testing.T) {
	t.Parallel()

	const (
		src          = "abcdefgàéèç"
		maxLookAhead = 3
	)

	reader := NewRewindRuneReader(src, WithLookAhead(maxLookAhead))
	runes := make([]rune, 0, len(src))

	t.Run("ReadRune eof error", func(t *testing.T) {
		for _, expectedRune := range src {
			runes = append(runes, expectedRune)

			ch, _, err := reader.ReadRune()
			require.NoError(t, err)
			require.Equalf(t, expectedRune, ch, "expected %c got %c", expectedRune, ch)
		}

		_, _, err := reader.ReadRune()
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("ReadRune backward error", func(t *testing.T) {
		require.NoError(t, reader.UnreadRune())
		lowerBound := len(runes) - maxLookAhead

		for i := len(runes) - 1; i >= lowerBound; i-- {
			ch, _, err := reader.ReadRune()
			require.NoError(t, err)

			expectedRune := runes[i]
			require.Equalf(t, expectedRune, ch, "expected rune [%d] to  be %c, got %c", i, expectedRune, ch)

			err = reader.UnreadRune()
			require.NoError(t, err)

			if i > lowerBound {
				require.NoError(t, reader.UnreadRune())
			}
		}

		require.Error(t, reader.UnreadRune())
	})
}
