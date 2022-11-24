# Protocol Design

## Ports

| Version | Source Port | Destination Port |
| ------- | ----------- | ---------------- |
| 2       | `63321`     | `63322`          |
| 1       | `63323`     | `63324`          |

## Structure

| Field                 | Type      | Description                      |
| --------------------- | --------- | -------------------------------- |
| Version               | byte      | always `1`                       |
| [Command](#command)   | byte      |                                  |
| [Result](#result)     | uint16 BE |                                  |
| Failure TLV           | 4 bytes   | Failure TLV                      |
| Manager ID            | 6 bytes   | Device MAC-address (Source)      |
| [Agent ID](#agent-id) | 6 bytes   | Device MAC-address (Destination) |
| Sequence              | uint32 BE | Sequence ID                      |
| Signature             | 4 bytes   | alaways `NSDP` as ASCII-encoded  |
| [Tags](#tags)         | TLVs      |                                  |

### Command

| Value | Meaning        |
| ----- | -------------- |
| 0x01  | Read Request   |
| 0x02  | Read Response  |
| 0x03  | Write Request  |
| 0x04  | Write Response |

### Result

| Value  | Meaning          |
| ------ | ---------------- |
| 0x0000 | Success          |
| 0x7000 | Invalid Password |

### Agent ID

If agent id value is `00:00:00:00:00:00` is used as multicast address, request will be proceeded by all devices, which
would receive it

## Tags

### Tag Structure

| Field  | Type      |
| ------ | --------- |
| Tag    | uint16 BE |
| Length | uint16 BE |
| Value  | N         |

### Tag Definitons

| Tag      | Description                       | Value Type         |
| -------- | --------------------------------- | ------------------ |
| `0x0000` | Start of Mark                     | (empty)            |
| `0x0001` | Device Model                      | `string`           |
| `0x0002` |                                   |                    |
| `0x0003` | Device Given Name                 | `string`           |
| `0x0004` | Device MAC-address                | `net.HardwareAddr` |
| `0x0005` | Device system location            | `string`           |
| `0x0006` | IP Address                        | `net.IP`           |
| `0x0007` | IP Mask                           | `net.IP`           |
| `0x0008` | Gateway                           | `net.IP`           |
| `0x0009` | New Password                      |                    |
| `0x000a` | Administration Password           |                    |
| `0x000b` | DHCP Mode                         | `uint8`            |
| `0x000c` | Active Firmware                   | `uint8`            |
| `0x000d` | Firmware (Slot 1)                 | `string`           |
| `0x000e` | Firmware (Slot 2)                 | `string`           |
| `0x000f` | Next Active Firmware              | `uint8`            |
| `0x0013` | Reboot                            |                    |
| `0x0014` |                                   |                    |
| `0x0017` | Auth v2 Password Salt             |                    |
| `0x001a` | Auth v2 Password                  |                    |
| `0x0400` | Factory Reset                     |                    |
| `0x0c00` | Link status of ports              |                    |
| `0x1000` | Port Traffic Statistic (Request)  |                    |
| `0x1400` | Port Traffic Statistic (Reset)    |                    |
| `0x1800` | Test Cable Request                |                    |
| `0x1c00` | Test Cable Result                 |                    |
| `0x2000` | VLAN Support                      |                    |
| `0x2400` | VLAN Members (Port)               |                    |
| `0x2800` | VLAN Members (dot1q)              |                    |
| `0x2c00` | Delete VLAN (write only)          |                    |
| `0x3000` | VLAN PVID (dot1q)                 |                    |
| `0x3400` | QoS Engine                        |                    |
| `0x3800` | Port based QoS Priority           |                    |
| `0x4c00` | Ingress Bandwidth limit           |                    |
| `0x5000` | Egress Bandwidth limit            |                    |
| `0x5400` | Broadcast Filtering               |                    |
| `0x5800` | Broadcast Bandwidth               |                    |
| `0x5c00` | Port Mirroring                    |                    |
| `0x6000` | Number of Ports                   | `uint8`            |
| `0x6800` | IGMP Snooping Status              |                    |
| `0x6c00` | Block Unknown Multicast Traffic   |                    |
| `0x7000` | IGMPv3 IP Header Validation       |                    |
| `0x7400` |                                   |                    |
| `0x7800` | Serial Number                     | `string`           |
| `0x8000` | IGMP Snooping Static Router Ports |                    |
| `0x9000` | Loop detection                    |                    |
| `0xffff` | End of Mark                       | (empty)            |

## References

- <https://en.wikipedia.org/wiki/Netgear_Switch_Discovery_Protocol>
- <https://github.com/kamiraux/wireshark-nsdp/blob/master/NSDP_info>
- <https://github.com/yaamai/go-nsdp>
- <https://github.com/bengal/nsdpc>
