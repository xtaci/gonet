.PHONY: .FORCE
GO=go
DOT=dot

PROGS = hub \
	 event \
	 stats \
	 agent 

DOC_DIR = ./doc
SRCDIR = ./src

GRAPHS = $(DOC_DIR)/arch.png $(DOC_DIR)/fsm.png

all: $(PROGS) $(GRAPHS)

$(PROGS):
	$(GO) install -gcflags -m $@

$(DOC_DIR)/%.png: $(DOC_DIR)/%.dot
	$(DOT) -Tpng $< -o $@
		
clean:
	rm -rf bin pkg $(GRAPHS)
 
fmt:
	$(GO) fmt $(SRCDIR)/...
