.PHONY: .FORCE
GO=go
DOT=dot
GOYACC=$(GO) tool yacc

PROGS = hub \
	 stats \
	 agent 

DOC_DIR = ./doc
SRCDIR = ./src
INSPECTDIR = $(SRCDIR)/inspect

NEX=nex.go

GRAPHS = $(DOC_DIR)/arch.png $(DOC_DIR)/fsm.png

all: $(PROGS) $(GRAPHS)

$(PROGS):
	$(GO) install $@

$(DOC_DIR)/%.png: $(DOC_DIR)/%.dot
	$(DOT) -Tpng $< -o $@

clean:
	rm -rf bin pkg $(NEXBIN)
 
fmt:
	$(GO) fmt $(SRCDIR)/...
