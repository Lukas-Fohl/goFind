comp = go
buildPath = ./build
inputFile = ./main.go

build:
	@if [ ! -d $(buildPath) ]; then\
		mkdir $(buildPath);\
	fi
	$(comp) build -o $(buildPath)/main $(inputFile)

run:
	$(MAKE) build
	$(buildPath)/main "/home/lukas/code/td/" "package" -l 1 -c
