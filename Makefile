ifneq (${GOBIN},)
OUTDIR ?= ${GOBIN}/
endif

EX_DIR = ${CURDIR}/examples/lily58

all: ${OUTDIR}go-qmk-keymap

${OUTDIR}go-qmk-keymap: $(wildcard *.go) go.mod Makefile
	go build -o ${OUTDIR}

ifneq ($(shell command -v watchexec),)
test:
	watchexec -e go -e json -c -r --debounce 2s -w ${CURDIR} -w ${EX_DIR} "make --no-print-directory all && cat ${EX_DIR}/keymap.c | go-qmk-keymap -workdir ${EX_DIR}"
else
test: all
	cat ${EX_DIR}/keymap.c | go-qmk-keymap -workdir ${EX_DIR}
endif



