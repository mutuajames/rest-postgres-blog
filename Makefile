# Stack
.PHONY:	stop
stop:
	@docker compose -f stack.yml down -v

.PHONY:	prod
prod:
	@docker compose -f stack.yml down -v
	@docker compose -f stack.yml up --build

.PHONY: dev
dev:
	@docker compose -f stack.yml down -v
	@docker compose -f stack.yml -f stack.dev.yml up