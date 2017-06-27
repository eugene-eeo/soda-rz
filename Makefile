build:
	cd sim && go build && cd ..
	mv sim/sim bin/sim

install: build
	pipenv install
