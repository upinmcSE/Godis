package core

import "github.com/upinmcSE/godis/internal/data_structure"

var dictStore *data_structure.Dict

func init() {
	dictStore = data_structure.CreateDict()
}
