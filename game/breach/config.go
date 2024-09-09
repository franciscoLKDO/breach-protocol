package breach

type SequenceConfig struct {
	Description string
	Points      int
	Size        int
}

type Config struct {
	Buffer    int
	Matrix    int
	Sequences []SequenceConfig
}

var DefaultConfig = Config{
	Matrix: 5,
	Buffer: 10,
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

// {
// 	Matrix: 7,
// 	Buffer: 10,
// 	Sequences: []SequenceConfig{
// 		{
// 			Size:        3,
// 			Description: "Avoid firewall detection",
// 			Points:      10,
// 		},
// 		{
// 			Size:        5,
// 			Description: "Decrypt encrypted files",
// 			Points:      30,
// 		},
// 		{
// 			Size:        10,
// 			Description: "Burn Netrunner ice",
// 			Points:      70,
// 		},
// 	},
// },
// {
// 	Matrix: 3,
// 	Buffer: 5,
// 	Sequences: []SequenceConfig{
// 		{
// 			Size:        3,
// 			Description: "Lock out intruders",
// 			Points:      30,
// 		},
// 	},
// },
// {
// 	Matrix: 5,
// 	Buffer: 8,
// 	Sequences: []SequenceConfig{
// 		{
// 			Size:        5,
// 			Description: "Find Mikoshi source code",
// 			Points:      30,
// 		},
// 		{
// 			Size:        6,
// 			Description: "Escape Arasaka Netrunners",
// 			Points:      50,
// 		},
// 	},
// },
