package projectInit

import (
	"github.com/fghosth/peep/util"
	"github.com/rakyll/statik/fs"
	_ "github.com/fghosth/peep/statik"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	gitIgnore = ".gitignore"
	githookPreCommit = "pre-commit"
	githookCommitMsg = "commit-msg"
	githookGoCommitSh = "go_pre_commit.sh"
	gitGolangciConfig = ".golangci.yml"
	statikFS http.FileSystem
)

func init(){
	var err error
	statikFS, err = fs.New()
	if err != nil {
		log.Println(err)
	}
}

func CreateGitIgnore(name string) error {
	return util.WriteWithIoutil(name+"/"+gitIgnore, GITIGNORE_TMP)
}

func CreateGithookPreCommit(name string)error{
	r,err:=statikFS.Open("/githook_pre_commit.tmpl")
	if err != nil {
		log.Println(err)
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	return util.WriteWithIoutil(name+"/"+githookPreCommit, string(contents))
}

func CreateGithookGoCommit(name string)error{
	r,err:=statikFS.Open("/githook_go_sh.tmpl")
	if err != nil {
		log.Println(err)
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	return util.WriteWithIoutil(name+"/"+githookGoCommitSh, string(contents))
}

func CreateGolangCi(name string)error{
	r,err:=statikFS.Open("/golangci_yml.tmpl")
	if err != nil {
		log.Println(err)
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	return util.WriteWithIoutil(name+"/"+gitGolangciConfig, string(contents))
}

func CreateCommitMsg(name string)error{
	r,err:=statikFS.Open("/commit_msg.tmpl")
	if err != nil {
		log.Println(err)
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	return util.WriteWithIoutil(name+"/"+githookCommitMsg, string(contents))
}