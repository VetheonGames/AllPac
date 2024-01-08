package main

import (
	"fmt"
)

func handleUpdateError(updateOption string, err error) {
    if err != nil {
        fmt.Printf("Error occurred during '%s' update: %v\n", updateOption, err)
    } else {
        fmt.Printf("Update '%s' completed successfully.\n", updateOption)
    }
}
