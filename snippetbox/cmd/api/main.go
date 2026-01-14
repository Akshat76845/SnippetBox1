package main
import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)


const version = "1.0.0"

type config struct {
	port int
	env string
}

type application struct {
	config config
	logger *log.Logger
}
func main() {
	var cfg config
}