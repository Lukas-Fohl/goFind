comp = go
buildPath = ./build
inputFile = ./main.go
exeName = gfind

build:
	@if [ ! -d $(buildPath) ]; then\
		mkdir $(buildPath);\
	fi
	$(comp) build -o $(buildPath)/$(exeName) $(inputFile)

install:
	$(MAKE) build -B
	@if [ -f /usr/local/bin/$(exeName) ]; then\
		rm /usr/local/bin/$(exeName);\
	fi
	mv $(buildPath)/$(exeName) /usr/local/bin/$(exeName)

