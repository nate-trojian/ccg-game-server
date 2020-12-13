FROM gitpod/workspace-full

# Install new version of golang

### Go ###
ENV GO_VERSION=1.15.6 \
    GOPATH=$HOME/go-packages \
    GOROOT=$HOME/go \
    GO111MODULE=on
ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH
RUN curl -fsSL https://storage.googleapis.com/golang/go$GO_VERSION.linux-amd64.tar.gz | tar -xzv \
    && go get -u -v \
        github.com/acroca/go-symbols \
        github.com/cweill/gotests/... \
        github.com/davidrjenni/reftools/cmd/fillstruct \
        github.com/fatih/gomodifytags \
        github.com/haya14busa/goplay/cmd/goplay \
        github.com/josharian/impl \
        github.com/nsf/gocode \
        github.com/ramya-rao-a/go-outline \
        github.com/rogpeppe/godef \
        github.com/uudashr/gopkgs/cmd/gopkgs \
        github.com/zmb3/gogetdoc \
        golang.org/x/lint/golint \
        golang.org/x/tools/cmd/godoc \
        golang.org/x/tools/cmd/gorename \
        golang.org/x/tools/cmd/guru \
        sourcegraph.com/sqs/goreturns