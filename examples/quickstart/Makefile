.PHONY: bundle
bundle:
	docker run -it --rm \
	--user=$(id -u):$(id -g) \
	-v $$(PWD):/work -v $$(PWD)/../policy:/policy \
	openpolicyagent/opa:edge-static \
	build /policy -o /work/bundle.tar.gz
