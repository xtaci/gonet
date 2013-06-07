.PHONY: .FORCE
GO=go
DOT=dot

PROGS = hub \
	 event \
	 stats \
	 agent 

GRAPH = arch.png

SRCDIR = ./src

all: $(PROGS) $(GRAPH)

$(PROGS):
	$(GO) install $@

$(GRAPH): arch.dot
	$(DOT) -Tpng $< -o $@
		
clean:
	rm -rf bin pkg arch.png
 
fmt:
	$(GO) fmt $(SRCDIR)/...
