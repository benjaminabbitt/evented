branch_hash=$(shell git rev-parse --abbrev-ref HEAD | shasum | cut -c1-8)
topdir=$(shell git rev-parse --show-toplevel)