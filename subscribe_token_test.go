package mqtt

import "testing"

// Run with go test -race
func TestSubscribeToken_Result(t *testing.T) {
	s := SubscribeToken{
		baseToken: newBaseToken(),
		subs:      []string{"sdas"},
		subResult: map[string]byte{"sd": 0x1},
		messageID: 0,
	}
	go func() {
		s.Result()["sd"] = 0x2
	}()
	go func() {
		s.Result()["sd"] = 0x3
	}()
	if s.Result()["sd"] != 0x1 && s.Result()["sd"] != 0x2 && s.Result()["sd"] != 0x3 {
		t.Fatal("asd")
	}
}
