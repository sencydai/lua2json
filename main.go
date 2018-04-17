package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	luajson "github.com/layeh/gopher-json"
	"github.com/yuin/gopher-lua"
	"io/ioutil"
	"os"
	"strings"
)

type gConfig struct {
	LuaFile  string
	JsonFile string
}

func main() {
	defer func() {
		var s string
		fmt.Scanln(&s)
	}()

	config := gConfig{}
	if data, err := ioutil.ReadFile("config.json"); err != nil {
		fmt.Println(err.Error())
		return
	} else if err = json.Unmarshal(data, &config); err != nil {
		fmt.Println(err.Error())
		return
	}

	L := lua.NewState()
	defer L.Close()
	luajson.Preload(L)

	if err := os.MkdirAll(config.JsonFile, os.ModeDir); err != nil {
		fmt.Println(err.Error())
		return
	}

	if files, err := ioutil.ReadDir(config.LuaFile); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			name := file.Name()
			if !strings.HasSuffix(name, ".lua") {
				continue
			}
			err := L.DoFile(fmt.Sprintf("%s/%s", config.LuaFile, name))
			if err != nil {
				fmt.Printf("load file %s error: %s", name, err.Error())
				continue
			}
			v := L.GetGlobal(strings.TrimSuffix(name, ".lua"))
			buff, err := luajson.Encode(v)
			if err != nil {
				fmt.Printf("parse file %s error: %s", name, err.Error())
				continue
			}
			file, err := os.Create(fmt.Sprintf("%s/%s.json", config.JsonFile, strings.TrimSuffix(name, ".lua")))
			if err != nil {
				fmt.Printf("create json file %s error: %s", name, err.Error())
				continue
			}
			writer := bytes.NewBuffer([]byte{})
			json.Indent(writer, buff, "", "\t")
			file.Write(writer.Bytes())
			file.Close()
		}
	}

	fmt.Println("转换完成...")
}
