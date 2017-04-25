#current_dir := `pwd`
OS := $(shell uname)
ifeq ($(OS),$(filter $(OS), Darwin Linux))
    files := $(shell find . -name '*.thrift')
    rmcmd =-rm -rf `find . -type d -name '*-remote'`
    cleancmd =rm -f `find . -name '*.go'`;rm -rf `find . -type d -empty`;rm -rf `find . -type d -empty`
else
    files := $(wildcard *.thrift */*.thrift)
    rmcmd =-rm -rf $(basename $<)/*-remote
    cleancmd =-rm -rf $(dirs)
endif
objects := $(foreach file, $(files), $(basename $(file))/constants.go)
dirs := $(basename $(files))
pkg_prefix := ""
all : $(objects)

$(objects) : %/constants.go : %.thrift
	thrift -gen go:package_prefix=$(pkg_prefix) -r -out . -I . $<
	$(rmcmd)

.PHONY : clean
clean :
	$(cleancmd)
