/*
I would implement some sort of memory-intensive algorithm here like Equihash or ETHash.

TODO: implement kind of a client puzzle algos registry with at least one computational-intensive and at least one memory-intensive algos.
Choose randomly the algo type for each incoming TCP connection and send algo type with the puzzle itself to the client.
This algo rotation will make possible service attack even harder to implement as attacker will have to support both approaches
and ASIC[https://en.wikipedia.org/wiki/Application-specific_integrated_circuit] aimed just to memory or computation power will not succeed.
*/

package pow
