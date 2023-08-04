.PHONY: run stop

run: 
	docker-compose up -d --build

stop: 
	docker-compose stop

down: 
	docker-compose down

consumer:
	docker-compose up -d --build consumer

sender:
	docker-compose up -d --build sender

tcpSmoke:
	docker-compose up -d --build tcpsmoke

postgres:
	docker-compose up -d --build db_sql
	
front:
	docker-compose up -d --build my-react-app