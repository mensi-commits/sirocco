package types

type Command struct {
	Action string `json:"action"`
	Data   []byte `json:"data"`
}

type ShardConfig struct {
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}