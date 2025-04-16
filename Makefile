run:
	@docker-compose -f config/compose.yml -p kashyap-site up --build 
	@sleep 3
	# @docker exec ollama-ai bash -c "ollama pull smollm2"
stop:
	@docker-compose -f config/compose.yml -p kashyap-site down --remove-orphans
