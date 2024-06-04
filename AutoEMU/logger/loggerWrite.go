package logger

import (
	"fmt"
	"log"
	"os"
)

func Config[T any](msg T, terminal bool) {
	log.Println(msg)
	if terminal {
		fmt.Printf("%s%v\n", ConfigMessage, msg)
	}
}

func Success[T any](msg T, terminal bool) {
	log.Println(msg)
	if terminal {
		fmt.Printf("%s%v\n", SuccessMessage, msg)
	}
}

func Info[T any](msg T, terminal bool) {
	log.Println(msg)
	if terminal {
		fmt.Printf("%s%v\n", InfoMessage, msg)
	}
}

func Warn[T any](msg T, terminal bool) {
	log.Println(msg)
	if terminal {
		fmt.Printf("%s%v\n", warnMessage, msg)
	}
}

func Error[T any](msg T, terminal bool) {
	log.Println(msg)
	if terminal {
		fmt.Printf("%s%v\n", ErrorMessage, msg)
	}
}

func Fatal[T any](msg T, terminal bool) {
	log.Println(msg)
	if terminal {
		fmt.Printf("%s%v\n", FatalMessage, msg)
	}
	os.Exit(1)
}

func Write[T any](msg T, terminal bool) {
	log.Println(msg)
	if terminal {
		fmt.Println(msg)
	}
}
