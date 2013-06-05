.PHONY: .FORCE
GO=go

PROGS = hub \
	 event \
	 stats \
	 agent 

SRCDIR = ./src

all: $(PROGS)

$(PROGS):
	$(GO) install -race $@


clean:
	rm -rf bin pkg
 
fmt:
	$(GO) fmt $(SRCDIR)/...
