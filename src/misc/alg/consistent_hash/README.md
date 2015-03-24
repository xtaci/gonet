Consistent hashing is a special kind of hashing such that when a hash table is resized and consistent hashing is used, only K/n keys need to be remapped on average, where K is the number of keys, and n is the number of slots. In contrast, in most traditional hash tables, a change in the number of array slots causes nearly all keys to be remapped.

Consistent hashing achieves some of the same goals as Rendezvous hashing (also called HRW Hashing). The two techniques use different algorithms, and were devised independently and contemporaneously.
