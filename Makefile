GO=go

PROGS = hub \
	 gate \
	 cooldown

all: $(PROGS)

$(PROGS):
	$(GO) install $@
