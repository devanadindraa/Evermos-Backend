include .env
export $(shell sed 's/=.*//' .env)

##################
### Migrations ###
##################

create-migration:
	migrate create -ext sql -dir ./migrations -seq $(NAME)

migrate:
	migrate -path ./migrations -database "mysql://$(BACKEND_DATABASE_USERNAME):$(BACKEND_DATABASE_PASSWORD)@tcp($(BACKEND_DATABASE_HOST):$(BACKEND_DATABASE_PORT))/$(BACKEND_DATABASE_NAME)?multiStatements=true" up

migrate-down:
	migrate -path ./migrations -database "mysql://$(BACKEND_DATABASE_USERNAME):$(BACKEND_DATABASE_PASSWORD)@tcp($(BACKEND_DATABASE_HOST):$(BACKEND_DATABASE_PORT))/$(BACKEND_DATABASE_NAME)?multiStatements=true" down