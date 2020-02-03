package projectInit

import (
	"io/ioutil"
	"log"
	"github.com/fghosth/peep/util"
)

var (
	dokcerFile     = "Dockerfile"
	entrypointFile = "entrypoint.sh"
	readmeFile     = "README.md"
	makeFile       = "Makefile"
)

func CreateDockerfile(name string) error {
	r,err:=statikFS.Open("/dockerfile.tmpl")
	if err != nil {
		log.Println(err)
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	return util.WriteWithIoutil(name+"/"+dokcerFile, string(contents))
}

func CreateEntrypoint(name string) error {
	return util.WriteWithIoutil(name+"/"+entrypointFile, ENTRYPOINT_TMP)
}
func CreateMakefile(name string) error {
	r,err:=statikFS.Open("/makefile.tmpl")
	if err != nil {
		log.Println(err)
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	return util.WriteWithIoutil(name+"/"+makeFile, string(contents))
}
func CreateReadme(name string) error {
	return util.WriteWithIoutil(name+"/"+readmeFile, README_TMP)
}
