**************************
SCION Header Specification
**************************

This document contains the specification of the SCION packet header.

SCION Header Formats
====================
Header Alignment
----------------
The SCION Header is aligned to 4 bytes.

Common Header
-------------
The Common Header has the following format:
::
     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |Version|      QoS      |                FlowID                 |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |    NextHdr    |    HdrLen     |          PayloadLen           |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |    PathType   |DT |DL |ST |SL |              RSV              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Version
    The version of the SCION Header. Currently, only 0 is supported.
QoS
    8-bit traffic class field. The value of the Traffic Class bits in a received
    packet or fragment might be different from the value sent by the packet's
    source. The current use of the Traffic Class field for Differentiated
    Services and Explicit Congestion Notification is specified in `RFC2474
    <https://tools.ietf.org/html/rfc2474>`_ and `RFC3168
    <https://tools.ietf.org/html/rfc3168>`_
FlowID
    The 20-bit FlowID field is used by a source to
    label sequences of packets to be treated in the network as a single
    flow. It is **mandatory** to be set.
NextHdr
    Field that encodes the type of the first header after the SCION header. This
    can be either a SCION extension or a layer-4 protocol such as TCP or UDP.
    Values of this field respect and extend `IANA’s assigned internet protocol
    numbers <https://perma.cc/FBE8-S2W5>`_.
HdrLen
    Length of the SCION header in bytes (i.e., the sum of the lengths of the
    common header, the address header, and the path header). All SCION header
    fields are aligned to a multiple of 4 bytes. The SCION header length is
    computed as ``HdrLen * 4 bytes``. The 8 bits of the ``HdrLen`` field limit
    the SCION header to a maximum of 1024 bytes.
PayloadLen
    Length of the payload in bytes. The payload includes extension headers and
    the L4 payload. This field is 16 bits long, supporting a maximum payload
    size of 65'535 bytes.
PathType
    The PathType specifies the SCION path type with up to 256 different types.
    The format of each path type is independent of each other. The initially
    proposed SCION path types are SCION (0), OneHopPath (1), EPIC-HP (2),
    EPIC-SAPV (3) and COLIBRI (4).
DT/DL/ST/SL
    DT/ST and DL/SL encode host-address type and host-address length,
    respectively, for destination/ source. The possible host address length
    values are 4 bytes, 8 bytes, 12 bytes and 16 bytes. ST and DT additionally
    specify the type of the address. If some address has a length different from
    the supported values, the next larger size can be used and the address can
    be padded with zeros.
RSV
    These bits are currently reserved for future use.

Address Header
==============
The Address Header has the following format:
::
     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |            DstISD           |                                 |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+                                 +
    |                             DstAS                             |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |            SrcISD           |                                 |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+                                 +
    |                             SrcAS                             |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                    DstHostAddr (variable Len)                 |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                    SrcHostAddr (variable Len)                 |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

DstISD, SrcISD
    16-bit ISD identifier of the destination/source.
DstAS, SrcAS
    48-bit AS identifier of the destination/source.
DstHostAddr, SrcHostAddr
    Variable length host address of the destination/source. The length and type
    is given by the DT/DL/ST/SL flags in the common header.

Path Type: SCION
================
The path type SCION has the following layout:
::
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                          PathMetaHdr                          |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           InfoField                           |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                              ...                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           InfoField                           |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           HopField                            |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           HopField                            |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                              ...                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+`

It consists of a path meta header, up to 3 info fields and up to 64 hop fields.

PathMeta Header
---------------

The PathMeta field is a 4 byte header containing meta information about the
SCION path contained in the path header. It has the following format:
::
     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    | C |  CurrHF   |    RSV    |  Seg0Len  |  Seg1Len  |  Seg2Len  |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

(C)urrINF
    2-bits index (0-based) pointing to the current info field (see offset
    calculations below).
CurrHF
    6-bits index (0-based) pointing to the current hop field (see offset
    calculations below).
Seg{0,1,2}Len
    The number of hop fields in a given segment. :math:`Seg_iLen > 0` implies
    the existence of info field `i`.

Path Offset Calculations
^^^^^^^^^^^^^^^^^^^^^^^^

The number of info fields is implied by :math:`Seg_iLen > 0,\; i \in [0,2]`,
thus :math:`NumINF = N + 1 \: \text{if}\: Seg_NLen > 0, \; N \in [2, 1, 0]`. It
is an error to have :math:`Seg_XLen > 0 \land Seg_YLen == 0, \; 2 \geq X > Y
\geq 0`. If all :math:`Seg_iLen == 0` then this denotes an empty path, which is
only valid for intra-AS communication.

The offsets of the current info field and current hop field (relative to the end
of the address header) are now calculated as

.. math::
    \begin{align}
    \text{InfoFieldOffset} &= 4B + 8B \cdot \text{CurrINF}\\
    \text{HopFieldOffset} &= 4B + 8B \cdot \text{NumINF}  + 12B \cdot
    \text{CurrHF} \end{align}

To check that the current hop field is in the segment of the current
info field, the ``CurrHF`` needs to be compared to the ``SegLen`` fields of the
current and preceding info fields.

This construction allows for up to three info fields, which is the maximum for a
SCION path. Should there ever be a path type with more than three segments, this
would require a new path type to be introduced (which would also allow for a
backwards-compatible upgrade). The advantage of this construction is that all
the offsets can be calculated and validated purely from the path meta header,
which greatly simplifies processing logic.

Info Field
----------
InfoField has the following format:
::
     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |r r r r r r P C|      RSV      |             SegID             |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           Timestamp                           |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

r
    Unused and reserved for future use.
P
    Peering flag. If set to true, then the forwarding path is built as
    a peering path, which requires special processing on the dataplane.
C
    Construction direction flag. If set to true then the hop fields are arranged
    in the direction they have been constructed during beaconing.
RSV
    Unused and reserved for future use.
SegID
    SegID is a updatable field that is required for the MAC-chaining mechanism.
Timestamp
    Timestamp created by the initiator of the corresponding beacon. The
    timestamp is expressed in Unix time, and is encoded as an unsigned integer
    within 4 bytes with 1-second time granularity.  This timestamp enables
    validation of the hop field by verification of the expiration time and MAC.

Hop Field
---------
The Hop Field has the following format:
::
     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |r r r r r r I E|    ExpTime    |           ConsIngress         |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |        ConsEgress             |                               |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+                               +
    |                              MAC                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

r
    Unused and reserved for future use.
I
    ConsIngress Router Alert. If the ConsIngress Router Alert is set, the
    ingress router (in construction direction) will process the L4 payload in
    the packet.
E
    ConsEgress Router Alert. If the ConsEgress Router Alert is set, the egress
    router (in construction direction) will process the L4 payload in the
    packet.

    .. Note::

        A sender cannot rely on multiple routers retrieving and processing the
        payload even if it sets multiple router alert flags. This is entirely
        use case dependent and in the case of `SCMP traceroute` for example the
        router for which the traceroute request is intended will process it (if
        the corresponding router alert flag is set) and reply to the request
        without further forwarding the request along the path. Use cases that
        require multiple routers/hops on the path to process a packet should
        instead rely on a **hop-by-hop extension**.
ExpTime
    Expiry time of a hop field. The field is 1-byte long, thus there are 256
    different values available to express an expiration time. The expiration
    time expressed by the value of this field is relative, and an absolute
    expiration time in seconds is computed in combination with the timestamp
    field (from the corresponding info field) as follows

    .. math::
        Timestamp + (1 + ExpTime) \cdot \frac{24\cdot60\cdot60}{256}

ConsIngress, ConsEgress
    The 16-bits ingress/egress interface IDs in construction direction.
MAC
    6-byte Message Authentication Code to authenticate the hop field. For details on how this MAC is calculated refer to `Hop Field MAC Computation`_.

Hop Field MAC Computation
-------------------------
The MAC in each hop field has two purposes:

#. Authentication of the information contained in the hop field itself, in
   particular ``ExpTime``, ``ConsIngress``, and ``ConsEgress``.
#. Prevention of addition, removal, or reordering hops within a path segment
   created during beaconing.

To that end, MACs are calculated over the relevant fields of a hop field and
additionally (conceptually) chained to other hop fields in the path segment. In
the following, we specify the computation of a hop field MAC.

We write the *i*-th  hop field in a path segment (in construction direction) as

.. math::
    HF_i = \langle  Flags_i || ExpTime_i || InIF_i || EgIF_i || \sigma_i \rangle

:math:`\sigma_i` is the hop field MAC calculate as

.. math::
    \sigma_i = \text{MAC}_{K_i}(TS || ExpTime_i || InIF_i || EgIF_i || \beta_i)

where *TS* is the `Timestamp` and :math:`\beta_i` is the current ``SegID`` of
the info field. :math:`\beta_i` changes at each hop according to the following
rules:

.. math::
    \begin{align}
    \beta_0 &= \text{RND}()\\
    \beta_{i+1} &= \beta_i \oplus \sigma_i[:2]
    \end{align}

Here, :math:`\sigma_i[:2]` is the hop field MAC truncated to 2 bytes and
:math:`\oplus` denotes bitwise XOR.

During beaconing, the initial random value :math:`\beta_0` can be stored in the
info field and all subsequent segment identifiers can be added to the respective
hop entries, i.e., :math:`\beta_{i+1}` can be added to the *i*-th hop entry. On
the data plane, the `SegID` field must contain :math:`\beta_{i+1}/\beta_i` for a
segment in up/down direction before being processed at the *i*th hop (this also
applies to core segments).

Peering Links
^^^^^^^^^^^^^

Peering hop fields can still be "chained" to the AS' standard up/down hop field
via the use of :math:`\beta_{i+1}`:

.. math::
    \begin{align}
    HF^P_i &= \langle  Flags^P_i || ExpTime^P_i || InIF^P_i || EgIF^P_i ||
    \sigma^P_i \rangle\\
    \sigma^P_i &= \text{MAC}_{K_i}(TS || ExpTime^P_i || InIF^P_i || EgIF^P_i || \beta_{i+1})
    \end{align}

Path Calculation
^^^^^^^^^^^^^^^^

**Initialization**

The paths must be initialized correctly for the border routers to verify the hop
fields in the data plane. `SegID` is an updatable field and is initialized based
on the location of sender in relation to path construction.



Initialization cases:

- The non-peering path segment is traversed in construction direction. It starts
  at the `i`-th AS of the full segment discovered in beaconing:

  :math:`SegID := \beta_{i}`

- The peering path segment is traversed in construction direction. It starts at
  the `i`-th AS of the full segment discovered in beaconing:

  :math:`SegID := \beta_{i+1}`

- The path segment is traversed against construction direction. The full segment
  discovered in beaconing has `n` hops:

  :math:`SegID := \beta_{n}`

**AS Traversal Operations**

Each AS on the path verifies the hop fields with the help of the current value
in `SegID`. The operations differ based on the location of the AS on the path.
Each AS has to set the `SegID` correctly for the next AS to verify its hop
field. These operations also have to be done by ASes that deliver the packet
to a local end host to ensure that path can be used in the reverse direction.

Each operation is described form the perspective of AS `i`.

Against construction direction (up, i.e., ConsDir == 0):
   #. `SegID` contains :math:`\beta_{i+1}` at this point.
   #. Compute :math:`\beta'_{i} := SegID \oplus \sigma_i[:2]`
   #. Compute :math:`\sigma'_i` with the formula above by replacing
      :math:`\beta_{i}` with :math:`\beta'_{i}`.
   #. Check that the MAC in the hop field matches :math:`\sigma'_{i}`.
   #. Update `SegID` for the next hop:

      :math:`SegID := \beta'_{i}`
   #. `SegID` now contains :math:`\beta_{i}`.

In construction direction (down, i.e., ConsDir == 1):
   #. `SegID` contains :math:`\beta_{i}` at this point.
   #. Compute :math:`\sigma'_i` with the formula above by replacing
      :math:`\beta_{i}` with `SegID`.
   #. Check that the MAC in the hop field matches :math:`\sigma'_{i}`.
   #. Update `SegID` for the next hop:

      :math:`SegID := SegID \oplus \sigma_i[:2]`
   #. `SegID` now contains :math:`\beta_{i+1}`.

The computation for ASes where a peering link is crossed between path segments
is special cased. A path containing a peering link contains exactly two path
segments, one in construction direction (down) and one against construction
direction (up). On the path segment in construction direction, the peering AS is
the first hop of the segment. Against construction direction (up), the peering
AS is the last hop of the segment.

Against construction direction (up):
   #. `SegID` contains :math:`\beta_{i+1}` at this point.
   #. Compute :math:`{\sigma^P_i}'` with the formula above by replacing
      :math:`\beta_{i+1}` with `SegID`.
   #. Check that the MAC in the hop field matches :math:`{\sigma^P_i}'`.
   #. Do not update `SegID` as it already contains :math:`\beta_{i+1}`.

In construction direction (down):
   #. `SegID` contains :math:`\beta_{i+1}` at this point.
   #. Compute :math:`{\sigma^P_i}'` with the formula above by replacing
      :math:`\beta_{i+1}` with `SegID`.
   #. Check that the MAC in the hop field matches :math:`{\sigma^P_i}'`.
   #. Do not update `SegID` as it already contains :math:`\beta_{i+1}`.

Path Type: OneHopPath
=====================

The OneHopPath path type is a special case of the SCION path type. It is used to
handle communication between two entities from neighboring ASes that do not have
a forwarding path. Currently, it's only used for bootstrapping beaconing between
neighboring ASes.

A OneHopPath has exactly one info field and two hop fields with the speciality
that the second hop field is not known apriori, but is instead created by the
corresponding BR upon processing of the OneHopPath.
::
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           InfoField                           |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           HopField                            |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           HopField                            |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Because of its special structure, no PathMeta header is needed. There is only a
single info field and the appropriate hop field can be processed by a border
router based on the source and destination address, i.e., ``if srcIA == self.IA:
CurrHF := 0`` and ``if dstIA == self.IA: CurrHF := 1``.

.. -------------------------------------------------------------------

Path Type: EPIC-HP
==================
The EPIC-HP (EPIC for Hidden Paths) header provides improved path authorization for the last hop of the path. 
In standard SCION, an attacker that once observed or brute-forced the hop authenticators for some path can use 
them to send arbitrary traffic along this path. EPIC-HP solves this problem on the last hop, which is particularly important for the 
security of hidden paths.

The EPIC-HP header has the following structure:
   - A *PacketTimestamp* field (8 bytes)
   - The path header for the standard SCION Path Type, where one bit of the Path Meta Header is used to indicate whether the sender accepts SCION response packets.
   - A 4-byte *LHVF* (Last Hop Verification Field) 

The EPIC-HP header contains the full SCION header, and also the calculation of the MAC is identical. This allows 
the destination host to directly send back a SCION answer packet to the source by inverting the path.
This is allowed from a security perspective, because the SCION answer packets do not leak information that would allow unauthorized entities to use the hidden path.
To protect the services behind the hidden path from DoS-attacks (only authorized entities should be able to access the services, prevent downgrade to standard SCION), ASes need to be able to configure the border routers such that only certain Path Types are allowed (see configuration_ section). 

::

    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                        PacketTimestamp                        |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                          PathMetaHdr                          |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           InfoField                           |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                              ...                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           InfoField                           |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           HopField                            |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                              ...                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           HopField                            |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                             LHVF                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Path Meta Header
----------------

::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    | C |  CurrHF   |S|   RSV   |  Seg0Len  |  Seg1Len  |  Seg2Len  |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

SCION-Response (S)
  Indicates whether the sender accepts SCION response packets. A sender that is not behind a hidden path can set this flag so that the service knows it has to answer with SCION traffic. A sender that is protected by a hidden path itself does not set this flag, as its AS likely drops standard SCION packets - the service knows that it will have to answer with EPIC-HP instead.

Packet Timestamp
----------------
::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                             TsRel                             |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                             PckId                             |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

TsRel
  A 4-byte timestamp relative to the (segment) Timestamp in the first 
  Info Field. TsRel is calculated by the source host as follows:
 
.. math::
    \begin{align}
        \text{Timestamp}_{\mu s} &= \text{Timestamp [s]} 
            \times 10^6 \\
        \text{Ts} &= \text{current unix timestamp [\mu s]}  \\
        \text{q} &= \left\lceil\left(\frac{24 \times 60 \times 60 
            \times 10^6}{2^{32}}\right)\right\rceil\text{\mu s}
            = \text{21 \mu s}\\
        \text{TsRel} &= \text{max} \left\{0, 
            \frac{\text{Ts - Timestamp}_{\mu s}}
            {\text{q}} -1 \right\} \\
        \textit{Get back the time when} &\textit{the packet 
        was timestamped:} \\
        \text{Ts} &= \text{Timestamp}_{\mu s} + (1 + \text{TsRel}) 
            \times \text{q} 
    \end{align}

TsRel has a precision of :math:`\text{21 \mu s}` and covers at least  
one day (1 day and 63 minutes). When sending packets at high speeds 
(more than one packet every :math:`\text{21 \mu s}`) or when using 
multiple cores, collisions may occur in TsRel. To solve this 
problem, the source further identifies the packet using PckId.

PckId
  A 4-byte identifier that allows to distinguish two packets with 
  the same TsRel. Every source is free to set PckId arbitrarily, but 
  we recommend to use the following structure:

::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |    CoreID     |                  CoreCounter                  |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
        
CoreID
  Unique identifier representing one of the cores of the source host.

CoreCounter
  Current value of the core counter belonging to the core specified 
  by CoreID. Every time a core sends an EPIC packet, it increases 
  its core counter (modular addition by 1).

Note that the Packet Timestamp is at the very beginning of the header, this allows other components (like the replay suppression system) to access it without having to go through any parsing overhead. To achieve an even higher precision of the timestamp, the source is free to allocate additional bits from the PckId to TsRel for this purpose.

Last Hop Validation Field (LHVF)
----------------------------------
::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                             LHVF                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

This 4-byte field contains the Hop Validation Field of the last hop of the last segment. 

EPIC Header Length Calculation
------------------------------
The length of the EPIC Path header is the same as the SCION Path
header plus 8 bytes (Packet Timestamp), and plus 4 bytes for the LHVF.

Procedures
----------
**Control plane:**
The beaconing process is the same as for SCION, but the last AS not 
only adds the 6 bytes of the truncated MAC, but further appends the 
remaining 10 bytes, which together define the 16-byte authenticator 
:math:`{\sigma_{\text{LH}}}` for the last hop (LH). 

**Data plane:**
The source fetches the path, including all the 6-byte short hop 
authenticators and the remaining 10 bytes of the last authenticator, 
from the path server. It copies the short authenticators to the 
corresponding MAC-subfield of the Hop Fields as in standard SCION 
and adds the current Packet Timestamp. 
In addition, it calculates the Last Hop Validation Field as follows:

.. math::    
    \begin{align}
    \text{Origin} &= \text{(SrcISD, SrcAS, SrcHostAddr)} \\
    \text{LHVF} &= \text{MAC}_{\sigma_{\text{LH}}}
        (\text{PacketTimestamp}, 
        \text{Origin}, \text{PayloadLen})~\text{[0:4]}
    \end{align}

The border routers of the on-path ASes validate and forward the 
data plane packets as in standard SCION (recalculate 
:math:`\sigma_{i}` and compare to the MAC field in the packet). In 
addition, the last hop of the last segment recomputes and verifies 
the LHVF field (:math:`\sigma_{\text{LH}} = \sigma_{i}`, where i is 
the last hop). If the verification fails, the packet is dropped.

How to only allow EPIC-HP traffic on a hidden path (and not standard 
SCION packets) is described in the configuration_ section.

.. -------------------------------------------------------------------

Path Type: EPIC-SAPV
====================
The Path Type EPIC-SAPV (EPIC Source Authentication and Path Validation) contains the following parts:
   - An 8-byte Packet Timestamp (same as for EPIC-HP).
   - A slightly modified SCION header.
   - A 16-byte *DVF* (Destination Validation Field).

::

    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                        PacketTimestamp                        |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                          PathMetaHdr                          |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           InfoField                           |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                              ...                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           InfoField                           |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           HopField                            |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                              ...                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           HopField                            |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                              DVF                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

SCION Header Modifications
--------------------------
EPIC-SAPV contains the standard SCION header with the following adaptations:
   - Two reserved bits of the Meta Header are used to indicate the EPIC Version (EV).
   - The size of the MAC (six bytes in standard SCION) inside the Hop Fields is reduced to two bytes, the four bytes of freed space are used for the Hop Validation Field (HVF). 

Path Meta Header
^^^^^^^^^^^^^^^^
::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    | C |  CurrHF   |  RSV  |EV |  Seg0Len  |  Seg1Len  |  Seg2Len  |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

EPIC Version (EV):
   - **EV = 0:** Provides per-packet source authentication: every AS on the path can verify that the packet source is authentic.
   - **EV = 1:** *unused (may be used for path validation in the future)*
   - **EV = 2:** *unused*
   - **EV = 3:** *unused*

Hop Field
^^^^^^^^^
We reduce the size of the MAC field to 2 bytes and assign a 4-byte Hop Validation Field (HVF) to the freed space.
The total size of the Hop Field stays the same (12 bytes).

::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |r r r r r r I E|    ExpTime    |           ConsIngress         |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |        ConsEgress             |              MAC              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                              HVF                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Destination Validation Field (DVF)
----------------------------------
::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                                                               |
    +                                                               +
    |                                                               |
    +                              DVF                              +
    |                                                               |
    +                                                               +
    |                                                               |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

The 16-byte Destination Validation Field is only present if the EPIC 
Version is 1. The DVF contains the MAC calculated by the source host 
to authenticate itself to the destination host.

EPIC-SAPV Header Length Calculation
-----------------------------------
The length of the EPIC Path header is the same as the SCION Path
header plus 8 bytes (Packet Timestamp), and plus 16 bytes for the DVF.

Procedures
----------
**Control plane:**
The beacons have to additionally carry a new 16-byte authenticator 
for each AS on the path. The new EPIC-authenticator is calculated in 
almost the same way as in standard SCION, but with an additional 
constant :math:`C_{\text{EPIC}}` prepended to the MAC input:    

.. math::    
    \begin{align}
    \sigma_i^{\text{EPIC}} &= \text{MAC}_{K_i}(C_{\text{EPIC}} || 
        TsPath || ExpTime_i || InIF_i || EgIF_i || \beta_i^\text{EPIC}) \\
    \sigma_i^{\text{P, EPIC}} &= \text{MAC}_{K_i}(C_{\text{EPIC}} || 
        TsPath || ExpTime_i^P || InIF_i^P || 
        EgIF_i^P || \beta_{i+1}^\text{EPIC}) \\
    where \\
    C_{\text{EPIC}} &= 0x391c
    \end{align}

Note that the input to the MAC function now consists of 15 instead of 
11 bytes, but still requires only one encryption operation for block 
cipher based MACs. Also note that :math:`\beta_{0}^EPIC = \beta_{0}`. 
Because the MAC-calculations are different, :math:`\beta_{i}^EPIC 
\neq \beta_{i}` (for i > 0), however.

**Data plane:**
The source host fetches the path, including all the 16-byte EPIC 
authenticators, from the path server. It also retrieves the 
host-level DRKeys (:math:`K_i^S` between the source host and the 
on-path ASes, and :math:`K_{SD}` between the source host and the 
destination host) from the certificate server. The source then 
calculates the Hop Validation Fields (:math:`V_i`) and the 
Destination Validation Field (:math:`V_{\text{SD}}`):

.. math::    
    \begin{align}
    V_i &= \text{MAC}_{K_i^{\text{S}}}(\text{PacketTimestamp}, 
        \text{Origin}, \sigma_i^{\text{EPIC}}, \text{PayloadLen})
        ~\text{[0:4]} \\
    V_{\text{SD}} &= \text{MAC}_{K_{\text{SD}}}
        (\text{PacketTimestamp, Path, Payload}) \\
    where \\
    \text{Path} &= (\text{TsPath}, \text{Address Header}, 
        HI_1, ..., HI_n)\\
    HI_i &= (\text{ExpTime}_i, \text{ConsIngress}_i, 
        \text{ConsEgress}_i) \\
    \end{align}

Depending on the implementation of the MAC, we also need to prepend 
the length of the input (additional 2 bytes). As in SCION, the source 
writes all the :math:`\sigma_i^{\text{EPIC}}` / 
:math:`\sigma_i^{\text{P, EPIC}}` to the MAC-subfield of the Hop 
Fields, but in this case truncates the MAC to 2 bytes instead of 6 
bytes. The :math:`V_i` are subsequently stored in the HVF-subfield of 
the Hop Fields, and :math:`V_{\text{SD}}` in the DVF field. 
The source host writes the necessary 
:math:`\beta` to the SegID of the Info Fields as in standard SCION.

The border routers perform the same operations as in SCION (see 
"Path Calculation" in the SCION Path Type section), but using 
:math:`\beta_i^\text{EPIC}`, :math:`\sigma_i^\text{EPIC}` and 
:math:`\sigma_i^\text{P, EPIC}` instead of :math:`\beta_i`, 
:math:`\sigma_i` and :math:`\sigma_i^P`. This is possible, because 
the Hop Fields in EPIC-SAPV still contain the two bytes of the MAC 
that are necessary for the chaining of the hops.
In addition, the border routers derive the necessary DRKey 
(:math:`K_i^S`), and recompute and validate the :math:`V_i`.

Upon receiving a packet, the destination host fetches :math:`K_{SD}` 
from its local certificate server, recomputes :math:`V_{\text{SD}}` 
and performs validation by comparing it to the DVF in the packet. 

Details of the EPIC Path Types
==============================
.. _configuration:

Configuration
-------------
Network operators should be able to clearly define which kind of 
traffic (SCION, EPIC-HP, EPIC-SAPV, and other protocols) they want to
allow. 
Therefore, for each AS and every interface pair, an AS can be 
configured with flags to allow only certain types of traffic: 

.. math::    
    \begin{align}
    \text{AllowedTraffic(If_1, If_2)} &= \text{(flag_{SCION}, 
    flag_\text{EPIC-HP}, flag_\text{EPIC-SAPV})} 
    \end{align}

The order of the interfaces defines the direction. This allows to 
specify different types of traffic depending on the direction.
To exclusively allow SCION traffic (default) between interfaces 'x' 
and 'y' we would set:

.. math::    
    \begin{align}
    \text{AllowedTraffic(x, y)} &= \text{(1, 0, 0)} 
    \end{align}

And similarly to only allow EPIC traffic (EPIC-HP and EPIC-SAPV):

.. math::    
    \begin{align}
    \text{AllowedTraffic(x, y)} &= \text{(0, 1, 1)} 
    \end{align}

This affects both the control and the data plane. For example on pure 
SCION paths, beacons do not collect the EPIC authenticators, and 
during forwaring every non-SCION packet gets dropped.
If an AS only wants to allow EPIC traffic, it still uses the normal  
SCION beaconing mechanism, extended with the EPIC authenticators, 
but drops non-EPIC packets in the data plane. 
The beaconing mechanism has to be extended such that the beacon 
contains information on which path types are supported for each AS,
which means that the beacons also contain those flags.

Summary of additional beacon extensions
---------------------------------------
A beacon has to additionally carry the following fields:
  - The remaining 10 bytes of the MAC of the very last hop (for EPIC-HP).

It also contains the following per-AS fields:
  - AllowedD: Indicates which path types are allowed in the down 
    direction.
  - AllowedU: Indicates which path types are allowed in the up 
    direction.
  - The EPIC-SAPV authenticators: 
    :math:`\sigma_i^{\text{EPIC}}` and 
    :math:`\sigma_i^{\text{P, EPIC}}` respectively (16 bytes).
  - The 2-byte beta field: :math:`\beta_{i+1}^\text{EPIC}` (:math:`\beta_{0}` is already stored in the beacon, every AS i :math:`\in \{0, 1, ...\}` adds :math:`\beta_{i+1}^\text{EPIC}`)

Cryptographic Primitives
------------------------

.. _CASA: ./casa.rst

In EPIC, hosts and ASes need to agree on what implementation of the 
MAC they want to use. Different ASes may not necessarily agree on 
one globally fixed MAC algorithm however. EPIC therefore leverages 
CASA_ (Cryptographic Agility for SCION ASes), where each AS promotes 
the supported MAC algorithm in the beacons. This way the border 
routers are still very efficient, as they do only have to support 
the MAC specified by their AS. The source hosts know the required 
MAC algorithm of each on-path AS (this information is fetched 
together with the path) and support all the different MAC algorithms.
Note that the structure of the data plane packets does not need to 
be changed, i.e., there is no field necessary to specify the MAC 
algorithm.
