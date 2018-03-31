all: build 

build: 
	test -e vendor || glide up
	go install 

