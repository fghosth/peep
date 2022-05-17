package projectInit

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/fghosth/peep/util"
	"io/ioutil"
	"log"
)

var (
	dokcerFile                  = "Dockerfile"
	entrypointFile              = "entrypoint.sh"
	readmeFile                  = "README.md"
	makeFile                    = "Makefile"
	dsshFile                    = "dssh.go"
	dsshTestFile                = "dssh_test.go"
	apolloFile                  = "apollo.go"
	apolloTestFile              = "apollo_test.go"
	filexFile                   = "filex.go"
	filexTestFile               = "filex_test.go"
	apolloOpenApiClientFile     = "client.go"
	apolloOpenApiClientTestFile = "client_test.go"
	apolloOpenApiAPiFile        = "openapi.go"
)

func CreateApolloOpenApiFile(name string) error {
	r, err := statikFS.Open("/apolloOpenApi/client.tmpl")
	if err != nil {
		log.Println(err)
	}
	r_test, err := statikFS.Open("/apolloOpenApi/client_test.tmpl")
	if err != nil {
		log.Println(err)
	}
	rapi_test, err := statikFS.Open("/apolloOpenApi/openapi.tmpl")
	if err != nil {
		log.Println(err)
	}
	defer func() {
		r.Close()
		r_test.Close()
		rapi_test.Close()
	}()
	contents, err := ioutil.ReadAll(r)
	err = util.WriteWithIoutil(name+"/"+apolloOpenApiClientFile, string(contents))
	if err != nil {
		return err
	}
	contents, err = ioutil.ReadAll(r_test)
	err = util.WriteWithIoutil(name+"/"+apolloOpenApiClientTestFile, string(contents))
	if err != nil {
		return err
	}
	contents, err = ioutil.ReadAll(rapi_test)
	return util.WriteWithIoutil(name+"/"+apolloOpenApiAPiFile, string(contents))
}

func CreateFilexFile(name string) error {
	r, err := statikFS.Open("/filex.tmpl")
	if err != nil {
		log.Println(err)
	}
	r_test, err := statikFS.Open("/filex_test.tmpl")
	if err != nil {
		log.Println(err)
	}
	defer func() {
		r.Close()
		r_test.Close()
	}()
	contents, err := ioutil.ReadAll(r)
	err = util.WriteWithIoutil(name+"/"+filexFile, string(contents))
	if err != nil {
		return err
	}
	contents, err = ioutil.ReadAll(r_test)
	return util.WriteWithIoutil(name+"/"+filexTestFile, string(contents))
}

func CreateDsshFile(name string) error {
	r, err := statikFS.Open("/dssh.tmpl")
	if err != nil {
		log.Println(err)
	}
	r_test, err := statikFS.Open("/dssh_test.tmpl")
	if err != nil {
		log.Println(err)
	}
	defer func() {
		r.Close()
		r_test.Close()
	}()
	contents, err := ioutil.ReadAll(r)
	err = util.WriteWithIoutil(name+"/"+dsshFile, string(contents))
	if err != nil {
		return err
	}
	contents, err = ioutil.ReadAll(r_test)
	return util.WriteWithIoutil(name+"/"+dsshTestFile, string(contents))
}

func CreateApolloFile(name string, module string) error {
	ctx := map[string]interface{}{
		"module": module,
	}

	r, err := statikFS.Open("/apollo.tmpl")
	if err != nil {
		log.Println(err)
	}
	r_test, err := statikFS.Open("/apollo_test.tmpl")
	if err != nil {
		log.Println(err)
	}
	defer func() {
		r.Close()
		r_test.Close()
	}()
	contents, err := ioutil.ReadAll(r)
	confstr, err := raymond.Render(string(contents), ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = util.WriteWithIoutil(name+"/"+apolloFile, string(confstr))
	if err != nil {
		return err
	}
	contents, err = ioutil.ReadAll(r_test)
	confstr, err = raymond.Render(string(contents), ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return util.WriteWithIoutil(name+"/"+apolloTestFile, string(confstr))
}

func CreateDockerfile(name string) error {
	r, err := statikFS.Open("/dockerfile.tmpl")
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

func CreateMakefile(name string, module, path string) error {
	r, err := statikFS.Open("/makefile.tmpl")
	if err != nil {
		log.Println(err)
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	ctx := map[string]interface{}{
		"module": module,
		"path":   path,
	}
	confstr, err := raymond.Render(string(contents), ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return util.WriteWithIoutil(name+"/"+makeFile, confstr)
}
func CreateReadme(name string) error {
	return util.WriteWithIoutil(name+"/"+readmeFile, README_TMP)
}
