package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	luajson "github.com/layeh/gopher-json"
	"github.com/yuin/gopher-lua"
	"os"
)

func main() {
	L := lua.NewState()
	defer L.Close()
	luajson.Preload(L)

	err := L.DoFile("TestConfig.lua")
	if err != nil {
		fmt.Println("DoFile:", err.Error())
		return
	}
	v := L.GetGlobal("TestConfig")
	buff, err := luajson.Encode(v)
	if err != nil {
		fmt.Println("Encode:", err.Error())
		return
	}
	file, err := os.Create("TestConfig.json")
	if err != nil {
		fmt.Println("Create:", err.Error())
		return
	}
	defer file.Close()
	writer := bytes.NewBuffer([]byte{})
	json.Indent(writer, buff, "", "\t")
	file.Write(writer.Bytes())
}
