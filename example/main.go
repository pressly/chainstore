package main

import (
	"fmt"
	"time"

	"github.com/nulayer/chainstore"
	"github.com/nulayer/chainstore/boltdb"
	"github.com/nulayer/chainstore/s3"
)

func main() {
	diskStore, err := boltdb.NewStore("/tmp/store.db", "myBucket")
	if err != nil {
		panic(err)
	}
	diskLru, err := chainstore.NewLRUManager(diskStore, 500*1024*1024) // 500MB of working data
	if err != nil {
		panic(err)
	}

	// NOTE: you'll have to supply your own keys in order for this example to work properly
	remoteStore, err := s3.NewStore("myBucket", "access-key", "secret-key")
	if err != nil {
		panic(err)
	}
	dataStore, err := chainstore.New(diskLru, remoteStore) // flows top-down
	if err != nil {
		panic(err)
	}

	//--

	// Save the object in the chain. It will be Put() synchronously into diskStore,
	// the boltdb engine, and then immediately dispatch background Put()'s to the
	// other stores down the chain, in this case S3.
	fmt.Println("Example 1...")
	obj := []byte{1, 2, 3}
	dataStore.Put("k", obj)
	fmt.Println("Put 'k':", obj, "in the chain")

	v, _ := dataStore.Get("k")
	fmt.Println("Grabbing 'k' from the chain:", v) // => [1 2 3]

	// For demonstration, let's grab the key directly from the store instead of
	// through the chain. This is pretty much the same as above, as the chain's Get()
	// stops once it finds the object.
	v, _ = diskStore.Get("k")
	fmt.Println("Grabbing 'k' directly from boltdb:", v) // => [1 2 3]

	// lets pause for a moment and then try to retrieve the value from the s3 store
	time.Sleep(1e9)

	// Grab the object from s3
	v, _ = remoteStore.Get("k")
	fmt.Println("Grabbing 'k' directly from s3:", v) // => [1 2 3]

	// Delete the object from everywhere
	dataStore.Del("k")
	time.Sleep(1e9) // pause for s3 demo
	v, _ = dataStore.Get("k")
	fmt.Println("Deleted 'k' from the chain (all stores). Get(k) returns:", v)

	//--

	// Another interesting behavior of the chain is when doing a Get(), it goes down
	// the entire chain looking for the value, and when found, it will Put() that
	// object back up the chain for subsequent retrievals. Lets see..
	fmt.Println("Example 2...")
	obj = []byte("hope you enjoy")
	dataStore.Put("hi", obj)
	fmt.Println("Put 'hi':", obj, "in the chain")
	time.Sleep(1e9) // lets wait for s3 again with more then enough time

	diskStore.Del("hi")
	v, _ = diskStore.Get("hi")
	fmt.Println("Delete 'hi' from boltdb. diskStore.Get(k) returns:", v)

	v, _ = dataStore.Get("hi")
	fmt.Println("Let's ask the chain for 'hi':", v)
	time.Sleep(1e9) // pause for bg routine to fill our local cache

	// The diskStore now has the value again from remoteStore lower down the chain.
	v, _ = diskStore.Get("hi")
	fmt.Println("Now, let's ask our diskStore again! diskStore.Get(k) returns:", v)

	// Also.. even though it hasn't been demonstrated here, the diskStore will only
	// store a max of 500MB (as defined with diskLru) worth of objects. Give it a shot.
}

/* OUTPUT:

Example 1...
Put 'k': [1 2 3] in the chain
Grabbing 'k' from the chain: [1 2 3]
Grabbing 'k' directly from boltdb: [1 2 3]
Grabbing 'k' directly from s3: [1 2 3]
Deleted 'k' from the chain (all stores). Get(k) returns: []
Example 2...
Put 'hi': [104 111 112 101 32 121 111 117 32 101 110 106 111 121] in the chain
Delete 'hi' from boltdb. diskStore.Get(k) returns: []
Let's ask the chain for 'hi': [104 111 112 101 32 121 111 117 32 101 110 106 111 121]
Now, let's ask our diskStore again! diskStore.Get(k) returns: [104 111 112 101 32 121 111 117 32 101 110 106 111 121]

*/
