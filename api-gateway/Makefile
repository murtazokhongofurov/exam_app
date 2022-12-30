pull_submodule:
	git submodule update --init --recursive

update_submodule:	
	git submodule update --remote --merge
	
run:
	go run cmd/main.go

create_proto_submodule:
	git submodule add git@gitlab.com:customergroup1/protos.git

run_script:
	./script/gen-proto.sh

swag:
	swag init -g ./api/router.go -o api/docs

