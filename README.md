# Bridger 

<img src="bridge.svg" width="20%" > 

Bridger is an http service that bridges http requests, allowing clients to transfer information between each other.

```
[Client A] <-> [Bridger] <-> [Client b]
```

The two clients must know about each other a head of time. The clients are bridged together by specifying unique tokens that idenify a pre approved session. Bridging sessions may be reused by the two clients until the bridging tokens expire.


## Unidirectional Data Transfer

One client sends data and another client receive data.

```
[Client A] [Bridger] [Client b]
    |         |         |
    +-ConnB-->|         |      # 1
    |         |<--ConnA-+      # 2
    |         |         |
    +-File1-->|         +      # 3
    |         +-File1-->|     
    |         |         |   
    -         -         -
```

1. _Client A_ connects to the bridger, ready to send _File1_ to _Client B_. _Client A_ waits for _Client B_ to connect to the bridger.
2. _Client B_ connects to the bridger.
3. _Client A_ is signaled to begin data transfer to _Client B_.



## Bidirectional Data Transfer

Two clients send and receive data at the same time.

```
[Client A] [Bridger] [Client b]
    |         |         |
    +-ConnB-->|         |      # 1
    |         |<--ConnA-+      # 2
    |         |         |
    +-File1-->|<--File2-+      # 3
    |<--File2-+-File1-->|     
    |         |         |   
    -         -         -
```

1. _Client A_ connects to the bridger, ready to send _File1_ to _Client B_. _Client A_ waits for _Client B_ to connect to the bridger.
2. _Client B_ connects to the bridger, ready to send _File2_ to _Client A_.
3. _Client A_ and _Client B_ are signaled to begin data transfer to eachother.


