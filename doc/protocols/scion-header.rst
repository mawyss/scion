**************************
SCION Header Specification
**************************

.. _header-specification:

This document contains the specification of the SCION packet header.

SCION Header Formats
====================
Header Alignment
----------------
The SCION Header is aligned to 4 bytes.

.. _header-specification_common-header:

Common Header
-------------
The Common Header has the following format::

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
    Values of this field respect the :ref:`Assigned SCION Protocol Numbers
    <assigned-protocol-numbers>`.
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
    proposed SCION path types are Empty (0), SCION (1), OneHopPath (2), EPIC (3)
    and COLIBRI (4).
DT/DL/ST/SL
    DT/ST and DL/SL encode host-address type and host-address length,
    respectively, for destination/ source. The possible host address length
    values are 4 bytes, 8 bytes, 12 bytes and 16 bytes. ST and DT additionally
    specify the type of the address. If some address has a length different from
    the supported values, the next larger size can be used and the address can
    be padded with zeros.
RSV
    These bits are currently reserved for future use.


.. _header-specification_address-header:

Address Header
==============
The Address Header has the following format::

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
The path type SCION has the following layout::

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
SCION path contained in the path header. It has the following format::

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
InfoField has the following format::

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
The Hop Field has the following format::

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
    6-byte Message Authentication Code to authenticate the hop field. For
    details on how this MAC is calculated refer to :ref:`hop-field-mac-computation`.

.. _hop-field-mac-computation:

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

We write the `i`-th  hop field in a path segment (in construction direction) as

.. math::
    HF_i = \langle  Flags_i || ExpTime_i || InIF_i || EgIF_i || \sigma_i \rangle

:math:`\sigma_i` is the hop field MAC calculated from the following input data::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |               0               |            Beta_i             |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           Timestamp                           |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |       0       |    ExpTime    |          ConsIngress          |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |          ConsEgress           |               0               |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

.. math::
    \sigma_i = \text{MAC}_{K_i}(InputData)

where :math:`\beta_i` is the current ``SegID`` of the info field.
The above input data layout comes from the 8 Bytes of the Info field and the
first 8 Bytes of the Hop field with some fields zeroed out.

:math:`\beta_i` changes at each hop according to the following rules:

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
the data plane, the *SegID* field must contain :math:`\beta_{i+1}/\beta_i` for a
segment in up/down direction before being processed at the *i*-th hop (this also
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
field.

Each operation is described form the perspective of AS `i`.

Against construction direction (up, i.e., ConsDir == 0):
   #. `SegID` contains :math:`\beta_{i+1}` at this point.
   #. Compute :math:`\beta'_{i} := SegID \oplus \sigma_i[:2]`
   #. At the ingress router update `SegID`:

      :math:`SegID := \beta'_{i}`
   #. `SegID` now contains :math:`\beta'_{i}`
   #. Compute :math:`\sigma_i` with the formula above by replacing
      :math:`\beta_{i}` with :math:`SegID`.
   #. Check that the MAC in the hop field matches :math:`\sigma_{i}`. If the
      MAC matches it follows that :math:`\beta'_{i} == \beta_{i}`.

In construction direction (down, i.e., ConsDir == 1):
   #. `SegID` contains :math:`\beta_{i}` at this point.
   #. Compute :math:`\sigma'_i` with the formula above by replacing
      :math:`\beta_{i}` with `SegID`.
   #. Check that the MAC in the hop field matches :math:`\sigma'_{i}`.
   #. At the egress router update `SegID` for the next hop:

      :math:`SegID := SegID \oplus \sigma_i[:2]`
   #. `SegID` now contains :math:`\beta_{i+1}`.

An example of how processing is done in up and down direction is shown in the
illustration below:

.. image:: fig/seg-id-calculation.png

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

Path Type: EmptyPath
====================

Empty path is used to send traffic within the AS. It has no additional fields,
i.e., it consumes 0 bytes on the wire.

Path Type: OneHopPath
=====================

The OneHopPath path type is a special case of the SCION path type. It is used to
handle communication between two entities from neighboring ASes that do not have
a forwarding path. Currently, it's only used for bootstrapping beaconing between
neighboring ASes.

A OneHopPath has exactly one info field and two hop fields with the speciality
that the second hop field is not known a priori, but is instead created by the
corresponding BR upon processing of the OneHopPath::

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

.. _pseudo-header-upper-layer-checksum:

Pseudo Header for Upper-Layer Checksum
======================================

Upper-layer protocols that include the addresses from the SCION header in the
checksum computation should use the following pseudo header:

.. code-block:: text

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
    |                    Upper-Layer Packet Length                  |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                      zero                     |  Next Header  |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

DstISD, SrcISD, DstAS, SrcAS, DstHostAddr, SrcHostAddr
    The values are taken from the SCION Address header.
Upper-Layer Packet Length
    The length of the upper-layer header and data. Some upper-layer protocols
    define headers that carry the length information explicitly (e.g., UDP).
    This information is used as the upper-layer packet length in the pseudo
    header for these protocols. For the remaining protocols, that do not carry
    the length information directly (e.g., SCMP), the value is defined as the
    ``PayloadLen`` from the SCION header, minus the sum of the extension header
    lengths.
Next Header
    The protocol identifier associated with the upper-layer protocol (e.g., 1
    for SCMP, 17 for UDP). This field can differ from the ``NextHdr`` field in
    the SCION header, if extensions are present.

Path Type: EPIC-HP
==================
EPIC-HP (EPIC for Hidden Paths) provides improved path authorization
for the last link of the path. For the SCION path type, an attacker
that once observed or brute-forced the hop authenticators for some
path can use them to send arbitrary traffic along this path. EPIC-HP
solves this problem on the last link, which is particularly
important for the security of hidden paths.

The EPIC-HP header has the following structure:
   - A *PacketTimestamp* field (8 bytes)
   - A 4-byte *PHVF* (Penultimate Hop Validation Field)  and a
     4-byte *LHVF* (Last Hop Validation Field)
   - The complete SCION path type header

::

    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                        PacketTimestamp                        |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                             PHVF                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                             LHVF                              |
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

The EPIC-HP header contains the full SCION path type header. The
calculation of the hop field MAC is identical. This allows the
destination host to directly send back (many) SCION path type answer
packets to the source. This can be done by extracting and reversing
the SCION path type header contained in the EPIC-HP packet.

This is allowed from a security perspective, because the SCION path
type answer packets do not leak information that would allow
unauthorized entities to use the hidden path. In particular, a SCION
path type response packet only contains strictly less information
than the previously received EPIC-HP packet, as the response packet
does not include the PacketTimestamp, the PHVF, and the LHVF.

If the sender is reachable through a hidden path itself, then it is
likely that its AS will not accept SCION path type packets, which
means that the destination can only respond using EPIC-HP traffic.
The destination is responsible to configure or fetch the necessary
EPIC-HP authenticators.

To protect the services behind the hidden link (only authorized
entities should be able to access the services, downgrade to the
SCION path type should be prevented, etc.), ASes need to be able to
configure the border routers such that only certain Path Types are
allowed. This is further described in the accompanying
`EPIC design document`_.

.. _`EPIC design document`: ../EPIC.html

Packet Timestamp
----------------

This timestamp represents the precise time at which a packet was sent.
It is relative to the Timestamp in the first `Info Field`_.
Together with the (ISD, AS, host) triple of the packet source and
the Timestamp in the first Info Field, the Packet Timestamp uniquely
identifies a packet. Unique packet identifiers are a requirement for
replay suppression.
The Packet Timestamp further allows the border router to discard
packets that exceed their lifetime or lie in the future.

::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                             TsRel                             |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                             PckId                             |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

TsRel
  A 4-byte timestamp relative to the (segment) Timestamp in the
  first Info Field. TsRel is calculated by the source host as
  follows:

.. math::
    \begin{align}
        \text{Timestamp}_{\mu s} &= \text{Timestamp [s]}
            \times 10^6 \\
        \text{Ts} &= \text{current unix timestamp [}\mu s\text{]}  \\
        \text{q} &= \left\lceil\left(\frac{24 \times 60 \times 60
            \times 10^6}{2^{32}}\right)\right\rceil\mu s
            = \text{21}\mu s\\
        \text{TsRel} &= \text{max} \left\{0,
            \frac{\text{Ts - Timestamp}_{\mu s}}
            {\text{q}} -1 \right\} \\
        \textit{Get back the time when} &~\textit{the packet
        was timestamped:} \\
        \text{Ts} &= \text{Timestamp}_{\mu s} + (1 + \text{TsRel})
            \times \text{q}
    \end{align}

TsRel has a precision of :math:`21 \mu s` and covers at least
one day (1 day and 63 minutes). When sending packets at high speeds
(more than one packet every :math:`21 \mu s`) or when using
multiple cores, collisions may occur in TsRel. To solve this
problem, the source further identifies the packet using PckId.

PckId
  A 4-byte identifier that allows to distinguish packets with
  the same TsRel. Every source is free to set PckId arbitrarily
  (it only needs to be unique for all packets with the same TsRel),
  but we recommend to use the following structure:

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

Note that the Packet Timestamp is at the very beginning of the
header, this allows other components (like the replay suppression
system) to access it without having to go through any parsing
overhead.

Hop Validation Fields (PHVF and LHVF)
-------------------------------------
::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                             PHVF                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                             LHVF                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Those 4-byte fields are the Hop Validation Fields of the
penultimate and the last hop of the last segment.
They contain the output of a MAC function (truncated to 4 bytes).
Before an EPIC-HP packet is sent, the source computes the MACs and
inserts them into the PHVF and the LHVF.
When the packet arrives at the border router of the penultimate AS,
the border router recomputes and validates the PHVF, and when the
packet arrives at the border router of the last AS on the path, its
border router recomputes and validates the LHVF.

The specification of how the MACs for the Hop Validation Fields are
calculated can be found in the `EPIC Procedures`_ section.

EPIC Header Length Calculation
------------------------------
The length of the EPIC Path header is the same as the SCION Path
header plus 8 bytes (Packet Timestamp), and plus 8 bytes for the
PHVF and LHVF.

.. _EPIC Procedures:

Procedures
----------
**Control plane:**
The beaconing process is the same as for SCION, but the ASes not
only add the 6 bytes of the truncated MAC to the beacon, but further
append the remaining 10 bytes.

**Data plane:**
The source fetches the path, including all the 6-byte short hop
authenticators and the remaining 10 bytes of the authenticators,
from a (hidden) path server. We will refer to the fully assembled 16-byte
authenticators of the penultimate and last hop on the path as
:math:`{\sigma_{\text{PH}}}` for the penultimate hop (PH) and
:math:`{\sigma_{\text{LH}}}` for the last hop (LH), respectively.

The source then copies the short authenticators to the corresponding
MAC-subfield of the Hop Fields as for SCION path type packets and
adds the current Packet Timestamp. In addition, it calculates the
PHVF and LHVF as follows:

.. math::
    \begin{align}
    \text{Origin} &= \text{(SrcISD, SrcAS, SrcHostAddr)} \\
    \text{PHVF} &= \text{MAC}_{\sigma_{\text{PH}}}
        (\text{Flags}, \text{Timestamp}, \text{PacketTimestamp},
        \text{Origin}, \text{PayloadLen})~\text{[0:4]} \\
    \text{LHVF} &= \text{MAC}_{\sigma_{\text{LH}}}
        (\text{Flags}, \text{Timestamp}, \text{PacketTimestamp},
        \text{Origin}, \text{PayloadLen})~\text{[0:4]} \\
    \end{align}

Here, "Timestamp" is the Timestamp from the first `Info Field`_ and
"Flags" is a 1-byte field structured as follows:
::

     0 1 2 3 4 5 6 7 8
    +-+-+-+-+-+-+-+-+-+
    |SL |      0      |
    +-+-+-+-+-+-+-+-+-+

"SL" denotes the source host address length as defined in the
`Common Header`_.
Because the length of the source host address varies based on SL,
also the length of the input to the MAC is dynamic.

The border routers of the on-path ASes validate and forward the
EPIC-HP data plane packets as for SCION path type packets
(recalculate :math:`\sigma_{i}` and compare it to the MAC field in
the packet).

In addition, the penultimate hop of the last segment recomputes and
verifies the PHVF field.
As it has already calculated the 16-byte authenticator
:math:`\sigma_{\text{PH}}` in the previous step, the penultimate hop
only needs to extract the Flags, Timestamp, PacketTimestamp and
Origin fields from the EPIC-HP packet, and the PayloadLen from
the Common Header, which is all the information it
needs to recompute the PHVF. If the verification fails, i.e., the
calculated PHVF is not equal to the PHVF field in the EPIC-HP
packet, the packet is dropped. In the case of an authorized source
(a source that knows the :math:`\sigma_{\text{PH}}` and
:math:`\sigma_{\text{LH}}`), the recomputed PHVF and the PHVF
in the packet will always be equal (assuming the packet has not been
tampered with on the way).

Similarly, the last hop of the last segment recomputes and
verifies the LHVF field. Again, if the verification fails, the
packet is dropped.

How to only allow EPIC-HP traffic on a hidden path (and not SCION
path type packets) is described in the `EPIC design document`_.

Path Type: COLIBRI
==================

The COLIBRI path type is a bit different than the regular SCION in that it has
only one info field::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                        PacketTimestamp                        |
    |                                                               |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                                                               |
    |                                                               |
    |                           InfoField                           |
    |                                                               |
    |                                                               |
    |                                                               |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           HopField                            |
    |                                                               |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                           HopField                            |
    |                                                               |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                              ...                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

The sizes of the packet timestamp, the info field and the individual hop fields
are fixed and the fields always exist, although the number of hop fields
is variable.

Colibri Packet Timestamp
------------------------
::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                             TsRel                             |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                             PckId                             |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Both fields ``TsRel`` and ``PckId`` contain arbitrary data when ``C=1``
(defined in the InfoField).
This is so because these fields are only used for E2E data plane traffic,
which means ``C=0``; thus they only need to be set for ``C=0``.

TsRel
  A 4-byte timestamp relative to the Expiration Tick in the InfoField minus 16
  seconds.
  TsRel is calculated by the source host as follows:

.. math::
    \begin{align}
        \text{Timestamp}_{ns} &= (4\times \text{ExpirationTick} - 16)
            \times 10^9 \\
        \text{Ts} &= \text{current unix timestamp [ns]}  \\
        \text{q} &= \left\lceil\left(\frac{16
            \times 10^9}{2^{32}}\right)\right\rceil\text{ns}
            = \text{4 ns}\\
        \text{TsRel} &= \text{max} \left\{0,
            \frac{\text{Ts - Timestamp}_{ns}}
            {\text{q}} -1 \right\} \\
        \textit{Get back the time when }&\textit{the packet
        was timestamped:} \\
        \text{Ts} &= \text{Timestamp}_{ns} + (1 + \text{TsRel})
            \times \text{q}
    \end{align}

TsRel has a precision of :math:`\text{4 ns}` and covers at least
17 seconds. When sending packets at high speeds
(more than one packet every :math:`\text{4 ns}`) or when using
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
  by CoreID. Every time a core sends a COLIBRI packet, it increases
  its core counter (modular addition by 1).

Note that the Packet Timestamp is at the very beginning of the
header, this allows other components (like the replay suppression
system) to access it without having to go through any parsing
overhead.

Colibri Info Field
------------------
The only info field has the following format::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |C R S r r r r r r r r r|  Ver  |     CurrHF    |    HFCount    |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                                                               |
    |                     Reservation ID Suffix                     |
    |                                                               |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                        Expiration Tick                        |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |      BWCls    |      RLC      |    Original Payload Length    |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

r
    Unused and reserved for future use.

(C)ontrol
    This is a control plane packet. On each border router it will be
    forwarded to the COLIBRI anycast address.
(R)everse
    This packet travels in the reverse direction of the reservation.
(S)egment Reservation
    This is a Segment Reservation Packet.
    If `S` is set, `C` must be set as well. Otherwise the packet is invalid.
    This flag is set every time the Reservation ID is of type Segment ID.
Ver
    The version of this reservation.
CurrHF
    The index of the current HopField.
HFCount
    The number of total HopFields.
Reservation ID Suffix
    Uses 12 bytes. Either an E2E Reservation ID suffix or a
    Segment Reservation ID suffix,
    depending on `S`. If :math:`S=1`, the Segment Reservation ID suffix
    is padded with zeroes until using all 12 bytes. If :math:`S=0`
    the 12 bytes from the E2E Reservation ID suffix are included.
Expiration Tick
    The value represents the "tick" where this packet is no longer valid.
    A tick is four seconds, so :math:`\text{Expiration Time} = 4 \times
    \text{Expiration Tick}` seconds after Unix epoch.
BWCls
    The bandwidth class this reservation has.
RLC
    The Request Latency Class this reservation has.
Original Payload Length
    The field only has a meaning if S = 0.
    If R = 0, then this field is the same as the PayloadLen field in
    the SCION common header. If R = 1 (the packet is a response),
    then the field contains the original payload length, i.e., the
    payload lenght of the preceding COLIBRI packet with R = 0.

Reservation ID Reconstruction
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
The reservation ID is encoded in two parts in the packet header.

- The ASID, which is always the initial part of the ID, is encoded in the
  regular SCION address header, either at the `SrcAS` or the `DstAS` field.
- The suffix is present in the `Reservation ID Suffix` field.

The process of reconstructing the reservation ID is simple. It depends
on the value of ``R`` and ``S``:

.. code-block:: go

    var ASID [6]byte
    var Suffix []byte
    if R == 0 {
        ASID = AddressHeader.SrcAS
    } else {
        ASID = AddressHeader.DstAS
    }
    if S == 0 {
        Suffix = InfoField.IDSuffix
    } else {
        Suffix = InfoField.IDSuffix[:4]
    }
    ReservationID = append(ASID, Suffix)

These steps need only to be carried out by entities that need the
complete reservation ID, which excludes the border router
(which only needs to derive the correct ASID using ``R``).


Hop Field
---------
The Hop Field has the following format::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |        Ingress ID             |         Egress ID             |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                              MAC                              |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Hop fields appear in the forwarding order.


Hop Field MAC Computation
-------------------------
There is an explanation about the rationale of the MAC computation on
:ref:`colibri-mac-computation`.
Here we only detail how to perform the two different MAC computations.
The two different MAC flavors are the *static MAC* and the *per-packet MAC*
(also known as *HVF*).

The `InputData` is common for both types::

     0                   1                   2                   3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                                                               |
    |                        Reservation ID                         |
    |                                                               |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                      Expiration Tick                          |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |      BWCls    |      RLC      |        0      |  Ver  |C|  0  |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |            Ingress            |            Egress             |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                                                               |
    |                 ASID          +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                               |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Most of the fields come from the COLIBRI *InfoField*,
with the exception of *Ingress* and *Egress* that come from the *HopField*,
and ``ASID``, which was used to derive the
full reservation ID. Depending on the value of ``R``, is derived as:

.. code-block:: go

    var ASID [6]byte
    if R == 0 {
        ASID = AddressHeader.SrcAS
    } else {
        ASID = AddressHeader.DstAS
    }

When ``C=1`` we compute the *static MAC*:

.. math::
    \text{MAC}_i^{C=1} \equiv \text{MAC}_{K_i} (\text{InputData})

When ``C=0`` we have :math:`\text{MAC}_{i}^{C=0}` which is also called
:math:`\sigma_i`:

.. math::
    \sigma_i = \text{MAC}_{K_i}(InputData, DT, DL, ST, SL,
      SrcHostAddr, DstHostAddr)

(SrcHostAddr and DstHostAddr are defined in the
:ref:`header-specification_address-header`,
SL, DL, ST and DT defined in the
:ref:`header-specification_common-header`,
both present in every SCION packet).

In the case of ``C=0``, we want to use the :math:`\sigma_i` defined above
to compute the *per-packet MAC*,
also known as HopField Validation Field (*HVF*):

.. math::
    \text{HVF}_i = \text{MAC}_{\sigma_i}(\text{PacketTimestamp}, \text{Original Packet Size})

With:

PacketTimestamp
    The Timestamp described on `Colibri Packet Timestamp`_.
Original Packet Size
    The total size of the packet. It is the sum of the common header, the address header, the
    Colibri header, and the original payload size as indicated in the Info Field (Original
    Payload Length).

The per packet MACs (or *HVFs*) are used only when ``C=0``, which implies
that the S flag is also not set (``S=0``). The computation of
the HVFs for all HopFields happens at the *stamping* service in the source AS,
and they are verified at each transit AS, one HVF per transit AS.


.. _colibri-forwarding-process:

Forwarding Process
------------------
There is a unique way of forwarding a COLIBRI packet, regardless of
whether it is control plane or data plane.
This should simplify the design and implementation of the COLIBRI
part in the border router.
There are, though, slight modifications on the forward process depending
on the ``C`` and ``R`` flags, as is noted below.

The validation process checks that all of the following conditions are true:

- The time derived from the expiration tick is less than the current time.
- The consistency of the flags: if `S` is set, `C` must be set as well.
- HFCount is at least 2, :math:`\text{HFCount} \geq 2`.
- The `CurrHF` is not beyond bounds.
  I.e. :math:`\text{CurrHF} \lt \text{HFCount}`

If the packet is valid, we continue to validate the current Hop Field.
Depending on ``R``, ingress and egress in the packet actually represent
the opposite (for forwarding purposes).
The current hop field is located at
`Offset(COLIBRI_header) + Len(TS) + Len(InfoField) + CurrHF` :math:`\times 8`:

- Its `Ingress ID` field is checked against the actual ingress interface.
- Its MAC is computed according to :ref:`colibri-mac-computation`
  and checked against the `MAC` field. If ``C=0`` the *HVF* is computed and
  checked instead of the *static MAC*.

If the packet is valid:

- If `C = 1`, the packet is delivered to the local COLIBRI anycast address.
- If `C = 0` and this AS is the destination AS (last hop):
  - Check `DestIA` against this IA.
- If `C = 0` and this AS is not the destination:

  - Its `CurrHF` field is incremented by 1 if
    :math:`\text{CurrHF} \lt \sum_{i=0}^2 SegLen_i - 1`.
  - It is forwarded to its `Egress ID` interface.
