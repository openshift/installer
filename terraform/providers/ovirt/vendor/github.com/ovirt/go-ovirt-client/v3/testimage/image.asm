; This is a minimalistic assembler program that prints "Hello oVirt!"
; and then stops the CPU. It is intended to be in the boot sector of
; a disk image for testing purposes.
;
; You can compile this program by running "nasm image.asm". This will
; create a file called "image" which is the raw image. This is committed
; for convenience.

ORG 0x7C00 ; Starting address of the boot loader
BITS 16 ; Start program in 16 bit mode

MOV SI, text ; Move pointer to the text label into the SI register
CALL printText
JMP halt

; printText prints a null-terminated string from the address passed in SI using the BIOS facility
printText:
    CLI ; Clear interrupts
    MOV AH, 0x0E ; Set BIOS / INT 10 printing facility to text output
.printChar:
    LODSB ; Load byte from the address in SI into AL and advance SI by one
    CMP AL, 0 ; Check if AL is 0
    JE .printReturn ; If yes, jump to the return
    INT 0x10 ; Trigger BIOS print method
    JMP .printChar ; Repeat for next byte
.printReturn:
    RET ; Return from the printText function

halt:
    HLT ; Halt CPU

text:
    DB "Hello oVirt!", 0 ; Embed data into binary, zero-terminated.

TIMES 510 - ($ - $$) DB 0 ; Fill up 510 bytes
DW 0xAA55 ; Write magic bytes for boot loader
