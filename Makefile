.PHONY: install-tools
install-tools:
	@./scripts/install-tools.sh

.PHONY: run-app
run-app:
	@./scripts/run.sh  

.PHONY: build
build:
	@./scripts/build.sh 

.PHONY: install-dependencies
install-dependencies:
	@./scripts/install-dependencies.sh 

# .PHONY: integration-test
# integration-test:
# 	@./scripts/test.sh integration
#
# .PHONY: e2e-test
# e2e-test:
# 	@./scripts/test.sh e2e

#.PHONY: load-test
#load-test:
#	@./scripts/test.sh load-test

.PHONY: format
format:
	@./scripts/format.sh 

.PHONY: lint
lint:
	@./scripts/lint.sh 

.PHONY: update-dependencies
update-dependencies:
	@./scripts/update-dependencies.sh

.PHONY: seed-users
seed-users:
	@echo "Seeding test users..."
	@if [ -f .env ]; then \
		if [ -f test_users_credentials.txt ]; then \
			echo ""; \
			echo "⚠️  WARNING: test_users_credentials.txt already exists!"; \
			echo ""; \
			echo "Running this command will:"; \
			echo "  1. Generate NEW random passwords"; \
			echo "  2. Overwrite the credentials file"; \
			echo "  3. Skip users that already exist in the database"; \
			echo ""; \
			echo "If you want to use the NEW passwords, you must first delete existing users."; \
			echo "To delete all test users, run:"; \
			echo "  docker exec -it algorithmia_postgres_dev psql -U algorithmia -d algorithmia -c \"DELETE FROM user_roles; DELETE FROM users WHERE email LIKE '%@algorithmia.com';\""; \
			echo ""; \
			read -p "Do you want to continue? (y/N): " confirm; \
			if [ "$$confirm" != "y" ] && [ "$$confirm" != "Y" ]; then \
				echo "Cancelled."; \
				exit 0; \
			fi; \
		fi; \
		export $$(grep -v '^#' .env | xargs) && go run ./cmd/tools/seed-users/main.go; \
	else \
		echo "Error: .env file not found"; \
		exit 1; \
	fi

.PHONY: reset-users
reset-users:
	@echo "🗑️  Deleting all test users from database..."
	@docker exec -it algorithmia_postgres_dev psql -U algorithmia -d algorithmia -c "DELETE FROM user_roles; DELETE FROM users WHERE email LIKE '%@algorithmia.com';" || (echo "Error: Could not connect to database. Is it running?"; exit 1)
	@echo "✅ All test users deleted"
	@echo ""
	@echo "Now run 'make seed-users' to create fresh users with new passwords"

.PHONY: reseed-users
reseed-users: reset-users seed-users

.PHONY: hash-password
hash-password:
	@echo "Password hasher utility"
	@go run ./cmd/tools/hashpw/main.go
