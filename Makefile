#database
migrate_create:
	migrate create -ext sql -dir .\database\migrations $(name)
migrate_up:
	migrate -database "mysql://root@tcp(127.0.0.1:3306)/payment_gateway_lite?charset=utf8mb4&parseTime=True&loc=Local" -path database/migrations up $(filter-out $@, $(MAKECMDGOALS))
migrate_down:
	migrate -database "mysql://root@tcp(127.0.0.1:3306)/payment_gateway_lite?charset=utf8mb4&parseTime=True&loc=Local" -path database/migrations down $(filter-out $@, $(MAKECMDGOALS))
migrate_force:
	migrate -database "mysql://root@tcp(127.0.0.1:3306)/payment_gateway_lite?charset=utf8mb4&parseTime=True&loc=Local" -path database/migrations force $(version)
migrate_version:
	migrate -database "mysql://root@tcp(127.0.0.1:3306)/payment_gateway_lite?charset=utf8mb4&parseTime=True&loc=Local" -path database/migrations version
