export GOPATH:=$(shell pwd):$(GOPATH)
export GOBIN:=$(shell pwd)/bin

agent:
	@go install src/nagios_agent.go
client:
	@go install src/nagios_client.go
