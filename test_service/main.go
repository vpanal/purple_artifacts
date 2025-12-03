package main

import (
	"log"
	"time"

	"golang.org/x/sys/windows/svc"
)

type myService struct{}

func (m myService) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {

	s <- svc.Status{State: svc.StartPending}
	s <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	// Bucle principal del servicio
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				s <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			}
		default:
			// Simula actividad
			time.Sleep(2 * time.Second)
		}
	}

	s <- svc.Status{State: svc.StopPending}
	return false, 0
}

func main() {
	isService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("Error comprobando servicio: %v", err)
	}

	if isService {
		svc.Run("GoTestService", myService{})
		return
	}

	// Si se ejecuta fuera del SCM (normal)
	log.Println("Ejecutándose en modo consola (no como servicio)")
	for {
		time.Sleep(2 * time.Second)
	}
}
