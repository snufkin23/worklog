FLUTTER := fvm flutter
DART := fvm dart

.PHONY: fvm-check setup get gen local generate run analyze test format build-apk clean doctor

fvm-check:
	@command -v fvm >/dev/null || (echo "fvm not found. Install: brew tap leoafarias/fvm && brew install fvm" && exit 1)
	@fvm --version

setup: fvm-check
	fvm install
	$(FLUTTER) pub get
	$(MAKE) generate

get: fvm-check
	$(FLUTTER) pub get

gen: fvm-check
	$(DART) run build_runner build

local: fvm-check
	$(FLUTTER) gen-l10n

generate: fvm-check
	$(MAKE) local
	$(MAKE) gen

run: fvm-check
	$(FLUTTER) run

analyze: fvm-check
	$(FLUTTER) analyze

test: fvm-check
	$(FLUTTER) test

format: fvm-check
	$(DART) format lib test

build-apk: fvm-check
	$(FLUTTER) build apk

clean: fvm-check
	$(FLUTTER) clean

doctor: fvm-check
	$(FLUTTER) doctor