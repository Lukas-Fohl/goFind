comp = go
buildPath = ./build
inputFile = ./main.go
exeName = gfind

build:
	@if [ ! -d $(buildPath) ]; then\
		mkdir $(buildPath);\
	fi
	$(comp) build -o $(buildPath)/$(exeName) $(inputFile)

run:
	$(MAKE) build
	$(buildPath)/gfind "/home/lukas/code/td/" "package" -l 1 -c

install:
	$(MAKE) build -B
	@if [ -f /usr/bin/$(exeName) ]; then\
		rm /usr/bin/$(exeName);\
	fi
	mv $(buildPath)/$(exeName) /usr/bin/$(exeName)
