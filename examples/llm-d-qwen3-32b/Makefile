.PHONY: check render prepare-router run-direct run-matrixhub cleanup cleanup-model cleanup-router

check:
	@./scripts/common.sh
	@./tests/validate-router-integration.sh

render:
	@./scripts/render.sh

prepare-router:
	@./scripts/prepare-router.sh

run-direct:
	@./scripts/run.sh direct

run-matrixhub:
	@./scripts/run.sh matrixhub

cleanup:
