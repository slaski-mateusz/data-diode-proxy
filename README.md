# Data Diode Proxy

Core components - Implementation in Golang
Supporting tools - Implementation in Python

## Introduction

[Data Diode or Unidirectional Network](https://en.wikipedia.org/wiki/Unidirectional_network) device (called later DD / Data Diode) is most secure network connection.

The disadvantages are:

* Package receive can't be confirmed so only UDP may be used.
* For the same reason transmitter can't resend data when some package is lost.
  * It is easy to get inconsistent data.
  * Solution is redundancy - Sending the same data cyclically multiple times.
* Data-diode devices has often limited package size. It may be configurable.

Such device requires dedicated proxy software to transmit data.

Core of software would be two applications:

* Transmission Agent - TX
* Receiver Agent - RX


## Implementation

I decided to start from easiest and most universal case Data-Diode-file-proxy.

It may be used or adapted for many cases - For example:

* Logs replication
* Database (MariaDB) binlog replication

To solve five applications are needed:

* **tx-sim** - Application simulating files data source.
* **tx-agent** - Data Diode Transmission Agent described above.
* **dd-proxy** - Simple proxy simulating data diode.
* **rx-agent** - Data Diode Receiver Agent described above.
* **rx-ver** - Application verifying if received files are consistent.

Simulation and verification are supporting / test purpose tools so they are Python applications.
TX and RX Agent need have good performance so they are Golang applications.

### tx-sim

Developed in Python.

Application simulates incoming data stored in files.

Program writes files named according the configuration. Files are filled with 32 characters lines:

* 8 characters with zero filled file number
* 23 characters with zero filled counter in file
* 1 newline

Script can be executed with option "-c", "--conf" and config file name.

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
errors:
  skip-file-chance: 0
  skip-line-chance: 0
```

Is placed in YAML file with the same name as script.

To be able to test **rx-ver** application simulator needs to generate dome errors.

It would be activated with options dictionary **errors** from config file. This section is optional. Without it chances are set to zero - No errors generated.

* **"skip-file-chance"**
  With value 0-100. Sometimes some file number would be skipped with propability given by this value.
* **"skip-line-chance"**
  With value 0-100. Sometimes line in file would be skipped with propability given by number.

### tx-agent

Developed in Go

This program has to:

* Read source data.
* Split it into packages fit to Data Diode allowed package size.
* Put it into sending buffer with configured TTL (time to live) and with unique ascending ids.
* Send it cyclically until TTL is reached.
* Remove outdated package from buffer
* Have option to resend outdated packages for selected time on demand - See RX functionallity.

### dd-proxy

Developed in Go

ApplicationListening on one interface. Sending packages on other.

The only more sophisticated function is possibility to add some errors:

* Loose some amount of randomly selected packages
  Option **"-p --packets_drops"**. With value 0-100. Number set percentage possibility to drop package.
* Eventually schedule period of inactivity to simulate device failure).

### rx-agent

Developed in Go

This program has to:

* Receive packages
* If has no such package in receive buffer - Put in receive buffer.
* If package is already in receive buffer - Drop it.
* Monitor if for configured TTL data in receive buffer is consistent - no gaps in id numbers.
  * In case of inconsistency indicate it (logs / ui / sending some message)
    * If TTL for package is actual than only show some statistics
    * If TTL for package is outdated strongly indicate and notify that data has to be restored manually.  
  * Indication has to contains information when inconsistencies appeared.

### rx-ver

Developed in Python

Program check if each incoming file:

* has subsequent number
* is consistent - contains file number and counter values are incremented without gaps.
