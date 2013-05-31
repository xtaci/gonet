字符串相关的函数：

foo_bar <-> FooBar 

字符串hash函数：

http://www.isthe.com/chongo/tech/comp/fnv/#

http://en.wikipedia.org/wiki/Fowler_Noll_Vo_hash

FNV-1a hash

The FNV-1a hash differs from the FNV-1 hash by only the order in which the multiply and XOR is performed:[7]      
<pre>
   hash = FNV_offset_basis       
   for each octet_of_data to be hashed      
      hash = hash XOR octet_of_data       
      hash = hash × FNV_prime       
   return hash        
</pre>
