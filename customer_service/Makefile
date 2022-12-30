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

create_migrate:
	migrate create -ext sql -dir migrations -seq create_customers_table

migrate_up:
	migrate -source file://migrations -database postgres://murtazokhon:12345@database-1.c9lxq3r1itbt.us-east-1.rds.amazonaws.com:5432/customerdata?sslmode=disable up

migrate_down:
	migrate -path migrations/ -database postgres://murtazokhon:12345@database-1.c9lxq3r1itbt.us-east-1.rds.amazonaws.com:5432/customerdata?sslmode=disable down

migrate_force:
	migrate -path migrations/ -database postgres://murtazokhon:12345@database-1.c9lxq3r1itbt.us-east-1.rds.amazonaws.com:5432/customerdata force 1


# migrate_up:
# 	migrate -source file://migrations -database postgres://developer:2002@localhost:5432/customer?sslmode=disable up

# migrate_down:
# 	migrate -path migrations/ -database postgres://developer:2002@localhost:5432/customer?sslmode=disable down

# migrate_force:
# 	migrate -path migrations/ -database postgres://developer:2002@localhost:5432/customer force 1
