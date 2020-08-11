******************
CASA Specification
******************

This document contains the specification of the CASA (Cryptographic 
Agility for SCION ASes) protocol.

Background
==========
In SCION, several protocols like DRKey and EPIC require different 
entities (ASes, hosts) to agree on the same cryptographic 
algorithms like PRFs and MACs.

The simplest solution would be to globally define the algorithm 
(e.g. AES-based CMAC) that is used by such a protocol. If an 
algorithm becomes obsolete, the protocol can be redefined using a 
new Path Type, and another algorithm may be introduced. This 
solution assumes that all the ASes in the whole Internet can agree 
on the same cryptographic algorithms.

However, ASes may want to be more flexible in the choice of the 
algorithms they support. A more sophistic approach would be that 
each AS promotes the subset of algorithms (from some globally 
defined set of algorithms) that it supports, similarly to the 
exchange of cipher suites in TLS handshakes. This means however, 
that the data plane packets need to introduce an additional field 
for each AS, which identifies the selected algorithm. For DRKey and 
EPIC, border routers would need to potentially support many 
different algorithms, which hinders scalability and performance.
Also, border routers are likely less agile in adopting new 
cryptographic algorithms as they often contain dedicated hardware for packet 
processing.


CASA can provide most of the advantages of both those two options, 
while avoiding their drawbacks: it is a tradeoff between simplicity 
and agility.

Specification
=============

Description
-----------
In CASA, every AS chooses exactly **one algorithm** per algorithm 
category. Those categories are globally defined and can either be 
general (PRF, MAC, ...) or protocol-specific (PRF-DRKEY, 
MAC-EPIC-HVF, MAC-EPIC-DVF, ...).

::

    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    | MAC:          | CMAC-AES    |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    | PRF:          | CBC-MAC-AES |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    | MAC-EPIC-HVF: | Poly1305    |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    | MAC-EPIC-DVF: | Poly1305    |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    | PRF-DRKEY:    | CBC-MAC-AES |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    | ...           | ...         |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

The idea behind the general categories is to cover a large set of 
protocols that only need algorithm agreement between different 
entities.
On the other hand, the protocol-specific categories serve the needs 
of protocols that have additional requirements like high 
performance. 
In case an algorithm for a protocol-specific category is announced, 
it has priority over the one specified in the general category. 

Those algorithms are then primarily announced in the beacons, but 
can also be exchanged in protocol-specific channels, for example 
when ASes fetch each others first-level DRKeys.

Note the asymmetry in the CASA design: 
There is in general one side with higher and one with lower 
efficiency requirements. The side that requires higher efficiency 
chooses the algorithm to be used, as it is the one that is primarily 
affected by the algorithm's performance. This mirrors the design of 
DRKey, where one side has to fetch the key (slow), and the other side 
can recompute the key on the fly (fast). 

Beaconing Extensions
--------------------
For each category, the beacon needs to be extended by a field that 
indicates the concrete algorithm implementation that an AS supports. 
With a field size of one byte, there are 256 possible algorithms per 
category.

For a scenario with ten different categories and a path consisting 
of eight ASes, the communication overhead introduced in the beacons 
due to CASA amounts to 80 bytes.

Advantages
----------
With CASA, every AS is free to choose which cryptographic algorithms 
it wants to support, there is no global agreement necessary, 
except for the different categories and their options regarding the 
algorithms. 
Furthermore, CASA introduces only a small overhead in the beacons 
and does not affect the layout of the data plane packets in any way, 
meaning that no additional fields are necessary.
The border routers of the ASes only need to implement the small 
number of algorithms promoted by their AS, which minimizes their 
overhead and prevents performance degradation.

Algorithm Categories
--------------------
The following categories are supported:

MAC
^^^
This category contains the following MAC algorithms: 

PRF
^^^
This category contains the following PRF algorithms: 
