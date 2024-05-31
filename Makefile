ifneq (${GOBIN},)
OUTDIR ?= ${GOBIN}/
endif

EX_DIR = ${CURDIR}/examples/lily58

all: ${OUTDIR}go-qmk-keymap

${OUTDIR}go-qmk-keymap: main.go go.mod Makefile
	go build -o ${OUTDIR}

ifneq ($(shell command -v watchexec),)
test:
	watchexec -e go -c -r -w ${CURDIR} "make --no-print-directory all && cat ${EX_DIR}/keymap.c | go-qmk-keymap -workdir ${EX_DIR}"
else
test: all
	cat ${EX_DIR}/keymap.c | go-qmk-keymap -workdir ${EX_DIR}
endif



