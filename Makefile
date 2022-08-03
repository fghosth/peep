Version := beta
default:
	@echo 'Usage of make: [ build | linux_build | windows_build | clean ]'
build: tempstatik
	@go build -mod=vendor -tags netgo -ldflags "-X main.Version=$(Version) -X 'main.BuildTime=`date`' -X 'main.GoVersion=`go version`'" -o ./dist/peep ./cmd
	@tar czvf peep.tar.gz -C ./dist .
	@shasum -a 256 peep.tar.gz && cp ./dist/peep /usr/local/bin
linux_build: confile
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor  -tags netgo -ldflags "-X main.Version=$(Version) -X 'main.BuildTime=`date`' -X 'main.GoVersion=`go version`'"  -o ./dist/peep ./cmd
windows_build: confile
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -mod=vendor -o ./dist/peep.exe ./cmd
run: build
	@./dist/peep pcreate --module "fghosth.net/reportssss" --path /Users/derek/project/git/goProject/tmp
stop:
	@kill -9 $(ps -ef | grep "./dist/server" |gawk '$0 !~/grep/ {print $2}' |tr -s '\n' ' ')

install: build
	@mv ./dist/server $(GOPATH)/dist/

clean:
	@rm -rf ./dist
confile:
	@mkdir -p ./dist/conf
	#@cp ./conf/config.yaml ./dist/conf
tempstatik:
	@rm -rf ./statik && statik -f -src=template
gitcommtool:
	@brew install node
	@npm install -g cnpm --registry=https://registry.npm.taobao.org
	@cnpm install -g commitizen
	@cnpm install -g conventional-changelog
	@cnpm install -g conventional-changelog-cli
	@cnpm init --yes
	@commitizen init cz-conventional-changelog --save --save-exact
gitlog:
	@conventional-changelog -p angular -i CHANGELOG.md -s


