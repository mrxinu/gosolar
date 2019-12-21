# Solarcmd

CLI tool based on [gosolar](https://github.com/stobias123/gosolar). Makes interacting with solarwinds IPAM much easier.

Error handling is probably pretty poor.

## Installation

```
go install github.com/stobias123/gosolar/solarcmd
```

## Examples

**Find a Subnet**
```
$ solarcmd subnet list | jq '.[] | [ .Address, .CIDR, .DisplayName ]'
{
    "192.168.1.0",
    "24",
    "Foobar Network"
}
...
```

**Get an IP address.**

Use the above information to get a free IP.
```
$ solarcmd ip get 192.168.1.0 --cidr 24
{"IpNodeID": 0, "IPAddress": "192.168.1.34", "Status": 0, "StatusString": "", "Comments","" }

```

**Reserve that IP**
```
$ solarcmd ip reserve 192.168.1.34
{"IpNodeID": 0, "IPAddress": "192.168.1.34", "Status": 0, "StatusString": "", "Comments","" }
```
