package handler

import (
	"fmt"
	"net"
	"regexp"
)

func (h *Handler) XaddHandler(conn net.Conn, buffer []byte) {
	fmt.Println("XADD Handler")

	params, err := h.parser.GetParams(buffer)
	if err != nil {
		fmt.Println("Error getting params: ", err)
		return
	}

	if len(params) < 9 {
		fmt.Println("Invalid number of arguments")
		return
	}

	key := string(params[4])
	entryID := string(params[6])
	if entryID == "*" {
		fmt.Println("Entry ID is *")
	}

	// get the key-values pairs
	args, err := h.parser.GetArgs(buffer, 6)
	if err != nil {
		fmt.Println("Error getting args: ", err)
		return
	}

	fmt.Println("args: ", args)

	// create stream
	entry := h.store.NewEntry(entryID, args)

	r, _ := regexp.Compile(`(\d+)-(\d+|\*)|(^\*$)`)
	matches := r.FindStringSubmatch(entryID)
	if matches == nil {
		fmt.Println("Invalid entry ID format")
		_, err = conn.Write(h.parser.WriteError("Invalid entry ID format"))
		if err != nil {
			fmt.Println("Error writing to client: ", err)
		}
		return
	}

	if matches[0] == "*" {
		entryID, err = h.store.SetEntryWithAutoGeneratedID(key, entry)
		if err != nil {
			fmt.Println("Error setting entry with auto generated ID: ", err)
			_, err = conn.Write(h.parser.WriteError(err.Error()))
		}

		_, err = conn.Write(h.parser.WriteString(entryID))
		if err != nil {
			fmt.Println("Error writing to client: ", err)
			return
		}
		fmt.Println("\n\nADDED ENTRY: ", entry, " TO KEY: ", key, "\n\n\n")
		return
	}

	if matches[2] == "*" {
		entryID, err = h.store.SetEntryWithAutoGeneratedSequence(key, matches[1], entry)
		if err != nil {
			fmt.Println("Error setting entry with auto generated sequence: ", err)
			_, err = conn.Write(h.parser.WriteError(err.Error()))
		}

		_, err = conn.Write(h.parser.WriteString(entryID))
		if err != nil {
			fmt.Println("Error writing to client: ", err)
			return
		}
		fmt.Println("\n\nADDED ENTRY: ", entry, " TO KEY: ", key, "\n\n\n")
		return
	}

	if matches[0] == "0-0" {
		fmt.Println("ERR The ID specified in XADD must be greater than 0-0")
		_, err = conn.Write(h.parser.WriteError("ERR The ID specified in XADD must be greater than 0-0"))
		if err != nil {
			fmt.Println("Error writing to client: ", err)
		}
		return
	}

	// validate id
	entryID, err = h.store.SetEntryWithID(key, entryID, entry)
	if err != nil {
		fmt.Println("Error setting entry with ID: ", err)
		_, err = conn.Write(h.parser.WriteError(err.Error()))
		if err != nil {
			fmt.Println("Error writing to client: ", err)
			return
		}
		return
	}

	fmt.Println("\n\nADDED ENTRY: ", entry, " TO KEY: ", key, "\n\n\n")

	_, err = conn.Write(h.parser.WriteString(entryID))
	if err != nil {
		fmt.Println("Error writing to client: ", err)
		return
	}

	return
}
