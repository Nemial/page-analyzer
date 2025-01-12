up:
	docker compose --env-file .env -f ./deploy/docker-compose.yml up --build -d --force-recreate
down:
	docker compose --env-file .env -f ./deploy/docker-compose.yml down
run:
	go run page-analyzer