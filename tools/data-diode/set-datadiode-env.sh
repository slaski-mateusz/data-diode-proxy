#!/usr/bin/bash

# Receiving interface
export DATADIODE_IN_IP=127.0.0.1
export DATADIODE_IN_PORT=8001

# Sending Interface
export DATADIODE_OUT_IP=127.0.0.1
export DATADIODE_OUT_PORT=8002

# Simulating network bad quality. By dropping some percentage of packets
export DATADIODE_DROP_PCT=0

# Simulating network Failutes. By stopping sending packages for defined time in seconds.
export DATADIODE_STOP_CYCLE=0   # Cycle in seconds to repeat stopping communication
export DATADIODE_STOP_RND=0     # Cycle randomisation in seconds
export DATADIODE_STOP_MIN=0     # Simulated failure minimum time in seconds
export DATADIODE_STOP_MAX=0     # Simulated failure maximum time in seconds
