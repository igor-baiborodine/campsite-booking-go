###############################################################################
# Variables
###############################################################################

###############################################################################
# Target: pprof-get-campsites
###############################################################################
.PHONY: pprof-get-campsites
pprof-get-campsites:
	@curl --output ./tests/perf/get-campsites-profile.pprof "http://localhost:6060/debug/pprof/profile?seconds=10" \
		& SERVER_ADDR=localhost:8085 go test -bench BenchmarkGetCampsites ./tests/perf \
		& wait
	@echo "completed profiling GetCampsite endpoint"

###############################################################################
# Target: pprof-get-campsites-data
###############################################################################
.PHONY: pprof-get-campsites-data
pprof-get-campsites-data:
	@go tool pprof ./tests/perf/get-campsites-profile.pprof
