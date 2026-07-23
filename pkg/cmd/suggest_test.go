package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestJaroDistance(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 1.0, jaroDistance("", ""))
	assert.Equal(t, 0.0, jaroDistance("", "tweet"))
	assert.Equal(t, 0.0, jaroDistance("tweet", ""))
	assert.Equal(t, 0.0, jaroDistance("abc", "xyz"))
	assert.Equal(t, 1.0, jaroDistance("tweet", "tweet"))
	assert.InDelta(t, 0.944, jaroDistance("martha", "marhta"), 0.01)
	assert.Greater(t, jaroDistance("tweets", "tweet"), 0.8)
}

func TestJaroWinkler(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 0.0, jaroWinkler("abc", "xyz"))
	assert.Equal(t, 1.0, jaroWinkler("tweet", "tweet"))
	assert.Greater(t, jaroWinkler("followers", "follower"), jaroDistance("followers", "follower"))
	assert.InDelta(t, 0.667, jaroWinkler("prefix", "presto"), 0.01)
}

func TestSuggestCommand(t *testing.T) {
	t.Parallel()

	commands := []*cli.Command{
		{Name: "search", Aliases: []string{"find"}},
		{Name: "followers"},
	}
	assert.Equal(t, "Did you mean 'search'?", suggestCommand(commands, "serch"))
	assert.Equal(t, "Did you mean 'followers'?", suggestCommand(commands, "folowers"))
	assert.Equal(t, "Did you mean ''?", suggestCommand(nil, "unknown"))
}
