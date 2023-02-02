all: test-coverage

test-coverage:
	go test -race -coverprofile=profile.out -covermode=atomic ./...
