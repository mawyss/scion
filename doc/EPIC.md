# EPIC Design 

This file documents the design rationales and best practices for EPIC-HP.

## Introduction

Explain EPIC, cite paper
Explain EPIC-HP (main application: hidden paths, only last link protected)

## EPIC-HP Overview


<p align="center">
  <img src="fig/EPIC/SCION-reponse-flag.png" width="600">
</p>

Add image: Really secure hidden path vs. prioritized path

## Procedures

### Control Plane
Add remaining 10 bytes of authenticator

### Data Plane


## Configuration
Network operators should be able to clearly define which kind of traffic (SCION, EPIC-HP, EPIC-SAPV, and other protocols) they want to allow. 
Therefore, for each AS and every interface pair, an AS can be configured with flags to allow only certain types of traffic: 

AllowedTraffic(If_1, If_2) = (flag_SCION, flag_EPIC-HP, flag_COLIBRI, ...)

The order of the interfaces, (If_1, If_2) vs. (If_2, If_1), allows to enable and disable different types of traffic depending on the direction. 
To exclusively allow SCION traffic (default) between interfaces 'x' and 'y' we would set:

AllowedTraffic(x, y) = (1, 0, 0, ...)

And similarly to only allow EPIC-HP traffic:

AllowedTraffic(x, y) = (0, 1, 0, ...)


## Best Practices

### Highly Secure Hidden Paths
The last and penultimate AS on the hidden path only allow EPIC-HP traffic on the affected interface pair.
Add image 

### 

