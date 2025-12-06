package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: dll_loader <ruta_dll> <export> [param1] [param2]...")
		fmt.Println("\nEjemplos que funcionan:")
		fmt.Println("  dll_loader kernel32.dll GetTickCount")
		fmt.Println("  dll_loader user32.dll MessageBeep 0")
		fmt.Println("  dll_loader kernel32.dll Sleep 2000")
		return
	}

	dllPath := os.Args[1]
	exportName := os.Args[2]

	dll, err := syscall.LoadDLL(dllPath)
	if err != nil {
		fmt.Println("Error al cargar la DLL:", err)
		return
	}
	defer dll.Release()

	proc, err := dll.FindProc(exportName)
	if err != nil {
		fmt.Println("Error: export no encontrado:", err)
		return
	}

	// Convertir argumentos a uintptr
	var params []uintptr
	for i := 3; i < len(os.Args); i++ {
		val, _ := strconv.ParseUint(os.Args[i], 10, 64)
		params = append(params, uintptr(val))
	}

	// Llamar función
	var ret uintptr
	var callErr error

	switch len(params) {
	case 0:
		ret, _, callErr = proc.Call()
	case 1:
		ret, _, callErr = proc.Call(params[0])
	case 2:
		ret, _, callErr = proc.Call(params[0], params[1])
	case 3:
		ret, _, callErr = proc.Call(params[0], params[1], params[2])
	case 4:
		ret, _, callErr = proc.Call(params[0], params[1], params[2], params[3])
	default:
		ret, _, callErr = proc.Call(params...)
	}

	// Mostrar resultado
	if callErr.Error() == "Solicitud no compatible." {
		fmt.Printf("Ret: %v (0x%x)\n", ret, ret)
	} else {
		fmt.Printf("Ret: %v (0x%x)  Error: %v\n", ret, ret, callErr)
	}
}
