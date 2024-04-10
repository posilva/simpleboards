.PHONY: up down full fmt plan init clean 

init:
	tflocal init

fmt:
	tflocal fmt 

plan: init
	tflocal plan 

up: fmt 
	tflocal apply -auto-approve 

down: 
	tflocal destroy -auto-approve

full: down up

clean: 
	rm -rf .terraform terraform.tfstate* 
