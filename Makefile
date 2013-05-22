GO=go

PROGS = hub \
	 gate \
	 event

all: $(PROGS)

$(PROGS):
	$(GO) install $@


clean:
	rm -rf bin pkg
