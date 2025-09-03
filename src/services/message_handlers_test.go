package services

import (
	"testing"

	"github.com/4nonch/echochamber-dc/src/cache"
	"github.com/bwmarrin/discordgo"
)

// TestFormatGuildEmojis covers the function behavior described:
// - replaces valid :name: with the guild emoji replacement when present in cache
// - leaves unknown names unchanged
// - handles adjacent/overlapping colons correctly
// - handles short strings / no-colon strings unchanged
// (yes, AI test. Too boring to write my own.)
func TestFormatGuildEmojis(t *testing.T) {
	// Prepare emoji cache: name "emoji" => id "123".
	// (This matches the examples in the prompt.)
	cache.Emojis.Set([]*discordgo.Emoji{
		{ID: "123", Name: "emoji"},
		{ID: "1", Name: "e"},
		{ID: "999", Name: "smile"},
		{ID: "222", Name: "emoji2"},
	})
	tests := []struct {
		name string
		in   string
		want string
	}{
		// basic
		{"simple replacement", "Some :emoji: here", "Some <:emoji:123> here"},
		{"two adjacent emojis", ":emoji::emoji:", "<:emoji:123><:emoji:123>"},
		{"emoji at start", ":emoji: start", "<:emoji:123> start"},
		{"emoji at end", "end :emoji:", "end <:emoji:123>"},
		{"only colons", "::::", "::::"},
		{"short string (<3) unchanged", ":a", ":a"},
		{"single colon at end", ":emoji", ":emoji"},
		{"unknown emoji left unchanged", "Hello :unknown: world", "Hello :unknown: world"},
		{"case sensitivity - not found", "Look :Emoji: here", "Look :Emoji: here"},

		// more complex / combinations
		{"leading/trailing extra colons", "::emoji::", ":<:emoji:123>:"},
		{"single-char name replacement", ":e:", "<:e:1>"},
		{"adjacent with single-char name", ":emoji::e::emoji2:", "<:emoji:123><:e:1><:emoji2:222>"},
		{"different emojis same string", "hi :smile: and :emoji2:!", "hi <:smile:999> and <:emoji2:222>!"},
		{"overlapping-ish (no separator)", ":emoji:emoji:", "<:emoji:123>emoji:"},
		{"multiple same & unknown in middle", ":emoji::unknown::smile:", "<:emoji:123>:unknown:<:smile:999>"},
		{"sequence of single-char names", ":e::e::e:", "<:e:1><:e:1><:e:1>"},
		{"mixed inline text", "start:emoji:mid:e:end", "start<:emoji:123>mid<:e:1>end"},
		{"embedded words with colons", "a:e:b", "a<:e:1>b"},
		{"multiline / complex (from prompt)", `Lorem:emoji::emoji:Ipsum
TEST:emoji:emoji:TEST

:emoji::emoji::emoji: TEST :::3::::4:::::5
::emoji:::emoji:10:
:emoji::emoji::emoji:23123:emoji::emoji::emoji::::::`,
			// expected result based on implemented algorithm and cache above
			`Lorem<:emoji:123><:emoji:123>Ipsum
TEST<:emoji:123>emoji:TEST

<:emoji:123><:emoji:123><:emoji:123> TEST :::3::::4:::::5
:<:emoji:123>:<:emoji:123>10:
<:emoji:123><:emoji:123><:emoji:123>23123<:emoji:123><:emoji:123><:emoji:123>:::::`},
	}
	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			got := formatGuildEmojis(tc.in)
			if got != tc.want {
				t.Fatalf("formatGuildEmojis mismatch\nname: %s\ninput:\n%q\nwant:\n%q\ngot:\n%q\n",
					tc.name, tc.in, tc.want, got)
			}
		})
	}
}
