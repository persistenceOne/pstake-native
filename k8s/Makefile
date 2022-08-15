PROCESS=gaia

apply:
	kubectl apply -k $(PROCESS)

delete:
	kubectl delete -k $(PROCESS)

local-cluster:
	kind create cluster --name "localnet"
	kubectl config use-context kind-localnet

# debug:
    # port-forward
    # init -genesis
    # get seeds - by calling exposer
    # get genesis and paste in .$BINARY folder
    #
