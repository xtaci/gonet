.PHONY: .FORCE
GO=go
DOT=dot
GOYACC=$(GO) tool yacc
NEXBIN=./nex

PROGS = hub \
	 event \
	 stats \
	 agent 

DOC_DIR = ./doc
SRCDIR = ./src
INSPECTDIR = $(SRCDIR)/inspect

NEX=nex.go

GRAPHS = $(DOC_DIR)/arch.png $(DOC_DIR)/fsm.png

INSPECT = $(SRCDIR)/inspect/inspect.go

all: $(NEXBIN) $(INSPECT) $(PROGS) $(GRAPHS)

$(NEXBIN): $(NEX)
	$(GO) build $<

$(PROGS):
	$(GO) install $@

$(DOC_DIR)/%.png: $(DOC_DIR)/%.dot
	$(DOT) -Tpng $< -o $@

$(INSPECT): $(INSPECTDIR)/inspect.nex $(INSPECTDIR)/inspect.y
	$(NEXBIN) $(INSPECTDIR)/inspect.nex
	$(GOYACC) -o $(INSPECTDIR)/inspect.y.go $(INSPECTDIR)/inspect.y 
	$(GO) fmt $(INSPECTDIR)/...
		
clean:
	rm -rf bin pkg $(NEXBIN)
 
fmt:
	$(GO) fmt $(SRCDIR)/...
