package handlers

import (
	"testing"
)

func TestRefineScoreWithRules(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Simple Arabic", "20金币", 20},
		{"Arabic with Space", "20 金币", 20},
		{"Simple Chinese", "二十金币", 20},
		{"Chinese Ten", "十金币", 10},
		{"Chinese One Hundred", "一百金币", 100},
		{"Chinese Complex", "二十五金币", 25},
		{"Context Arabic", "跑步5公里 1小时 20金币", 20},
		{"Context Chinese", "背单词30分钟 十积分", 10},
		{"Points Keyword", "10积分", 10},
		{"Coins Keyword", "5 coins", 5},
		{"Progress Keyword", "进步数量5", 0}, // Regex expects number BEFORE keyword
		{"Progress Keyword 2", "5进步数量", 5},
		{"No Score", "跑步5公里", 0},
		{"Duration Confusion", "30分钟", 0},
		{"Distance Confusion", "5公里", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := refineScoreWithRules(tt.input)
			if got != tt.expected {
				t.Errorf("refineScoreWithRules(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseChineseNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"一", 1},
		{"十", 10},
		{"十一", 11},
		{"二十", 20},
		{"二十一", 21},
		{"一百", 100},
		{"一百零一", 101},
		{"一百一十", 110},
		{"三百", 300},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseChineseNumber(tt.input)
			if got != tt.expected {
				t.Errorf("parseChineseNumber(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}
