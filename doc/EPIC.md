# EPIC Design 

This file introduces EPIC and documents the design rationales and best practices for EPIC-HP.

## Introduction
One important security property of SCION is that end hosts are not allowed to send packets along arbitrary paths, but only along paths that were advertised by all on-path ASes. 
This property is called *path authorization*. 
The ASes only advertise paths that serve their economic interests, and path authorization protects those routing decisions against malicious end hosts.

In SCION, this is implemented by having the ASes create authenticators during beaconing, which end hosts then have to include in the MAC field of their packets. Those MAC fields prove to the ASes, that the packets are allowed to traverse them.
The MAC fields are static, meaning that they are the same for every packet on that path.

However, this implementation of path authorization is insufficient against more sophisticated adversaries that are able to derive MAC fields, for example by brute-forcing them or by observing them in the data plane packets. Because the MACs are static, once derived MACs for one packet can be reused by the adversary to send arbitrarily many other packets (until the authenticators expire).

The EPIC (Every Packet Is Checked) protocol [[1]](#1) solves this problem by introducing per-packet MACs.
Even if an adversary is able to brute-force the MACs for one packet, the MACs cannot be reused to send any other traffic.
Observing MACs from packets in the data plane does also not help an attacker, as he still needs to derive new MACs for every packet he wants to send over an unauthorized path.
This is especially important for hidden paths [[2]](#2). Hidden paths are paths which are not publicly announced, but only communicated to a group of authorized sources. If one of those sources sends traffic on the hidden path using SCION path type packets, an on-path adversary can observe the MACs and reuse them to send traffic on the hidden path himself. This allows the adversary to reach services that were meant to be hidden, or to launch DoS attacks directed towards them.
EPIC precludes such attacks, making hidden paths more secure.

## EPIC-HP Overview
EPIC-HP (EPIC for Hidden Paths) provides the improved security properties of EPIC on the very last inter-AS link of a path. It is meant as a lightweight EPIC version and is specifically designed to better protect hidden paths.

### Assumptions
EPIC-HP makes the following assumptions necessary to provide a meaningful level of security:
- The AS protected by the hidden path is the last AS.
By "last AS" we mean that the beacon which defined the hidden path ends in the local interface of this AS. Or stated differently: the AS behind the hidden path does not forward the beacon defining the hidden path to further downstream/peering ASes. 
- On the interface-pairs (ingress/egress pair) that affect the hidden path, the last two ASes employ one of two different strategies in the data plane:
  1. Only allow EPIC-HP path type traffic. See use case "Highly Secure Hidden Paths" [here](#HighlySecureHiddenPaths). The path type filtering is further explained [here](#PathTypeFiltering).
  2. Prioritize EPIC-HP path type traffic. See use case "DOS-Secure Hidden Paths" [here](#DOSSecureHiddenPaths).
- The last two ASes of the hidden path have a duplicate-suppression system in place. This prohibits DOS attacks based on replayed packets.

### Example
The following figure illustrates those assumptions:
<p align="center">
  <img src="fig/EPIC/path-type-filtering.png" width="630">
</p>

Here, AS 6 is the AS protected by the hidden path (blue lines). The hidden path ends at the local interface (black dot) of AS 6, so AS 6 did not forward the beacon that defines the hidden path further down to AS 7. 
This is however still allowed for SCION path type traffic (green lines): there are SCION paths that enter AS 6 from AS 4. One of the two paths ends in the local interface of AS 6, while the other one is extended further to AS 7. 

In this example, the border routers of AS 6 and AS 5 (the last and penultimate ASes on the hidden path) further implement path type filtering (red dots). For example, AS 5 will block (red "X" in the figure) SCION path type traffic from AS 3 that is destined towards AS 6, as it would affect the hidden path. Instead of blocking non-EPIC-HP path type traffic, ASes could also prioritize EPIC-HP traffic, which would still satisfy the assumptions above.

Of course the ASes can always decide to be more restrictive, for example AS 6 could additionally disallow SCION path type traffic from AS 4, so that its local interface is reachable through the hidden path only.

### SCION Path Type Responses
EPIC-HP path type packets contain the full SCION path type header plus a timestamp and verification fields for the penultimate and last ASes on the path. 

The included SCION path type header allows the destination behind a hidden path to directly respond with SCION path type packets. The destination only has to extract the SCION path type header from the EPIC-HP header and reverse the path. For this, there is a "SCION-Response (S)" flag in the EPIC-HP path type header, which can be set by the source. Setting this flag (S = 1) makes sense if the source AS is not behind a hidden path, or if it is behind a hidden path but also allows SCION path type traffic (but with lower priority than EPIC-HP).

If the SCION-Response flag is not set (S = 0), this means that the source is itself behind a hidden path. The destination will therefore answer with a new EPIC-HP packet, provided it has the necessary authenticators for the hidden path towards the source.

<p align="center">
  <img src="fig/EPIC/SCION-reponse-flag.png" width="600">
</p>

## Procedures

### Control Plane
In the control plane, the ASes do not only append 6 bytes of the hop authenticators to the beacon, but also the remaining 10 bytes (the authenticator is the 16 byte long output of a MAC function).

### Data Plane
The data plane operations for EPIC-HP path type packets are the same as for SCION path type packets, but the source additionally computes two per-packet validation fields for the penultimate and last ASes on the path. 
The last two ASes need to validate the fields accordingly.
A more concise description can be found in the EPIC-HP path type specification.

## <a id="PathTypeFiltering"></a> Path Type Filtering
Network operators should be able to clearly define which kind of traffic (SCION, EPIC-HP, EPIC-SAPV, and other protocols) they want to allow. 
Therefore, for each AS and every interface pair, an AS can be configured with 1-bit flags to allow only certain types of traffic: 

<p align=center>AllowedTraffic(If<sub>1</sub>, If<sub>2</sub>) = (flag<sub>SCION</sub>, flag<sub>EPIC-HP</sub>, flag<sub>COLIBRI</sub>, ...)</p>

The order of the interfaces, (If<sub>1</sub>, If<sub>2</sub>) vs. (If<sub>2</sub>, If<sub>1</sub>), allows to enable and disable different types of traffic depending on the direction. 
To exclusively allow SCION path type traffic (default) between interfaces 'x' and 'y' we would set:

<p align=center>AllowedTraffic(x, y) = (1, 0, 0, ...)</p>

And similarly to only allow EPIC-HP path type traffic:

<p align=center>AllowedTraffic(x, y) = (0, 1, 0, ...)</p>

## Best Practices
There are two main applications for EPIC-HP:

### <a id="HighlySecureHiddenPaths"></a> Highly Secure Hidden Paths
The last and penultimate ASes on the hidden path only allow EPIC-HP traffic on the interface pairs that affect the hidden path. 
With such a setup it is not possible for unauthorized sources to reach the services in the last AS. Therefore, EPIC-HP effectively prevents adversaries from running attacks like denial of service, or attack preparations like scanning the services for vulnerabilities. 

If some host H1 inside an AS with such a setup wants to communicate with a host inside another AS that is also behind a hidden path, both hosts need to have valid authenticators to send traffic over the corresponding hidden paths. The hosts can exclusively communicate using EPIC-HP, and H1 needs to set the SCION-Response flag to zero.

Note that hosts behind a hidden path can send SCION path type packets towards hosts in other ASes, but that those hosts can not send a response back if they do not have the necessary authenticators.

### <a id="DOSSecureHiddenPaths"></a> DoS-Secure Hidden Paths
The last and penultimate ASes on the hidden path allow EPIC-HP and other path types simultaneously, but prioritize traffic using the EPIC-HP path type over the SCION path type.
This means that DOS attacks are not possible, because an adversary is limited to sending low-priority SCION path type packets. However, an adversary can still reach the services behind the hidden path.

In this scenario, the hosts protected by the hidden path can set the SCION-Response flag to one, so the destination will be able to answer with (arbitrariliy many) SCION path type packets.

## References
<a id="1">[1]</a> 
M. Legner, T. Klenze, M. Wyss, C. Sprenger, A. Perrig. (2020)
EPIC: Every Packet Is Checked in the Data Plane of a Path-Aware Internet
Proceedings of the USENIX Security Symposium 
[[Link]](https://netsec.ethz.ch/publications/papers/Legner_Usenix2020_EPIC.pdf)

<a id="2">[2]</a> 
Design Document for the Hidden Path Infrastructure
[[Link]](https://scion.docs.anapaya.net/en/latest/HiddenPaths.html)

