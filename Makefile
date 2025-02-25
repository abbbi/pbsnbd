all:
	go build -o nbdkit-pbs-plugin.so -buildmode=c-shared

clean:
	rm -f *.h
	rm -f *.so
