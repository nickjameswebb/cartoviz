.DEFAULT_GOAL: install

.PHONY: install
install:
	go mod tidy
	go install cmd/kubectl-cartoviz.go

.PHONY: smoke
smoke: install
	kubectl cartoviz
