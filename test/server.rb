require 'sinatra'
set :port, 17778

get '/' do
    'hello'
end

=begin        {
        "SubnetId": 1235,
        "Address": "10.199.154.0",
        "CIDR": "23",
        "FriendlyName": "DEV/QA NFS2",
        "DisplayName": "DEV/QA NFS2",
        "AvailableCount": 222,
        "ReservedCount": 6,
        "UsedCount": 191,
        "totalCount": 512,
        "Comments": "NFS - VLAN 412",
        "VLAN": 412,
        "AddressMask": "255.255.254.0"
    } 
=end
    
subnet_obj = '''
{
    "totalRows":1,
        "results": [
        {
            "SubnetId": 1234,
            "Address": "10.199.152.0",
            "CIDR": "23",
            "FriendlyName": "test subnet",
            "DisplayName": "test subnet",
            "AvailableCount": 200,
            "ReservedCount": 2,
            "UsedCount": 181,
            "totalCount": 512,
            "Comments": "NFS - VLAN 410",
            "VLAN": 410,
            "AddressMask": "255.255.254.0"
        }
    ]
}
'''

post '/SolarWinds/InformationService/v3/Json/Invoke/IPAM.SubnetManagement/GetFirstAvailableIp' do
  return '"192.168.1.23"'
end

post '/SolarWinds/InformationService/v3/Json/Query' do
    data = JSON.parse( request.body.read.to_s )
    ## Get the post params
    puts data["parameters"]
    ipam_query = data["query"].include? "IPAM"
    if ipam_query
        return '{"totalRows":1,
        "results": [
            {
                "SubnetId": 1234,
                "Address": "10.199.152.0",
                "CIDR": 23,
                "FriendlyName": "test subnet",
                "DisplayName": "test subnet",
                "AvailableCount": 200,
                "ReservedCount": 2,
                "UsedCount": 181,
                "totalCount": 512,
                "Comments": "NFS - VLAN 410",
                "VLAN": "410",
                "AddressMask": "255.255.254.0"
            }
        ]}'
    end
    if data['parameters']['name'].include? "test"
        return '{"totalRows":1,
        "results": [
            {
                "SubnetId": 1234,
                "Address": "10.199.152.0",
                "CIDR": "23",
                "FriendlyName": "test subnet",
                "DisplayName": "test subnet",
                "AvailableCount": 200,
                "ReservedCount": 2,
                "UsedCount": 181,
                "totalCount": 512,
                "Comments": "NFS - VLAN 410",
                "VLAN": 410,
                "AddressMask": "255.255.254.0"
            }
        ]}'
    end

    if data["query"].include? 'DisplayName == "DEV/QA NFS"'
        return '{"totalRows":1,
        "results": [
            {
                "SubnetId": 1234,
                "Address": "10.199.152.0",
                "CIDR": "23",
                "FriendlyName": "DEV/QA NFS",
                "DisplayName": "DEV/QA NFS",
                "AvailableCount": 200,
                "ReservedCount": 2,
                "UsedCount": 181,
                "totalCount": 512,
                "Comments": "NFS - VLAN 410",
                "VLAN": 410,
                "AddressMask": "255.255.254.0"
            }
        ]}'
    end
    return '{"results":[]}'
end

post '/SolarWinds/InformationService/v3/Json/Invoke/IPAM.SubnetManagement/ChangeIPStatus' do
    return '[{"IPNodeID": 123,"Address": "192.168.33.2","Status": 2,"StatusString": "Available","Comments": ""}]'
end