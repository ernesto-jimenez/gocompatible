package importers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunsLocally(t *testing.T) {
	i := &Local{}
	list, err := i.List("github.com/stretchr/testify", true)
	require.NoError(t, err)
	require.NotEqual(t, 0, len(list))
}
