package strgen

import (
	"testing"
	"time"
)

func TestStringGenerator(t *testing.T) {
	testChannel := make(chan string, 2)
	stringGen := New(testChannel)

	err := stringGen.Start()
	if err != nil {
		t.Fatalf("Failed to start StringGenerator: %v", err)
	}

	time.Sleep(3 * time.Second)
	stringGen.Stop()
	close(testChannel)

	channelMessages := []string{}
	for str := range testChannel {
		channelMessages = append(channelMessages, str)
	}

	if len(channelMessages) == 0 {
		t.Fatalf("No messages received")
	}

	for _, str := range channelMessages {
		if len(str) != 10 {
			t.Fatalf("Generated string has incorrect length: %s", str)
		}
	}
}
