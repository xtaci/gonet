.PHONY: .FORCE
GO=go

PROGS = hub \
	 gate \
	 event

SRCDIR = ./src
FILES = ${shell find $(SRCDIR) |grep '\.go'}

all: $(PROGS)

$(PROGS):
	$(GO) install $@


clean:
	rm -rf bin pkg
 
fmt:$(FILES)

.FORCE:

$(FILES): .FORCE
	$(GO) fmt $@
