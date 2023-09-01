.PHONY: run stop

run: 
	docker-compose up -d --build

stop: 
	docker-compose stop

down: 
	docker-compose down
	
login:
	docker-compose up -d --build login

consumer:
	docker-compose up -d --build consumer

sender:
	docker-compose up -d --build sender

tcpSmoke:
	docker-compose up -d --build tcpsmoke

postgres:
	docker-compose up -d --build db-sql
	
front:
	docker-compose up -d --build nginx