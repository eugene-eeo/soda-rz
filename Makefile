build:
	cd sim && go build && cd ..
	mv sim/sim bin/sim

install: build test
	pipenv install

test: build
	bin/sim | test/sanity
