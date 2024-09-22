###############################################################################
# Variables
###############################################################################
SERVER_ADDR = localhost:8085
DEBUG_PPROF_URL = http://localhost:6060/debug/pprof/profile?seconds=10

###############################################################################
# Target: pprof-get-campsites
###############################################################################
.PHONY: pprof-get-campsites
pprof-get-campsites:
	@curl --output ./tests/perf/get-campsites-profile.pprof "$(DEBUG_PPROF_URL)" \
		& SERVER_ADDR=$(SERVER_ADDR) go test -bench BenchmarkGetCampsites ./tests/perf \
		& wait
	@echo "completed profiling GetCampsite endpoint"

###############################################################################
# Target: pprof-get-campsites-data
###############################################################################
.PHONY: pprof-get-campsites-data
pprof-get-campsites-data:
	@go tool pprof ./tests/perf/get-campsites-profile.pprof

###############################################################################
# Target: pprof-get-vacant-dates
###############################################################################
.PHONY: pprof-get-vacant-dates
pprof-get-vacant-dates:
	@curl --output ./tests/perf/get-vacant-dates-profile.pprof "$(DEBUG_PPROF_URL)" \
		& SERVER_ADDR=$(SERVER_ADDR) go test -bench BenchmarkGetVacantDates ./tests/perf \
		& wait
	@echo "completed profiling GetVacantDates endpoint"

###############################################################################
# Target: pprof-get-vacant-dates-data
###############################################################################
.PHONY: pprof-get-vacant-dates-data
pprof-get-vacant-dates-data:
	@go tool pprof ./tests/perf/get-vacant-dates-profile.pprof
