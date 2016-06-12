
test:
	basedir=`pwd` ; export GOPATH=$$GOPATH:$$basedir ;echo $$GOPATH;cd src/unit_test && go test
