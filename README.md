# Data Diode Proxy

Core components - Implementation in Golang
Supporting tools - Implementation in Python

## Introduction

[Data Diode or Unidirectional Network](https://en.wikipedia.org/wiki/Unidirectional_network) device (called later DD) is most secure network connection.

The disadvantages are:

* Package receive can't be confirmed so only UDP may be used.
* For the same reason transmitter can't resend data when some package is lost.
  * It is easy to get inconsistent data.
  * Solution is redundancy - Sending the same data cyclically multiple times.
* Data-diode devices has often limited package size. It may be configurable.

Such device requires dedicated proxy software to transmit data.

Software must consist of two programs:

* Transmission Agent - TX
  This program has to:
  * Read source data.
  * Split it into packages fit to DD allowed package size.
  * Put it into sending buffer with configured TTL (time to live) and with unique ascending ids.
  * Send it cyclically until TTL is reached.
  * Remove outdated package from buffer
  * Have option to resend outdated packages for selected time on demand - See RX functionallity.
* Receiver Agent - RX
  This program has to:
  * Receive packages
  * If has no such package in receive buffer - Put in receive buffer.
  * If package is already in receive buffer - Drop it.
  * Monitor if for configured TTL data in receive buffer is consistent - no gaps in id numbers.
    * In case of inconsistency indicate it (logs / ui / sending some message)
    * Indication has to contains data when inconsistencies appeared.

## Implementation

I decided to start from easiest and most universal case Data-Diode-file-proxy.

It may be used or adapted for many cases - For example:

* Logs replication
* Database (MariaDB) binlog replication

There are four applications:

* **tx-sim** - Application simulating files data source.
* **tx-agent** - Data Diode Transmission Agent described above.
* **rx-agent** - Data Diode Receiver Agent described above.
* **rx-ver** - Application verifying if received files are consistent.

Simulation and verification are supporting / test purpose tools so they are Python applications.
TX and RX Agent need have good performance so they are Golang applications.

### tx-sim

Program writes files named according the configuration filled with 32 characters lines:

* 8 characters with zero filled file number
* 23 characters with zero filled counter in file
* 1 newline

**Default configuration:**

```yaml:no-line-numbers
out-file:
  path: ../tx-workdir
  name: sim
  ext: dat
  min-size-kb: 1
  max-size-kb: 2
cycle:
  seconds: 5
  random-sec-offset: 2
```

Is placed in YAML file with the same name as script.

Script can be executed with option "-c", "--conf" and config file name

### rx-ver

Program check if each incoming file:

* has subsequent number
* is consistent - contains file number and counter values are incremented without gaps.
