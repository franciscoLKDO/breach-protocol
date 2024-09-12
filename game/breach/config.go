package breach

import "time"

type SequenceConfig struct {
	Description string
	Points      int
	Size        int
}

type Config struct {
	Buffer    int
	Matrix    int
	Timer     time.Duration
	Sequences []SequenceConfig
}

var DefaultConfig = Config{
	Matrix: 5,
	Buffer: 10,
	Timer:  40 * time.Second,
	Sequences: []SequenceConfig{
		{
			Size:        3,
			Description: "Extract sensitive data",
			Points:      10,
		},
		{
			Size:        5,
			Description: "Steal user credentials",
			Points:      30,
		},
		{
			Size:        6,
			Description: "Corrupt database records",
			Points:      50,
		},
	},
}
