.PHONY: .FORCE
GO=go
DOT=dot

PROGS = hub \
	 event \
	 stats \
	 agent 

DOC_DIR = ./doc
GRAPH = $(DOC_DIR)/arch.png

SRCDIR = ./src

all: $(PROGS) $(GRAPH)

$(PROGS):
	$(GO) install $@

$(GRAPH): $(DOC_DIR)/arch.dot
	$(DOT) -Tpng $< -o $@
		
clean:
	rm -rf bin pkg $(DOC_DIR)/arch.png
 
fmt:
	$(GO) fmt $(SRCDIR)/...
