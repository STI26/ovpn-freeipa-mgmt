VERSION_BACK  = v1.1.0
VERSION_FRONT = v1.1.0
API_URL       = /api
BUILD_DATE    = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

all: build_back build_front

.PHONY: all

build_back:
	docker build \
		--file Dockerfile.back \
		--rm --compress \
		--build-arg VERSION="$(VERSION_BACK)" \
		--build-arg DATE="$(BUILD_DATE)" \
		--tag imagelist/ovpn_freeipa_mgmt:$(VERSION_BACK)-api \
		--tag imagelist/ovpn_freeipa_mgmt:latest-api \
		.

build_front:
	docker build \
		--file Dockerfile.front \
		--rm --compress \
		--build-arg VERSION="$(VERSION_FRONT)" \
		--build-arg API_URL="$(API_URL)" \
		--build-arg DATE="$(BUILD_DATE)" \
		--tag imagelist/ovpn_freeipa_mgmt:$(VERSION_FRONT)-ui \
		--tag imagelist/ovpn_freeipa_mgmt:latest-ui \
		.

push:
	docker push imagelist/ovpn_freeipa_mgmt:$(VERSION_BACK)-api
	docker push imagelist/ovpn_freeipa_mgmt:latest-api
	docker push imagelist/ovpn_freeipa_mgmt:$(VERSION_FRONT)-ui
	docker push imagelist/ovpn_freeipa_mgmt:latest-ui
