.PHONY: clear-swag swag-init

clear-swag:
	@echo "Clearing swag files..."
	del /Q docs\* 2>nul || exit 0

swag-init: clear-swag
	@echo "Generating swag files..."
	swag init -g main.go --parseDependency --parseInternal --parseDepth 1