PROCESS=gaia
NODES=3

apply:
	kubectl apply -k $(PROCESS)

delete:
	kubectl delete -k $(PROCESS)

delete-pvc:
	for (( i=0; i<$(NODES); i++ ));\
	do kubectl delete pvc \
	$(PROCESS)-validator-pv-storage-$(PROCESS)-validator-$$i --namespace dev-native;\
	done;

# debug:
    # port-forward
    # init -genesis
    # get seeds - by calling exposer
    # get genesis and paste in .$BINARY folder
    #
