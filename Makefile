.PHONY: build
build:
	@# In GitHub Actions you can't name a repo secret `GITHUB_TOKEN`.
	@# Hence I've had to rename that to `API_TOKEN_GITHUB` here.
	@# As the code looks for this in the environment.
	op signin --account my; \
	API_TOKEN_GITHUB=$$(op read "op://Programming/jct6fke2i6guaj3p7jtrunqkpa/credential") go run ./cmd/archiver/main.go
