.PHONY: .FORCE
GO=go
DOT=dot

PROGS = hub \
	 event \
	 stats \
	 agent 

DOC_DIR = ./doc
SRCDIR = ./src

all: $(PROGS) $(DOC_DIR)/arch.png $(DOC_DIR)/fsm.png

$(PROGS):
	$(GO) install $@

$(DOC_DIR)/%.png: $(DOC_DIR)/%.dot
	$(DOT) -Tpng $< -o $@
		
clean:
	rm -rf bin pkg $(DOC_DIR)/arch.png
 
fmt:
	$(GO) fmt $(SRCDIR)/...
