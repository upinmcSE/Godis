package core

import "strings"

/*
	{
		Cmd: "PING"
		Args: ["hello"]
	}

	{
		Cmd: "SET"
		Args: ["key", "10"]
	}
*/
type Command struct {
	Cmd  string
	Args []string
}

func ParseCmd(data []byte) (*Command, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}

	array := value.([]interface{})
	tokens := make([]string, len(array))
	for i := range tokens {
		tokens[i] = array[i].(string)
	}
	res := &Command{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}
	return res, nil
}
