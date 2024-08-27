package breach

type SequenceConfig struct {
	Description string
	Points      int
	Size        int
}

type Roundconfig struct {
	Buffer    int
	Matrix    int
	Sequences []SequenceConfig
}

var DefaultConfig = []Roundconfig{
	{
		Matrix: 5,
		Buffer: 10,
		Sequences: []SequenceConfig{
			{
				Size:        3,
				Description: "Initialize connection",
				Points:      10,
			},
			{
				Size:        5,
				Description: "Steal credentials",
				Points:      30,
			},
			{
				Size:        6,
				Description: "Extract databases",
				Points:      50,
			},
		},
	},
	{
		Matrix: 7,
		Buffer: 10,
		Sequences: []SequenceConfig{
			{
				Size:        3,
				Description: "Remove logs",
				Points:      10,
			},
			{
				Size:        5,
				Description: "Exploit Netwatch bug",
				Points:      30,
			},
			{
				Size:        10,
				Description: "Burn Netrunner ice",
				Points:      50,
			},
		},
	},
	{
		Matrix: 3,
		Buffer: 5,
		Sequences: []SequenceConfig{
			{
				Size:        3,
				Description: "Build own ice",
				Points:      10,
			},
		},
	},
	{
		Matrix: 5,
		Buffer: 8,
		Sequences: []SequenceConfig{
			{
				Size:        5,
				Description: "Find Mikoshi source code",
				Points:      10,
			},
			{
				Size:        6,
				Description: "Escape Arasaka Netrunners",
				Points:      10,
			},
		},
	},
}
