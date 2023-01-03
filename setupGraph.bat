go get golang.org/x/sys/execabs
go get golang.org/x/text/cases
go get golang.org/x/tools/go/ast/astutil
go get github.com/urfave/cli/v2
go get golang.org/x/mod/module
cd .\graph\
go generate
cd ..
echo .