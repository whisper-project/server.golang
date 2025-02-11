/*
 * Copyright 2024 Daniel C. Brotsky. All rights reserved.
 * All the copyrighted work in this repository is licensed under the
 * GNU Affero General Public License v3, reproduced in the LICENSE file.
 */

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/whisper-project/server.golang/platform"

	"github.com/whisper-project/server.golang/legacy/client"
	"github.com/whisper-project/server.golang/legacy/conversation"
	"github.com/whisper-project/server.golang/legacy/profile"
)

func saveObjects(what platform.ObjectMap) {
	var saved int
	var t string
	var as []any
	for t, as = range what {
		if len(as) > 0 {
			switch t {
			case "profiles":
				saved += saveTypedObjects(t, as, &profile.UserProfile{})
			case "clients":
				saved += saveTypedObjects(t, as, &client.Data{})
			case "conversations":
				saved += saveTypedObjects(t, as, &conversation.Data{})
			case "states":
				saved += saveTypedObjects(t, as, &conversation.State{})
			default:
				_, _ = fmt.Fprintf(os.Stderr, "Skipping objects of unknown type: %s", t)
			}
		}
	}
	if saved != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Saved %d objects.\n", saved)
	}
}

func saveTypedObjects[T platform.StructPointer](name string, oa []any, e T) int {
	var saved int
	if len(oa) >= 10 {
		_, _ = fmt.Fprintf(os.Stderr, "Starting to save %s...", name)
	}
	for _, o := range oa {
		s, err := e.Downgrade(o)
		if err != nil {
			panic(err)
		}
		if err = platform.SaveFields(context.Background(), s); err != nil {
			panic(err)
		}
		if saved++; saved%10 == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "\nSaved %d %s...", saved, name)
		}
	}
	if len(oa) >= 10 {
		_, _ = fmt.Fprintf(os.Stderr, "\n")
	}
	if saved != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Saved %d %s.\n", saved, name)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "Saved 1 %s.\n", name[0:len(name)-1])
	}
	return saved
}

// dumpObjectsToPath serializes the entire map to the given filepath
// A path of "-" means use the standard input. Otherwise, if the path does
// not have a JSON extension, one is added.
func dumpObjectsToPath(what platform.ObjectMap, where string) {
	if where == "-" {
		if err := platform.DumpObjectsToStream(what, os.Stdout); err != nil {
			panic(err)
		}
	} else {
		if !strings.HasSuffix(strings.ToLower(where), ".json") {
			where = where + ".json"
		}
		if err := platform.DumpObjectsToPath(what, where); err != nil {
			panic(err)
		}
		fmt.Printf("Objects dumped to %q\n", where)
	}
}

// loadObjectsFromPath loads the objects dumped to the given filepath
// A path of "-" means use the standard input. Otherwise, if the path does
// not have a JSON extension, one is added.
func loadObjectsFromPath(where string) platform.ObjectMap {
	var som platform.StoredObjectMap
	var err error
	if where == "-" {
		som, err = platform.LoadObjectsFromStream(os.Stdin)
	} else {
		if !strings.HasSuffix(strings.ToLower(where), ".json") {
			where = where + ".json"
		}
		som, err = platform.LoadObjectsFromPath(where)
	}
	if err != nil {
		panic(err)
	}
	return loadObjectsFromStorage(som)
}

func loadObjectsFromStorage(som platform.StoredObjectMap) platform.ObjectMap {
	result := make(platform.ObjectMap)
	var err error
	result["profiles"], err = platform.UnmarshalStoredObjects(profile.UserProfile{}, som["profiles"])
	if err != nil {
		panic(err)
	}
	result["clients"], err = platform.UnmarshalStoredObjects(client.Data{}, som["clients"])
	if err != nil {
		panic(err)
	}
	result["conversations"], err = platform.UnmarshalStoredObjects(conversation.Data{}, som["conversations"])
	if err != nil {
		panic(err)
	}
	result["states"], err = platform.UnmarshalStoredObjects(conversation.State{}, som["states"])
	if err != nil {
		panic(err)
	}
	return result
}
