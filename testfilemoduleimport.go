package main

import (
	"bytes"
	"fmt"

	"mig.ninja/mig/modules"
	_ "mig.ninja/mig/modules/file"
)

func main() {
	run := modules.Available["file"].NewRun()
	args := make([]string, 0)
	args = append(args, "-path", "/tmp")
	args = append(args, "-name", "^testfile$")
	args = append(args, "-maxdepth", "3")
	args = append(args, "-content", "somestring")
	param, err := run.(modules.HasParamsParser).ParamsParser(args)

	buf, err := modules.MakeMessage(modules.MsgClassParameters, param)
	if err != nil {
		panic(err)
	}
	rdr := bytes.NewReader(buf)

	res := run.Run(rdr)
	fmt.Printf("%s\n", res)
}
