package breach

import "time"

type SequenceConfig struct {
	Description string
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
		},
		{
			Size:        5,
			Description: "Steal user credentials",
		},
		{
			Size:        6,
			Description: "Corrupt database records",
		},
	},
}
